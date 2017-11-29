package api

import (
	"fmt"

	"gitlab.com/swarmfund/api/db2/api"
	"gitlab.com/swarmfund/api/render/hal"
	"gitlab.com/swarmfund/api/render/problem"
)

type CreateRecoveryRequestAction struct {
	Action

	Username        string
	RecoveryRequest *api.RecoveryRequest
}

func (action *CreateRecoveryRequestAction) JSON() {
	action.Do(
		action.ValidateBodyType,
		action.loadParams,
		action.performRequest,
		action.sendNotification,
		func() {
			hal.Render(action.W, problem.Success)
		},
	)
}

func (action *CreateRecoveryRequestAction) loadParams() {
	action.Username = action.GetNonEmptyString("username")
}

func (action *CreateRecoveryRequestAction) performRequest() {
	wallet, err := action.APIQ().Wallet().ByEmail(action.Username)
	if err != nil {
		action.Log.WithError(err).Error("failed to get wallet")
		action.Err = &problem.ServerError
		return
	}

	if wallet == nil {
		// since endpoint is public caller shouldn't know about wallet existence
		action.Err = &problem.Success
		return
	}

	if !wallet.Verified {
		// user should verify email before proceeding with recovery
		action.Err = &problem.Forbidden
		return
	}

	// check if there is recovery in progress
	action.RecoveryRequest, err = action.APIQ().Recoveries().ByAccountID(wallet.AccountID)
	if err != nil {
		action.Log.WithError(err).Error("failed to check recovery requests")
		action.Err = &problem.ServerError
		return
	}

	if action.RecoveryRequest != nil {
		// not setting any errors, so user will receive notification with same token again
		// rate limits should be enforced on notificator config
		return
	}

	action.RecoveryRequest, err = api.NewRecoveryRequest(wallet)
	if err != nil {
		action.Log.WithError(err).Error("failed to generate recovery request")
		action.Err = &problem.ServerError
		return
	}

	err = action.APIQ().Recoveries().Create(action.RecoveryRequest)
	if err != nil {
		action.Log.WithError(err).Error("failed to save recovery request")
		action.Err = &problem.ServerError
		return
	}
}

func (action *CreateRecoveryRequestAction) sendNotification() {
	response, err := action.Notificator().SendRecoveryRequest(action.Username, action.RecoveryRequest.EmailToken, "")
	if err != nil {
		action.Log.WithError(err).Error("failed to connect notificator")
		// we could try to clean up recovery request, but it might already be sent previously
		// instead we are failing, so when user retries request we will try again
		action.Err = &problem.ServerError
		return
	}

	if !response.IsSuccess() {
		retryIn := response.RetryIn()
		if retryIn != nil {
			p := problem.RecoveryRequestLimitExceeded
			p.Detail = fmt.Sprintf("Retry in %s", retryIn.String())
			action.Log.WithError(err).Warn("recovery notification to many requests")
			action.Err = &p
			return
		}
		action.Log.WithError(err).Error("recovery notification not sent")
		action.Err = &problem.ServerError
		return
	}

	err = action.APIQ().Recoveries().MarkSent(action.RecoveryRequest.EmailToken)
	if err != nil {
		// since email is sent we shouldn't fail just log
		action.Log.WithError(err).Error("failed to mark recovery as sent")
	}
}
