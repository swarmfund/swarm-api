package api

import (
	"gitlab.com/swarmfund/api/db2/api"
	"gitlab.com/swarmfund/api/render/hal"
	"gitlab.com/swarmfund/api/render/problem"
	"github.com/go-errors/errors"
)

type ResolveUserRecoveryRequestRequest struct {
	TX           string `json:"tx"`
	Approved     bool   `json:"approved"`
	WalletID     string `json:"wallet_id"`
	RejectReason string `json:"reject_reason"`
}

type ResolveUserRecoveryRequestAction struct {
	Action

	AccountID       string
	InitialWallet   *api.Wallet
	Wallet          *api.Wallet
	RecoveryRequest *api.RecoveryRequest

	Request ResolveUserRecoveryRequestRequest
}

func (action *ResolveUserRecoveryRequestAction) JSON() {
	action.Do(
		action.ValidateBodyType,
		action.loadParams,
		action.checkAllowed,
		action.loadRecoveryRequest,
		action.loadInitialWallet,
		action.loadWallet,
		action.performRequest,
		func() {
			hal.Render(action.W, problem.Success)
		},
	)
}

func (action *ResolveUserRecoveryRequestAction) loadParams() {
	action.AccountID = action.GetNonEmptyString("id")
	action.UnmarshalBody(&action.Request)
	if action.Request.Approved {
		if action.Request.WalletID == "" {
			action.SetInvalidField("wallet_id", errors.New("required"))
			return
		}
		if action.Request.TX == "" {
			action.SetInvalidField("tx", errors.New("required"))
			return
		}
	} else {
		if action.Request.RejectReason == "" {
			action.SetInvalidField("reject_reason", errors.New("required"))
			return
		}
	}

}

func (action *ResolveUserRecoveryRequestAction) checkAllowed() {
	action.checkSignerConstraints(
		SignerOf(action.App.CoreInfo.MasterAccountID),
	)
}

func (action *ResolveUserRecoveryRequestAction) loadRecoveryRequest() {
	recoveryRequest, err := action.APIQ().Recoveries().ByAccountID(action.AccountID)
	if err != nil {
		action.Log.WithError(err).Error("failed to get recovery request")
		action.Err = &problem.ServerError
		return
	}

	if recoveryRequest == nil {
		action.Err = &problem.NotFound
		return
	}
	action.RecoveryRequest = recoveryRequest
}

func (action *ResolveUserRecoveryRequestAction) loadInitialWallet() {
	wallet, err := action.APIQ().Wallet().ByID(action.RecoveryRequest.WalletID)
	if err != nil {
		action.Log.WithError(err).Error("failed to get initial wallet")
		action.Err = &problem.ServerError
		return
	}

	if wallet == nil {
		// something went wrong
		action.Log.WithField("recovery_request", action.RecoveryRequest.ID).
			Error("initial wallet not found")
		action.Err = &problem.ServerError
		return
	}

	action.InitialWallet = wallet
}

func (action *ResolveUserRecoveryRequestAction) loadWallet() {
	if action.Request.WalletID == "" {
		return
	}
	wallet, err := action.APIQ().Wallet().RecoveryWallet(action.Request.WalletID, action.RecoveryRequest.Username)
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

func (action *ResolveUserRecoveryRequestAction) performRequest() {
	var err error
	if action.Request.Approved {
		err = action.recoveryApprove(action.RecoveryRequest, action.InitialWallet, action.Wallet, action.Request.TX)
		if err != nil {
			action.Log.WithError(err).Error("failed to approve recovery request")
			action.Err = &problem.ServerError
			return
		}
	} else {
		// recovery request was rejected

		// flip upload state
		err = action.APIQ().Recoveries().MarkRejected(action.AccountID)
		if err != nil {
			action.Log.WithError(err).Error("failed to make recovery request rejected")
			action.Err = &problem.ServerError
			return
		}

		// resend recovery request notification
		// TODO make it separate notification type with own template
		_, err = action.Notificator().SendRecoveryRequest(
			action.RecoveryRequest.Username, action.RecoveryRequest.EmailToken, action.Request.RejectReason)
		if err != nil {
			action.Log.WithError(err).Error("failed to send recovery notification")
			action.Err = &problem.ServerError
			return
		}
	}

	// TODO implement
	//err = action.App.Storage().Delete(action.RecoveryRequest.AccountID, api.DocumentTypeRecoveryPhoto)
	//if err != nil {
	//	action.Log.WithError(err).Error("failed to delete doc")
	//	action.Err = &problem.ServerError
	//	return
	//}
}
