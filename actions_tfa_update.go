package api

import (
	"errors"

	"fmt"

	"gitlab.com/swarmfund/api/db2/api"
	"gitlab.com/swarmfund/api/render/hal"
	"gitlab.com/swarmfund/api/render/problem"
	"gitlab.com/swarmfund/api/tfa"
	"gitlab.com/swarmfund/go/xdr"
)

type UpdateTFABackendRequest struct {
	WalletID string `json:"wallet_id" valid:"required"`
	Priority int    `json:"priority"`
}

type UpdateTFABackendAction struct {
	Action

	Request UpdateTFABackendRequest

	backendID int64
	wallet    *api.Wallet
}

func (action *UpdateTFABackendAction) JSON() {
	action.Do(
		action.ValidateBodyType,
		action.loadParams,
		action.loadWallet,
		action.checkAllowed,
		action.checkTFA,
		action.performRequest,
		func() {
			hal.Render(action.W, problem.Success)
		},
	)
}

func (action *UpdateTFABackendAction) loadParams() {
	action.UnmarshalBody(&action.Request)
	action.backendID = action.GetInt64("tfa")
}

func (action *UpdateTFABackendAction) loadWallet() {
	wallet, err := action.APIQ().Wallet().ByWalletID(action.Request.WalletID)
	if err != nil {
		action.Log.WithError(err).Error("failed to load wallet")
		action.Err = &problem.ServerError
		return
	}

	if wallet == nil {
		action.SetInvalidField("wallet_id", errors.New("does not exists"))
		return
	}

	action.wallet = wallet
}

func (action *UpdateTFABackendAction) checkAllowed() {
	action.checkSignerConstraints(
		SignerOf(action.wallet.AccountID),
	)
	if action.Err == &problem.NotAllowed {
		account, err := action.App.horizon.AccountSigned(action.App.AccountManagerKP(), action.wallet.AccountID)
		if err != nil {
			action.Log.WithError(err).Error("failed to get account")
			action.Err = &problem.ServerError
			return
		}
		if account == nil {
			return
		}
		if account.AccountType != xdr.AccountTypeNotVerified {
			return
		}
		action.Err = nil
		action.checkSignerConstraints(
			SignerType(action.App.CoreInfo.MasterAccountID, xdr.SignerTypeNotVerifiedAccManager),
		)
	}
}

func (action *UpdateTFABackendAction) checkTFA() {
	// basically is copy-paste from action.consumeTFA but forces tfa method
	record, err := action.APIQ().TFA().Backend(action.backendID)
	if err != nil {
		action.Log.WithError(err).Error("failed to get backend")
		action.Err = &problem.ServerError
		return
	}

	if record == nil {
		action.Err = &problem.NotFound
		return
	}

	//backend, err := record.Backend()
	//if err != nil {
	//	action.Log.WithError(err).WithField("backend", record.ID).Error("failed to init backend")
	//	action.Err = &problem.ServerError
	//	return
	//}

	//switch record.BackendType {
	//case api.TFABackendGoogleTOTP:
	//default:
	//	action.Err = &problem.BadRequest
	//	return
	//}

	token := tfa.Token(int64(action.wallet.Id), fmt.Sprintf("update:%d", record.ID))

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

	if otp == nil {
		//otpData, err := backend.OTPData()
		//if err != nil {
		//	action.Log.WithError(err).WithField("backend", record.ID).Error("failed to create tfa")
		//	action.Err = &problem.ServerError
		//	return
		//}
		otp = &api.TFA{
			BackendID: record.ID,
			//OTPData:   otpData,
			Token: token,
		}

		err = action.APIQ().TFA().Create(otp)
		if err != nil {
			action.Log.WithError(err).Error("failed to store tfa")
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

func (action *UpdateTFABackendAction) performRequest() {
	err := action.APIQ().TFA().SetBackendPriority(action.backendID, action.Request.Priority)
	if err != nil {
		action.Log.WithError(err).Error("failed to update backend")
		action.Err = &problem.ServerError
		return
	}
}
