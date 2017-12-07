package api

import (
	"gitlab.com/swarmfund/api/db2/api"
	"gitlab.com/swarmfund/api/render/hal"
	"gitlab.com/swarmfund/api/render/problem"
)

// GetUserIdAction render a user's account_id by email.
type GetUserIdAction struct {
	Action
	Email  string
	Record *api.User
}

// JSON is a method for actions.JSON
func (action *GetUserIdAction) JSON() {
	action.Do(
		action.ValidateBodyType,
		action.loadParams,
		action.loadRecord,
		func() {
			res := map[string]string{
				"account_id": string(action.Record.Address),
			}
			hal.Render(action.W, res)
		})
}

func (action *GetUserIdAction) loadParams() {
	action.Email = action.GetNonEmptyString("email")
}

func (action *GetUserIdAction) loadRecord() {
	if action.Err != nil {
		return
	}

	action.Record, action.Err = action.APIQ().Users().ByEmail(action.Email)
	if action.Err != nil {
		action.Log.WithError(action.Err).Error("Failed to get user from db")
		action.Err = &problem.ServerError
		return
	}

	if action.Record == nil {
		action.Err = &problem.NotFound
		return
	}
}
