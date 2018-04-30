package api

import (
	"net/http"
	"net/url"

	"github.com/zenazn/goji/web"
	"gitlab.com/swarmfund/api/actions"
	"gitlab.com/swarmfund/api/db2"
	"gitlab.com/swarmfund/api/db2/api"
	"gitlab.com/swarmfund/api/httpx"
	"gitlab.com/swarmfund/api/log"
	"gitlab.com/swarmfund/api/render/problem"
	"gitlab.com/swarmfund/api/tfa"
	"gitlab.com/tokend/go/xdr"
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
