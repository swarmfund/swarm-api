package api

import (
	"gitlab.com/swarmfund/api/render/hal"
	"gitlab.com/swarmfund/api/render/problem"
)

type ResendTokenAction struct {
	Action
	username string
}

func (action *ResendTokenAction) JSON() {
	action.Do(
		action.ValidateBodyType,
		action.loadParams,
		action.performRequest,
		func() {
			hal.Render(action.W, problem.Success)
		},
	)
}

func (action *ResendTokenAction) loadParams() {
	action.username = action.GetNonEmptyString("username")
}

func (action *ResendTokenAction) performRequest() {
	wallet, err := action.APIQ().Wallet().ByEmail(action.username)
	if err != nil {
		action.Log.WithError(err).Error("Failed to get wallet by username")
		action.Err = &problem.ServerError
		return
	}

	if wallet == nil || wallet.Verified {
		// we should not leak any info about users
		return
	}

	//err = action.Notificator().SendVerificationLink(wallet.Username, wallet.VerificationToken)
	//if err != nil {
	//	action.Log.WithError(err).Error("Unable to send email")
	//	action.Err = &problem.ServerError
	//	return
	//}
}
