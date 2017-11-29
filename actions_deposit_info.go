package api

import (
	"gitlab.com/swarmfund/api/render/hal"
	"gitlab.com/swarmfund/api/render/problem"
)

type DepositInfoAction struct {
	Action

	Method   DepositMethod
	Resource interface{}
}

func (action *DepositInfoAction) JSON() {
	action.Do(
		action.loadParams,
		action.loadResource,
		func() {
			hal.Render(action.W, action.Resource)
		},
	)
}

func (action *DepositInfoAction) loadParams() {
	action.Method = DepositMethod(action.GetInt32("method"))
}

func (action *DepositInfoAction) loadResource() {
	switch action.Method {
	case DepositMethodStripe:
		// TODO
		//action.Resource = map[string]string{
		//	"pk": action.App.config.Deposit.StripePK,
		//}
	case DepositMethodSkrill:
		action.Resource = map[string]string{}
	default:
		action.Err = &problem.NotFound
		return
	}
}
