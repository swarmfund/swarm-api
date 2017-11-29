package api

import (
	"errors"
	"net/http"
	"net/url"
	"strings"

	"github.com/zenazn/goji/web"
	"gitlab.com/swarmfund/api/actions"
	"gitlab.com/swarmfund/api/db2"
	"gitlab.com/swarmfund/api/db2/api"
	"gitlab.com/swarmfund/api/httpx"
	"gitlab.com/swarmfund/api/ledger"
	"gitlab.com/swarmfund/api/log"
	"gitlab.com/swarmfund/api/render/problem"
	"gitlab.com/swarmfund/api/tfa"
	"gitlab.com/swarmfund/api/toid"
	"gitlab.com/swarmfund/go/xdr"
	horizon "gitlab.com/swarmfund/horizon-connector"
)

// Action is the "base type" for all actions in horizon.  It provides
// structs that embed it with access to the App struct.
//
// Additionally, this type is a trigger for go-codegen and causes
// the file at Action.tmpl to be instantiated for each struct that
// embeds Action.
type Action struct {
	actions.Base
	App *App
	Log *log.Entry

	apiQ api.QInterface
}

// APIQ provides access to queries that access the horizon api database
func (action *Action) APIQ() api.QInterface {
	if action.apiQ == nil {
		action.apiQ = &api.Q{Repo: action.App.APIRepo(action.Ctx)}
	}

	return action.apiQ
}

// GetPagingParams modifies the base GetPagingParams method to replace
// cursors that are "now" with the last seen ledger's cursor.
func (action *Action) GetPagingParams() (cursor string, order string, limit uint64) {
	if action.Err != nil {
		return
	}

	cursor, order, limit = action.Base.GetPagingParams()

	if cursor == "now" {
		tid := toid.ID{
			LedgerSequence:   ledger.CurrentState().HistoryLatest,
			TransactionOrder: toid.TransactionMask,
			OperationOrder:   toid.OperationMask,
		}
		cursor = tid.String()
	}

	return
}

// GetPageQuery is a helper that returns a new db.PageQuery struct initialized
// using the results from a call to GetPagingParams()
func (action *Action) GetPageQuery() db2.PageQuery {
	if action.Err != nil {
		return db2.PageQuery{}
	}

	r, err := db2.NewPageQuery(action.GetPagingParams())

	if err != nil {
		action.Err = err
	}

	return r
}

// Prepare sets the action's App field based upon the goji context
func (action *Action) Prepare(c web.C, w http.ResponseWriter, r *http.Request) {
	base := &action.Base
	base.Prepare(c, w, r)
	action.App = action.GojiCtx.Env["app"].(*App)

	if action.Ctx != nil {
		action.Log = log.Ctx(action.Ctx)
	} else {
		action.Log = log.DefaultLogger
	}
}

// ValidateCursorAsDefault ensures that the cursor parameter is valid in the way
// it is normally used, i.e. it is either the string "now" or a string of
// numerals that can be parsed as an int64.
func (action *Action) ValidateCursorAsDefault() {
	if action.Err != nil {
		return
	}

	if action.GetString(actions.ParamCursor) == "now" {
		return
	}

	action.GetInt64(actions.ParamCursor)
}

// ValidateCursorWithinHistory compares the requested page of data against the
// ledger state of the history database.  In the event that the cursor is
// guaranteed to return no results, we return a 410 GONE http response.
func (action *Action) ValidateCursorWithinHistory() {
	if action.Err != nil {
		return
	}

	pq := action.GetPageQuery()
	if action.Err != nil {
		return
	}

	// an ascending query should never return a gone response:  An ascending query
	// prior to known history should return results at the beginning of history,
	// and an ascending query beyond the end of history should not error out but
	// rather return an empty page (allowing code that tracks the procession of
	// some resource more easily).
	if pq.Order != "desc" {
		return
	}

	var cursor int64
	var err error

	// HACK: checking for the presence of "-" to see whether we should use
	// CursorInt64 or CursorInt64Pair is gross.
	if strings.Contains(pq.Cursor, "-") {
		cursor, _, err = pq.CursorInt64Pair("-")
	} else {
		cursor, err = pq.CursorInt64()
	}

	if err != nil {
		action.Err = err
		return
	}

	elder := toid.New(ledger.CurrentState().HistoryElder, 0, 0)

	if cursor <= elder.ToInt64() {
		action.Err = &problem.BeforeHistory
	}
}

// BaseURL returns the base url for this requestion, defined as a url containing
// the Host and Scheme portions of the request uri.
func (action *Action) BaseURL() *url.URL {
	return httpx.BaseURL(action.Ctx)
}

func (action *Action) getMasterSignerType() int32 {
	result := int32(0)
	for i := range xdr.SignerTypeAll {
		result |= int32(xdr.SignerTypeAll[i])
	}
	return result
}

// If `token` is `nil` then TFA is fulfilled, otherwise caller should prompt client to enter OTP
func (action *Action) consumeTFA(wallet *api.Wallet, tfaAction string) {
	token := tfa.Token(int64(wallet.Id), tfaAction)
	// get active wallet tfa backends
	backends, err := action.APIQ().TFA().Backends(wallet.WalletId)
	if err != nil {
		action.Log.WithError(err).Error("failed to get tfa backend")
		action.Err = &problem.ServerError
		return
	}

	if len(backends) == 0 {
		// no backends are enabled, wallet is not tfa protected
		return
	}

	// try to consume tfa token
	ok, err := action.APIQ().TFA().Consume(token)
	if err != nil {
		action.Log.WithError(err).Error("failed to consume tfa")
		action.Err = &problem.ServerError
		return
	}

	if ok {
		// tfa token was already verified and now consumed
		return
	}

	// check if there is active token already
	otp, err := action.APIQ().TFA().Get(token)
	if err != nil {
		action.Log.WithError(err).Error("failed to get tfa")
		action.Err = &problem.ServerError
		return
	}

	var backend tfa.Backend
	if otp == nil {
		// no active tfa, let's go through backends and try create new one
		for _, record := range backends {
			if record.Priority <= 0 {
				continue
			}

			backend, err = record.Backend()
			if err != nil {
				action.Log.WithError(err).WithField("backend", record.ID).Error("failed to init backend")
				continue
			}
			//otpData, err := backend.OTPData()
			//if err != nil {
			//	action.Log.WithError(err).WithField("backend", record.ID).Error("failed to create tfa")
			//	continue
			//}
			otp = &api.TFA{
				BackendID: record.ID,
				//OTPData:   otpData,
				Token: token,
			}
			break
		}

		if otp == nil { // TODO FIX && len(backends) > 0 {
			// we failed to create tfa for account with backends enabled
			//action.Log.Error("failed to provide any tfa")
			//action.Err = &problem.ServerError
			return
		}

		err = action.APIQ().TFA().Create(otp)
		if err != nil {
			action.Log.WithError(err).Error("failed to store tfa")
			action.Err = &problem.ServerError
			return
		}
	}

	if backend == nil {
		record, err := action.APIQ().TFA().Backend(otp.BackendID)
		if err != nil {
			action.Log.WithError(err).WithField("backend", otp.BackendID).Error("failed to load backend")
			action.Err = &problem.ServerError
			return
		}

		if record == nil {
			// ok, we have tfa for backend not in the db
			// could happen if something went wrong while deleting it or someone messed up db
			// either way it's broken state
			action.Log.WithField("tfa", otp.ID).Error("couldn't find backend")
			action.Err = &problem.ServerError
			return
		}

		backend, err = record.Backend()
		if err != nil {
			action.Log.WithError(err).WithField("backend", otp.BackendID).Error("failed to init backend")
			action.Err = &problem.ServerError
			return
		}
	}

	// try to deliver notification by backend means
	//details, err := backend.Deliver(otp.OTPData)
	//if err != nil {
	//	action.Log.WithError(err).WithField("tfa", otp.ID).Error("failed to deliver")
	//	action.Err = &problem.ServerError
	//	return
	//}
	//action.Err = problem.TFARequired(otp.Token, details)
}

func (action *Action) recoveryApprove(recoveryRequest *api.RecoveryRequest, initialWallet, recoveryWallet *api.Wallet, tx string) error {
	account, err := action.App.horizon.AccountSigned(action.App.AccountManagerKP(), recoveryRequest.AccountID)
	if err != nil {
		return err
	}

	if account == nil {
		return errors.New("core account does not exists")
	}

	if tx == "" {
		err = action.App.horizon.Transaction(&horizon.TransactionBuilder{Source: action.App.MasterKP()}).
			Op(&horizon.RecoverOp{
				AccountID: initialWallet.AccountID,
				OldSigner: initialWallet.CurrentAccountID,
				NewSigner: recoveryWallet.CurrentAccountID,
			}).Sign(action.App.AccountManagerKP()).Submit()
	} else {
		err = action.App.horizon.SubmitTX(tx)
	}
	if err != nil {
		if serr, ok := err.(horizon.SubmitError); ok {
			action.Log.
				WithField("tx code", serr.TransactionCode()).
				WithField("op codes", serr.OperationCodes())
		}
		return err
	}

	// remove all wallets except one created during recovery flow
	err = action.APIQ().Wallet().SetActive(recoveryRequest.AccountID, recoveryWallet.WalletId)
	if err != nil {
		return err
	}

	// update user state
	err = action.APIQ().Users().SetRecoveryState(recoveryRequest.AccountID, api.UserRecoveryStateNil)
	if err != nil {
		return err
	}

	// should go last, so if something above fails admin can re-run it
	err = action.APIQ().Recoveries().Delete(recoveryRequest.AccountID)
	if err != nil {
		return err
	}

	return nil
}
