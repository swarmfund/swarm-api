package api

import (
	"errors"

	"gitlab.com/swarmfund/api/db2/api"
	"gitlab.com/swarmfund/api/render/hal"
	"gitlab.com/swarmfund/api/render/problem"
	"gitlab.com/swarmfund/api/resource"
)

type GetTfaKeychainAction struct {
	Action
	Username string
	WalletId string
	wallet   *api.Wallet
}

func (action *GetTfaKeychainAction) JSON() {
	action.ValidateBodyType()
	action.Do(
		action.loadParams,
		action.performRequest,
		func() {
			var response resource.TfaKeychain
			//response.TfaKeychain = action.wallet.TfaKeychainData
			//response.TfaSalt = action.wallet.TfaSalt
			hal.Render(action.W, response)
		})
}

func (action *GetTfaKeychainAction) loadParams() {
	action.ValidateBodyType()
	action.Username = action.GetNonEmptyString("username")
	action.WalletId = action.GetByteArray("walletId", 32)
}

func (action *GetTfaKeychainAction) performRequest() {
	var err error
	action.wallet, err = action.APIQ().Wallet().ByEmail(action.Username)
	if err != nil {
		action.Log.WithField("err", err.Error()).Error("Unable to get wallet")
		action.Err = &problem.ServerError
		return
	}

	if action.wallet == nil {
		action.Err = &problem.NotFound
		return
	}

	if !action.wallet.Verified {
		action.Err = &problem.Forbidden
		return
	}

	if action.WalletId != action.wallet.WalletId {
		action.SetInvalidField("walletId", errors.New(" mismatched"))
		return
	}
}
