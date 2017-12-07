package notificator

import (
	"bytes"
	"fmt"
	"net/url"

	"github.com/pkg/errors"
	"gitlab.com/distributed_lab/notificator"
	"gitlab.com/swarmfund/api/config"
	"gitlab.com/swarmfund/api/db2/api"
	"gitlab.com/swarmfund/api/internal/clienturl"
)

const (
	NotificatorTypeVerificationEmail     = 3
	NotificatorTypeApprovalEmail         = 4
	NotificatorTypeTFA                   = 5
	NotificatorTypeRecoveryRequest       = 6
	NotificatorTypeAdminNotification     = 7
	NotificatorTypeKYCReviewPending      = 7
	NotificatorTypeLoginNotification     = 6
	NotificatorTypePaymentNotification   = 8
	NotificatorTypeOperationNotification = 9
)

type Connector struct {
	notificator *notificator.Connector
	conf        config.Notificator
}

func NewConnector(conf config.Notificator) *Connector {
	// TODO move this to config
	endpoint, err := url.Parse(conf.Endpoint)
	if err != nil {
		panic(err)
	}

	return &Connector{
		notificator: notificator.NewConnector(
			notificator.Pair{Secret: conf.Secret, Public: conf.Public},
			*endpoint,
		),
		conf: conf,
	}
}

func (c *Connector) SendVerificationLink(email string, payload clienturl.Payload) error {
	encoded, err := payload.Encode()
	if err != nil {
		return errors.Wrap(err, "failed to encode payload")
	}
	letter := &Letter{
		Header: "Swarm Fund Email Verification",
		Link:   fmt.Sprintf("%s/%s", c.conf.ClientRouter, encoded),
	}

	var buff bytes.Buffer
	err = c.conf.EmailConfirmation.Execute(&buff, letter)
	if err != nil {
		return errors.Wrap(err, "Error while populating template for notify approval")
	}
	msg := &notificator.EmailRequestPayload{
		Destination: email,
		Subject:     letter.Header,
		Message:     buff.String(),
	}

	return c.send(NotificatorTypeVerificationEmail, email, msg)
}

func (c *Connector) NotifyKYCReviewPending(recipient string) error {
	//header := "KYC reviews pending"
	//
	//var buff bytes.Buffer
	//err := c.conf.KYCReviewNotification.Template.Execute(&buff, struct{}{})
	//if err != nil {
	//	return errors.Wrap(err, "failed to render template")
	//}
	//
	//payload := &notificator.EmailRequestPayload{
	//	Destination: recipient,
	//	Subject:     header,
	//	Message:     buff.String(),
	//}
	//
	//if err = c.send(NotificatorTypeApprovalEmail, recipient, payload); err != nil {
	//	c.log.WithError(err).Error("failed to send email")
	//	return err
	//}

	return nil
}

func (c *Connector) NotifyApproval(email string) error {
	header := "Swarm Fund Account Approved"
	msg := "Your account was just approved at Swarm Fund. "

	letter := Letter{Header: header, Body: msg, Link: ""}
	return c.sendKycNotification(email, letter)
}

func (c *Connector) NotifyRejection(email string) error {
	header := "Swarm Fund Account Rejected"
	msg := "Your request was rejected by the administrator. Log in to your account for more details."

	letter := Letter{Header: header, Body: msg, Link: ""}
	return c.sendKycNotification(email, letter)
}

func (c *Connector) sendKycNotification(email string, letter Letter) error {
	//var buff bytes.Buffer
	//err := c.conf.KYCApproval.Template.Execute(&buff, letter)
	//if err != nil {
	//	log.WithField("err", err.Error()).Error("Error while populating template for notify approval")
	//	return err
	//}
	//
	//payload := &notificator.EmailRequestPayload{
	//	Destination: email,
	//	Subject:     letter.Header,
	//	Message:     buff.String(),
	//}
	//
	//err = c.send(NotificatorTypeApprovalEmail, email, payload)
	//if err != nil {
	//	c.log.WithError(err).Error("Failed to send email")
	//	return err
	//}

	return nil
}

func (c *Connector) SendRecoveryRequest(email, token, rejectReason string) (*notificator.Response, error) {
	//letter := Letter{
	//	Header: "Swarm Fund wallet recovery",
	//	Link: fmt.Sprintf("%s/wallets/recovery?token=%s",
	//		c.conf.RecoveryRequest.RedirectURL, token),
	//}
	//
	//if rejectReason != "" {
	//	letter.Body = fmt.Sprintf(`Your request was rejected: "%s"`, rejectReason)
	//}
	//
	//var buff bytes.Buffer
	//err := c.conf.RecoveryRequest.Template.Execute(&buff, letter)
	//if err != nil {
	//	log.WithField("err", err.Error()).Error("failed to render template")
	//	return nil, err
	//}
	//
	//payload := &notificator.EmailRequestPayload{
	//	Destination: email,
	//	Subject:     letter.Header,
	//	Message:     buff.String(),
	//}
	//
	//response, err := c.notificator.Send(NotificatorTypeRecoveryRequest, email, payload)
	//return response, err
	return nil, nil
}

func (c *Connector) SendNewDeviceLogin(email string, device api.AuthorizedDevice) error {
	//letter := LoginNoticeLetter{
	//	Header:       "Swarm Fund",
	//	BrowserFull:  device.Details.BrowserFull,
	//	BrowserShort: device.Details.Browser,
	//	Date:         time.Now().Format("Mon Jan _2 15:04:05 2006"),
	//	DeviceFull:   device.Details.OSFull,
	//	DeviceShort:  device.Details.OS,
	//	Ip:           device.Details.IP,
	//	Location:     device.Details.Location,
	//}
	//
	//var buff bytes.Buffer
	//err := c.conf.LoginNotification.Template.Execute(&buff, letter)
	//if err != nil {
	//	log.WithField("err", err.Error()).Error("failed to render template")
	//	return err
	//}
	//
	//payload := &notificator.EmailRequestPayload{
	//	Destination: email,
	//	Subject:     letter.Header,
	//	Message:     buff.String(),
	//}
	//
	//_, err = c.notificator.Send(NotificatorTypeLoginNotification, email, payload)
	//return err
	return nil
}

func (c *Connector) SendOperationNotice(opType int, letter TransferNoticeI) error {
	//if letter.GetEmail() == "" {
	//	// for case when participants is not present is api db
	//	return nil
	//}
	//
	//var buff bytes.Buffer
	//err := letter.AddLoginLink(c.conf.EmailConfirmation.ClientURL)
	//if err != nil {
	//	return err
	//}
	//
	//switch opType {
	//case PAYMENT:
	//	err = c.conf.OperationsNotification.Payment.Execute(&buff, letter)
	//case DEMURRAGE:
	//	err = c.conf.OperationsNotification.Demurrage.Execute(&buff, letter)
	//case DEPOSIT:
	//	err = c.conf.OperationsNotification.CoinsEmission.Execute(&buff, letter)
	//case OFFER:
	//	err = c.conf.OperationsNotification.Offer.Execute(&buff, letter)
	//case FORFEIT:
	//	err = c.conf.OperationsNotification.Forfeit.Execute(&buff, letter)
	//case INVOICE:
	//	err = c.conf.OperationsNotification.Invoice.Execute(&buff, letter)
	//default:
	//	err = errors.New("Unknown operation notification type")
	//}
	//
	//if err != nil {
	//	log.WithField("err", err.Error()).Error("failed to render payment notice template")
	//	return err
	//}
	//
	//payload := &notificator.EmailRequestPayload{
	//	Destination: letter.GetEmail(),
	//	Subject:     letter.GetHeader(),
	//	Message:     buff.String(),
	//}
	//// c.log.WithField("destination", payload.Destination).WithField("Subject", payload.Subject).Info("Try to send")
	//_, err = c.notificator.Send(NotificatorTypeOperationNotification, letter.GetToken(), payload)
	//return err
	return nil
}

func (c *Connector) send(requestType int, token string, payload notificator.Payload) error {
	response, err := c.notificator.Send(requestType, token, payload)
	if err != nil {
		return err
	}

	if !response.IsSuccess() {
		return errors.New("notification request not accepted")
	}
	return nil
}
