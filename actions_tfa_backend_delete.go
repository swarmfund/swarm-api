package api

import (
	"errors"

	"gitlab.com/swarmfund/api/db2/api"
	"gitlab.com/swarmfund/api/render/hal"
	"gitlab.com/swarmfund/api/render/problem"
	"gitlab.com/swarmfund/go/xdr"
)

type DeleteTFABackendsAction struct {
	Action

	username string
	wallet   *api.Wallet
}

func (action *DeleteTFABackendsAction) JSON() {
	action.Do(
		action.ValidateBodyType,
		action.loadParams,
		action.checkAllowed,
		action.loadWallet,
		action.checkAccountType,
		action.performRequest,
		func() {
			hal.Render(action.W, problem.Success)
		},
	)
}

func (action *DeleteTFABackendsAction) loadParams() {
	action.username = action.GetNonEmptyString("username")
}

func (action *DeleteTFABackendsAction) checkAllowed() {
	action.checkSignerConstraints(
		SignerType(action.App.CoreInfo.MasterAccountID, xdr.SignerTypeNotVerifiedAccManager),
	)
}

func (action *DeleteTFABackendsAction) loadWallet() {
	recovery, err := action.APIQ().Recoveries().ByUsername(action.username)
	if err != nil {
		action.Log.WithError(err).Error("failed to load recoveries")
		action.Err = &problem.ServerError
		return
	}

	if recovery == nil {
		action.wallet, err = action.APIQ().Wallet().ByEmail(action.username)
	} else {
		action.wallet, err = action.APIQ().Wallet().ByID(recovery.WalletID)
	}

	if err != nil {
		action.Log.WithError(err).Error("failed to load wallet")
		action.Err = &problem.ServerError
		return
	}

	if action.wallet == nil {
		action.SetInvalidField("username", errors.New("does not exists"))
		return
	}

}

func (action *DeleteTFABackendsAction) checkAccountType() {
	account, err := action.App.horizon.AccountSigned(action.App.AccountManagerKP(), action.wallet.AccountID)
	if err != nil {
		action.Log.WithError(err).Error("Can not load account info")
		action.Err = &problem.ServerError
		return
	}

	if account == nil || account.AccountType != xdr.AccountTypeNotVerified {
		action.Err = &problem.Forbidden
		return
	}
}

func (action *DeleteTFABackendsAction) performRequest() {
	err := action.APIQ().TFA().DeleteBackends(action.wallet)
	if err != nil {
		action.Log.WithError(err).Error("failed to delete backends")
		action.Err = &problem.ServerError
		return
	}
}
