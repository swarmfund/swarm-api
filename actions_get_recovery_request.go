package api

import (
	"gitlab.com/swarmfund/api/db2/api"
	"gitlab.com/swarmfund/api/render/hal"
	"gitlab.com/swarmfund/api/render/problem"
	"gitlab.com/swarmfund/api/resource"
	"gitlab.com/swarmfund/horizon-connector"
)

type GetRecoveryRequestAction struct {
	Action

	ID        int64
	Record    *api.RecoveryRequest
	Resource  resource.RecoveryRequest
	RecoverOp *horizon.RecoverOp
}

func (action *GetRecoveryRequestAction) JSON() {
	action.Do(
		action.ValidateBodyType,
		action.loadParams,
		action.loadRecord,
		action.loadRecoverOp,
		action.loadResource,
		func() {
			hal.Render(action.W, action.Resource)
		},
	)
}

func (action *GetRecoveryRequestAction) loadParams() {
	action.ID = action.GetInt64("id")
}

func (action *GetRecoveryRequestAction) loadRecord() {
	recoveryRequest, err := action.APIQ().Recoveries().ByID(action.ID)
	if err != nil {
		action.Log.WithError(err).Error("failed to get recovery request")
		action.Err = &problem.ServerError
		return
	}

	if recoveryRequest == nil {
		action.Err = &problem.NotFound
		return
	}

	action.Record = recoveryRequest
}

func (action *GetRecoveryRequestAction) loadRecoverOp() {
	if action.Record.RecoveryWalletID == nil {
		return
	}

	initialWallet, err := action.APIQ().Wallet().ByID(action.Record.WalletID)
	if err != nil {
		action.Log.WithError(err).Error("failed to load initial wallet")
		action.Err = &problem.ServerError
		return
	}

	if initialWallet == nil {
		action.Log.WithField("recovery_request", action.Record.ID).Error("wallet expected to exist")
		action.Err = &problem.ServerError
		return
	}

	recoveryWallet, err := action.APIQ().Wallet().RecoveryWallet(*action.Record.RecoveryWalletID, initialWallet.Username)
	if err != nil {
		action.Log.WithError(err).Error("failed to load recovery wallet")
		action.Err = &problem.ServerError
		return
	}

	if recoveryWallet == nil {
		action.Log.WithField("recovery_request", action.Record.ID).Error("recovery wallet expected to exist")
		action.Err = &problem.ServerError
		return
	}

	action.RecoverOp = &horizon.RecoverOp{
		AccountID: initialWallet.AccountID,
		OldSigner: initialWallet.CurrentAccountID,
		NewSigner: recoveryWallet.CurrentAccountID,
	}
}

func (action *GetRecoveryRequestAction) loadResource() {
	action.Resource.Populate(action.Record, action.RecoverOp)
}
