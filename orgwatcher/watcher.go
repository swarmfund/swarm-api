package orgwatcher

import (
	"encoding/json"
	"net/http"
	"net/url"
	"time"

	"gitlab.com/distributed_lab/logan/v3/errors"
	sse "gitlab.com/distributed_lab/sse-go"
	"gitlab.com/swarmfund/api/db2/api"
	"gitlab.com/swarmfund/api/log"
	"gitlab.com/swarmfund/go/keypair"
	horizon "gitlab.com/swarmfund/horizon-connector"
)

var (
	operationsEndpoint = url.URL{
		Path: "/operations",
	}
)

type SignerWatcher struct {
	sse        *sse.Listener
	signer     keypair.KP
	horizonURL url.URL
	cursor     string
	log        *log.Entry
	q          api.QInterface
}

func New(
	signer keypair.KP, log *log.Entry, q api.QInterface, horizonURL url.URL) *SignerWatcher {
	w := &SignerWatcher{
		signer:     signer,
		horizonURL: horizonURL,
		log:        log,
		q:          q,
	}
	w.sse = sse.NewListener(w.request)
	return w
}

func (w *SignerWatcher) request() (*http.Request, error) {
	u := w.horizonURL
	query := u.Query()
	query.Set("operation_type", "2")
	query.Set("cursor", w.cursor)
	u = *u.ResolveReference(&operationsEndpoint)
	u.RawQuery = query.Encode()
	return horizon.NewSignedRequest(
		w.horizonURL.String(), "GET", u.RequestURI(), w.signer)
}

func (w *SignerWatcher) Run() {
	ticker := time.NewTicker(30 * time.Second)
	for ; ; <-ticker.C {
		w.processEvents()
	}
}

func (w *SignerWatcher) processEvents() {
	defer func() {
		if r := recover(); r != nil {
			err := errors.FromPanic(r)
			w.log.WithError(err).Error("panicked")
		}
	}()
	var op SetOptionsOp
	var err error

	// get fresh cursor
	w.cursor, err = w.q.Wallet().OrganizationWatcherCursor()
	if err != nil {
		w.log.WithError(err).Error("failed to get default cursor")
		return
	}

	for event := range w.sse.Events() {
		if event.Err != nil {
			w.log.WithError(event.Err).Error("failed to get event")
			continue
		}

		err := json.NewDecoder(event.Data).Decode(&op)
		if err != nil {
			w.log.
				WithError(err).
				WithField("cursor", w.cursor).
				Error("failed to unmarshal op")
			continue
		}

		entry := w.log.WithField("operation_id", op.ID)
		entry.Debug("got operation")

		user, err := w.q.Users().ByAddress(op.SourceAccount)
		if err != nil {
			w.log.WithError(err).Error("failed to get user")
			return
		}

		if user == nil {
			entry.WithField("address", op.SourceAccount).Debug("user not found")
			continue
		}

		wallet, err := w.q.Wallet().ByCurrentAccountID(op.SignerKey)
		if err != nil {
			w.log.WithError(err).Error("failed to get wallet")
			return
		}

		if wallet == nil {
			entry.WithField("current_account_id", op.SignerKey).Debug("wallet not found")
			continue
		}

		isMaster := op.SourceAccount == wallet.AccountID
		isMinion := wallet.OrganizationAddress != nil && *wallet.OrganizationAddress == op.SourceAccount
		isOrgWallet := isMinion || isMaster

		entry.WithFields(log.F{
			"wallet_id":     wallet.Id,
			"is_master":     isMaster,
			"is_minion":     isMinion,
			"is_org_wallet": isOrgWallet,
		}).Debug("found wallet")

		if op.SignerWeight == 0 && isOrgWallet {
			if err = w.q.Wallet().Delete(wallet.Id); err != nil {
				w.log.WithError(err).Error("failed to delete wallet")
				return
			}
			continue
		}

		if op.SignerWeight != 0 && (wallet.OrganizationAddress == nil || *wallet.OrganizationAddress == op.SourceAccount) {
			if err = w.q.Wallet().UpdateOrganizationAttachment(int64(wallet.Id), op.SourceAccount, op.ID); err != nil {
				w.log.WithError(err).Error("failed to update attachment")
				return
			}
			continue
		}
	}
}
