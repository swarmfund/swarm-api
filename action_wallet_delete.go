package api

import (
	"errors"

	"gitlab.com/swarmfund/api/render/hal"
	"gitlab.com/swarmfund/api/render/problem"
)

type DeleteWalletAction struct {
	Action
	Username string
}

func (action *DeleteWalletAction) JSON() {
	action.ValidateBodyType()
	action.Do(
		action.loadParams,
		action.checkIsAllowed,
		action.performRequest,
		func() {
			hal.Render(action.W, problem.Success)
		})
}

func (action *DeleteWalletAction) loadParams() {
	action.ValidateBodyType()
	action.Username = action.GetNonEmptyString("username")
}

func (action *DeleteWalletAction) checkIsAllowed() {
	action.checkSignerConstraints(
		SignerOf(action.App.CoreInfo.MasterAccountID),
	)
}

func (action *DeleteWalletAction) performRequest() {
	wallet, err := action.APIQ().Wallet().ByEmail(action.Username)
	if err != nil {
		action.Log.WithError(err).Error("Failed when try to get wallet")
		action.Err = &problem.ServerError
		return
	}

	if wallet == nil {
		action.Err = &problem.NotFound
		return
	}

	if wallet.Verified {
		action.SetInvalidField("username", errors.New("this user is already verified"))
		return
	}

	err = action.APIQ().Wallet().Delete(wallet.Id)
	if err != nil {
		action.Log.WithError(err).Error("Failed when try to delete wallet")
		action.Err = &problem.ServerError
		return
	}

}
