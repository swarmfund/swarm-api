package operation_notifier

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"gitlab.com/swarmfund/api/db2/api"
	apErrors "gitlab.com/swarmfund/api/errors"
	"gitlab.com/swarmfund/api/log"
	"gitlab.com/swarmfund/api/notificator"
	responseTP "gitlab.com/swarmfund/api/operation_notifier/response_types"
	"gitlab.com/swarmfund/go/keypair"
	"gitlab.com/swarmfund/go/xdr"
	"gitlab.com/swarmfund/horizon-connector"
	ssego "gitlab.com/distributed_lab/sse-go"
	"github.com/pkg/errors"
)

const PaymentsPath = "/payments?order=asc"

var ErrorUnsupportedOpType = errors.New("Unsupported operation type")

type Notifier struct {
	apiQ             api.QInterface
	horizonURL       string
	kp               keypair.KP
	cursor           string
	log              *log.Entry
	sender           *notificator.Connector
	sse              *ssego.Listener
	timeBound        time.Time
	skipOldOperation bool
}

func New(cursor, horizonURL string, kp keypair.KP, sender *notificator.Connector, apiQ api.QInterface) *Notifier {
	return &Notifier{
		apiQ:             apiQ,
		cursor:           cursor,
		horizonURL:       horizonURL,
		kp:               kp,
		log:              log.WithField("service", "operation-notifier"),
		sender:           sender,
		timeBound:        time.Now().AddDate(0, 0, -1),
		skipOldOperation: cursor == "",
	}
}
func (n *Notifier) NewListener() {
	n.sse = ssego.NewListener(n.makeRequest)
}

func (n *Notifier) makeRequest() (*http.Request, error) {
	path := fmt.Sprintf("%s&cursor=%s", PaymentsPath, n.cursor)
	n.log.WithField("cursor", n.cursor).WithField("path", path).Warn("Remake request")
	return horizon.NewSignedRequest(n.horizonURL, "GET", path, n.kp)
}

func (n *Notifier) Run() {
	if n.skipOldOperation {
		n.log.Warn(fmt.Sprintf("Last cursor is not set. App operations before %s will be ignored", n.timeBound.Format(notificator.TimeLayout)))
	}
	for event := range n.sse.Events() {
		if event.Err != nil {
			n.log.WithField("last_cursor", n.cursor).WithError(event.Err).Error("Failed to get event")
			time.Sleep(30 * time.Second)
			continue
		}

		rawOperation, err := ioutil.ReadAll(event.Data)
		if err != nil {
			n.log.WithField("last_cursor", n.cursor).WithError(err).Error("Failed to unmarshal op")
			time.Sleep(30 * time.Second)
			continue
		}

		c, err := n.processOperation(rawOperation)
		if err != nil {
			n.log.WithField("last_cursor", n.cursor).WithError(err).Error("Failed to process operation")
			time.Sleep(30 * time.Second)
			continue
		}
		n.cursor = c
		n.log.WithField("last_cursor", n.cursor).Info("Process operations in progress")
	}
}

func (n *Notifier) processOperation(rawOperation []byte) (cursor string, err error) {
	defer func() {
		if rec := recover(); rec != nil {
			err := apErrors.FromPanic(rec)
			n.log.WithStack(err).WithError(err).Error("Process Operation recovered")
		}
	}()

	base := new(responseTP.OperationBase)
	err = json.Unmarshal(rawOperation, base)
	if err != nil {
		n.log.WithError(err).Warn("Failed to unmarshal operation base")
		return cursor, errors.Wrap(err, "Failed to unmarshal operation base")
	}

	if n.skipOldOperation && base.LedgerCloseTime.Before(n.timeBound) {
		return base.ID, nil
	} else {
		n.skipOldOperation = false
	}

	err = base.Populate()
	if err != nil {
		n.log.WithError(err).Warn("Failed to normalize operation base")
		return cursor, errors.Wrap(err, "Failed to normalize operation base")
	}

	cursor = base.ID
	err = n.loadParticipantsDetails(base, "")
	if err != nil {
		n.log.WithError(err).Error("Failed to load participants details")
		return cursor, err
	}

	if len(base.Participants) == 0 {
		return cursor, nil
	}

	switch xdr.OperationType(base.TypeI) {
	case xdr.OperationTypePayment:
		return cursor, n.paymentNotification(base, rawOperation)
	case xdr.OperationTypeManageInvoice:
		return cursor, n.manageInvoiceNotification(base, rawOperation)
	case xdr.OperationTypeDemurrage:
		return cursor, n.demurrageNotification(base, rawOperation)
	case xdr.OperationTypeManageCoinsEmissionRequest, xdr.OperationTypeReviewCoinsEmissionRequest:
		return cursor, n.coinsEmissionNotification(base, rawOperation)
	case xdr.OperationTypeManageForfeitRequest:
		return cursor, n.manageForfeitRequestNotification(base, rawOperation)
	case xdr.OperationTypeManageOffer:
		return cursor, n.manageOfferNotification(base, rawOperation)
	default:
		return cursor, ErrorUnsupportedOpType
	}
}

func (n *Notifier) loadParticipantsDetails(op responseTP.OperationI, contactsFor string) error {
	//participantsMap := op.ParticipantsMap()
	//err := n.apiQ.Users().Participants(participantsMap, contactsFor)
	//if err != nil {
	//	return err
	//}
	//op.UpdateParticipants(participantsMap)
	return nil
}

func (n *Notifier) paymentNotification(base *responseTP.OperationBase, rawOperation []byte) error {
	payment := responseTP.Payment{}
	err := payment.FromBase(base, rawOperation)
	if err != nil {
		n.log.WithError(err).Warn("Failed to unmarshal payment details")
		return errors.Wrap(err, "Failed to unmarshal payment details")
	}

	// load participants as for sender
	// and send email to them
	err = n.loadParticipantsDetails(&payment, payment.From)
	if err != nil {
		n.log.WithError(err).Warn("Failed to load payment participants | 1")
		return errors.Wrap(err, "Failed to load payment participants | 1")
	}

	err = n.sender.SendOperationNotice(notificator.PAYMENT, payment.ToLetter(payment.From, true))
	if err != nil {
		n.log.WithError(err).Warn("Failed to send payment letter | 1")
		return errors.Wrap(err, "Failed to send payment letter | 1")
	}

	// load participants as for receiver
	// and send email to them
	err = n.loadParticipantsDetails(&payment, payment.To)
	if err != nil {
		n.log.WithError(err).Warn("Failed to load payment participants | 2")
		return errors.Wrap(err, "Failed to load payment participants | 2")
	}

	err = n.sender.SendOperationNotice(notificator.PAYMENT, payment.ToLetter(payment.To, false))
	if err != nil {
		n.log.WithError(err).Warn("Failed to send payment letter | 2")
		return errors.Wrap(err, "Failed to send payment letter | 2")
	}

	return nil
}

func (n *Notifier) manageInvoiceNotification(base *responseTP.OperationBase, rawOperation []byte) error {
	invoice := responseTP.ManageInvoice{}
	err := invoice.FromBase(base, rawOperation)
	if err != nil {
		n.log.WithError(err).Warn("Failed to unmarshal manage invoice details")
		return errors.Wrap(err, "Failed to unmarshal manage invoice details")
	}

	// load participants as for sender
	// and send email to them
	err = n.loadParticipantsDetails(&invoice, invoice.SourceAccount)
	if err != nil {
		n.log.WithError(err).Warn("Failed to load invoice participants | 1")
		return errors.Wrap(err, "Failed to load invoice participants | 1")
	}

	err = n.sender.SendOperationNotice(notificator.INVOICE, invoice.ToLetter(invoice.SourceAccount, false))
	if err != nil {
		n.log.WithError(err).Warn("Failed to send invoice letter | 1")
		return errors.Wrap(err, "Failed to send invoice letter | 1")
	}

	// load participants as for receiver
	// and send email to them
	err = n.loadParticipantsDetails(&invoice, invoice.Sender)
	if err != nil {
		n.log.WithError(err).Warn("Failed to load invoice participants | 2")
		return errors.Wrap(err, "Failed to load invoice participants | 2")
	}

	err = n.sender.SendOperationNotice(notificator.INVOICE, invoice.ToLetter(invoice.Sender, true))
	if err != nil {
		n.log.WithError(err).Warn("Failed to send invoice letter | 2")
		return errors.Wrap(err, "Failed to send invoice letter | 2")
	}

	return nil
}

func (n *Notifier) coinsEmissionNotification(base *responseTP.OperationBase, rawOperation []byte) error {
	deposit := responseTP.ReviewCoinsEmissionRequest{}
	err := deposit.FromBase(base, rawOperation)
	if err != nil {
		n.log.WithError(err).Warn("Failed to unmarshal review coins emission details")
		return errors.Wrap(err, "Failed to unmarshal review coins emission details")
	}

	isManage := xdr.OperationType(base.TypeI) == xdr.OperationTypeManageCoinsEmissionRequest
	err = n.sender.SendOperationNotice(notificator.DEPOSIT, deposit.ToLetter(isManage))
	if err != nil {
		n.log.WithError(err).Warn("Failed to send deposit letter")
		return errors.Wrap(err, "Failed to send deposit letter")
	}

	return nil
}

func (n *Notifier) manageForfeitRequestNotification(base *responseTP.OperationBase, rawOperation []byte) error {
	forfeit := responseTP.ManageForfeitRequest{}
	err := forfeit.FromBase(base, rawOperation)
	if err != nil {
		n.log.WithError(err).Warn("Failed to unmarshal manage forfeit details")
		return errors.Wrap(err, "Failed to unmarshal manage forfeit details")
	}

	err = n.sender.SendOperationNotice(notificator.FORFEIT, forfeit.ToLetter())
	if err != nil {
		n.log.WithError(err).Warn("Failed to send forfeit letter")
		return errors.Wrap(err, "Failed to send forfeit letter")
	}

	return nil
}

func (n *Notifier) demurrageNotification(base *responseTP.OperationBase, rawOperation []byte) error {
	for _, p := range base.Participants {
		if p.Email == "" {
			continue
		}

		letter, err := p.ToDemurrageLetter(base.ID, base.LedgerCloseTime)
		if err != nil {
			n.log.WithError(err).Warn("Failed to cast to demurrage letter")
			return errors.Wrap(err, "Failed to cast to demurrage letter")
		}

		err = n.sender.SendOperationNotice(notificator.DEMURRAGE, letter)
		if err != nil {
			n.log.WithError(err).Warn("Failed to send demurrage letter")
			return errors.Wrap(err, "Failed to send demurrage letter")
		}
	}
	return nil
}

func (n *Notifier) manageOfferNotification(base *responseTP.OperationBase, rawOperation []byte) error {
	for _, p := range base.Participants {
		if p.Email == "" {
			continue
		}

		letters, err := p.ToOfferLetters(base.ID, base.LedgerCloseTime)
		if err != nil {
			n.log.WithError(err).Warn("Failed to cast to offers letters")
			return errors.Wrap(err, "Failed to cast to offers letters")
		}

		for _, l := range letters {
			err = n.sender.SendOperationNotice(notificator.OFFER, &l)
			if err != nil {
				n.log.WithError(err).Warn("Failed to send offer letter")
				return errors.Wrap(err, "Failed to send offer letter")
			}
		}
	}
	return nil
}
