package api

import (
	"gitlab.com/swarmfund/api/db2/api"
	"gitlab.com/swarmfund/api/render/hal"
	"gitlab.com/swarmfund/api/render/problem"
	"gitlab.com/swarmfund/api/resource"
)

type GetWalletOrganizationAction struct {
	Action

	WalletID int64

	Wallet *api.Wallet
	User   *api.User

	Resource resource.User
}

func (action *GetWalletOrganizationAction) JSON() {
	action.Do(
		action.loadParams,
		action.loadWallet,
		action.checkAllowed,
		action.loadRecord,
		action.loadResource,
		func() {
			hal.Render(action.W, action.Resource)
		},
	)
}

func (action *GetWalletOrganizationAction) loadParams() {
	action.WalletID = action.GetInt64("id")
}

func (action *GetWalletOrganizationAction) loadWallet() {
	wallet, err := action.APIQ().Wallet().ByID(action.WalletID)
	if err != nil {
		action.Log.WithError(err).Error("failed to get wallet")
		action.Err = &problem.ServerError
		return
	}

	if wallet == nil {
		action.Err = &problem.NotFound
		return
	}

	action.Wallet = wallet
}

func (action *GetWalletOrganizationAction) checkAllowed() {
	action.checkSignerConstraints(
		SignedBy(action.Wallet.CurrentAccountID),
	)
}

func (action *GetWalletOrganizationAction) loadRecord() {
	if action.Wallet.OrganizationAddress == nil {
		// wallet has not been marked for multi-sign flow or has been connected yet
		action.Err = &problem.NotFound
		return
	}

	user, err := action.APIQ().Users().ByAddress(*action.Wallet.OrganizationAddress)
	if err != nil {
		action.Log.WithError(err).Error("failed to get user")
		action.Err = &problem.ServerError
		return
	}

	if user == nil {
		action.Log.WithField("address", *action.Wallet.OrganizationAddress).Error("user expected to exist")
		action.Err = &problem.ServerError
		return
	}

	action.User = user
}

func (action *GetWalletOrganizationAction) loadResource() {
	action.Resource.Populate(action.User)
}
