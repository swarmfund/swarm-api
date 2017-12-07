package api

import (
	"gitlab.com/swarmfund/api/db2/api"
	"gitlab.com/swarmfund/api/render/redirect"
	"gitlab.com/swarmfund/go/xdr"
	"gitlab.com/swarmfund/horizon-connector"
)

type ApproveRecoveryRequestAction struct {
	Action

	Token           string
	RedirectPayload *redirect.Payload
	RecoveryRequest *api.RecoveryRequest
	Wallet          *api.Wallet
}

func (action *ApproveRecoveryRequestAction) JSON() {
	// pipeline shouldn't set action.Err anywhere
	// use `RedirectPayload` instead
	action.Do(
		action.checkAvailable,
		action.loadParams,
		action.loadRecoveryRequest,
		action.loadWallet,
		action.checkUploaded,
		action.blockAccount,
		action.setUserState,
		action.markRequest,
		action.craftRedirect,
		func() {
			// TODO
			//redirectURL := *action.App.config.ClientRouter
			//encodedPayload, err := action.RedirectPayload.Encode()
			//if err != nil {
			//	// client should handle invalid or empty payloads as 500
			//	action.Log.WithError(err).Error("failed to encode payload")
			//} else {
			//	query := redirectURL.Query()
			//	query.Set("action", encodedPayload)
			//	redirectURL.RawQuery = query.Encode()
			//
			//}
			//hal.Redirect(action.W, action.R, redirectURL.String())
		},
	)
}

func (action *ApproveRecoveryRequestAction) checkAvailable() {
	if action.App.Config().Storage().Disabled {
		action.Log.Warn("storage service disabled")
		action.RedirectPayload = &redirect.Unavailable
		return
	}
}

func (action *ApproveRecoveryRequestAction) loadParams() {
	action.Token = action.GetString("token")
}

func (action *ApproveRecoveryRequestAction) loadRecoveryRequest() {
	if action.RedirectPayload != nil {
		return
	}

	recoveryRequest, err := action.APIQ().Recoveries().Get(action.Token)
	if err != nil {
		action.Log.WithError(err).Error("failed to get request")
		action.RedirectPayload = &redirect.ServerError
		return
	}

	if recoveryRequest == nil {
		action.RedirectPayload = &redirect.NotFound
		return
	}

	action.RecoveryRequest = recoveryRequest
}

func (action *ApproveRecoveryRequestAction) loadWallet() {
	if action.RedirectPayload != nil {
		return
	}

	wallet, err := action.APIQ().Wallet().ByID(action.RecoveryRequest.WalletID)
	if err != nil {
		action.Log.WithError(err).Error("failed to get wallet")
		action.RedirectPayload = &redirect.ServerError
		return
	}

	if wallet == nil {
		// shouldn't really happen, probably something went wrong
		action.Log.WithField("recovery_request", action.RecoveryRequest.ID).Error("failed to get wallet")
		action.RedirectPayload = &redirect.ServerError
		return
	}

	action.Wallet = wallet
}

func (action *ApproveRecoveryRequestAction) checkUploaded() {
	if action.RedirectPayload != nil {
		return
	}

	// TODO implement
	//exists, err := action.App.Storage().Exists(action.RecoveryRequest.AccountID, api.DocumentTypeRecoveryPhoto)
	//if err != nil {
	//	action.Log.WithError(err).Error("riak failed")
	//	action.RedirectPayload = &redirect.ServerError
	//	return
	//}
	//
	//if exists {
	//	action.RedirectPayload = &redirect.RecoveryRequestAlreadyUploaded
	//	return
	//}
}

func (action *ApproveRecoveryRequestAction) blockAccount() {
	if action.RedirectPayload != nil {
		return
	}

	coreAccount, err := action.App.horizon.AccountSigned(action.App.AccountManagerKP(), action.RecoveryRequest.AccountID)
	if err != nil {
		action.Log.WithError(err).Error("failed to get core account")
		action.RedirectPayload = &redirect.ServerError
		return
	}

	if coreAccount == nil {
		// if user is not present in core db we are "approving" request right away
		// by purging everything related to username and letting him sign up again

		err = action.PurgeUser(action.RecoveryRequest.AccountID)
		if err != nil {
			action.Log.WithError(err).Error("failed to purge user")
			action.RedirectPayload = &redirect.ServerError
			return
		}
		action.RedirectPayload = &redirect.SignupPayload
		return
	}

	if xdr.BlockReasons(coreAccount.BlockReasons)&xdr.BlockReasonsRecoveryRequest == 0 {
		err := action.App.horizon.Transaction(&horizon.TransactionBuilder{Source: action.App.MasterKP()}).
			Op(&horizon.ManageAccountOp{
				AccountID:   action.RecoveryRequest.AccountID,
				AccountType: xdr.AccountType(coreAccount.AccountType),
				AddReasons:  xdr.BlockReasonsRecoveryRequest,
			}).Sign(action.App.AccountManagerKP()).Submit()

		if err != nil {
			action.Log.WithError(err).Error("failed to submit tx")
			action.RedirectPayload = &redirect.ServerError
			return
		}
	}
}

func (action *ApproveRecoveryRequestAction) setUserState() {
	if action.RedirectPayload != nil {
		return
	}

	err := action.APIQ().Users().SetRecoveryState(action.RecoveryRequest.AccountID, api.UserRecoveryStatePending)
	if err != nil {
		action.Log.WithError(err).Error("failed to set user recovery state")
		action.RedirectPayload = &redirect.ServerError
		return
	}
}

func (action *ApproveRecoveryRequestAction) markRequest() {
	if action.RedirectPayload != nil {
		return
	}

	err := action.APIQ().Recoveries().MarkCodeShown(action.Token)
	if err != nil {
		action.Log.WithError(err).Error("failed to change request status")
		action.RedirectPayload = &redirect.ServerError
		return
	}
}

func (action *ApproveRecoveryRequestAction) craftRedirect() {
	if action.RedirectPayload != nil {
		return
	}

	account, err := action.App.horizon.AccountSigned(action.App.AccountManagerKP(), action.RecoveryRequest.AccountID)
	if err != nil {
		action.Log.WithError(err).Error("failed to get core account")
		action.RedirectPayload = &redirect.ServerError
		return
	}

	if account == nil {
		action.Log.WithField("address", action.RecoveryRequest.AccountID).Error("account expected to exist")
		action.RedirectPayload = &redirect.ServerError
		return
	}

	switch xdr.AccountType(account.AccountType) {
	case xdr.AccountTypeNotVerified:
		action.RedirectPayload = redirect.RecoveryRequestShowCode(action.RecoveryRequest.Username, action.RecoveryRequest.Code, true)
	case xdr.AccountTypeGeneral:
		action.RedirectPayload = redirect.RecoveryRequestShowCode(action.RecoveryRequest.Username, action.RecoveryRequest.Code, false)
	}
}
