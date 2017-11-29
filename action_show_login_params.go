package api

import (
	"gitlab.com/swarmfund/api/render/hal"
	"gitlab.com/swarmfund/api/render/problem"
	"gitlab.com/swarmfund/api/resource"
)

type ShowLoginParamsAction struct {
	Action
	username string
	response resource.LoginParams
}

func (action *ShowLoginParamsAction) JSON() {
	action.ValidateBodyType()
	action.Do(
		action.loadParams,
		action.performRequest,
		func() {
			hal.Render(action.W, action.response)
		})
}

func (action *ShowLoginParamsAction) loadParams() {
	action.ValidateBodyType()
	action.username = action.GetNonEmptyString("username")
}

func (action *ShowLoginParamsAction) performRequest() {
	wallet, err := action.APIQ().Wallet().ByEmail(action.username)

	if err != nil {
		action.Log.WithError(err).Error("Failed to get wallet")
		action.Err = &problem.ServerError
		return
	}
	if wallet == nil {
		action.Err = &problem.NotFound
		return
	}

	action.response.Populate(wallet)
}
