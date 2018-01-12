package api

import (
	"net/http"

	"gitlab.com/swarmfund/api/db2/api"
	"gitlab.com/swarmfund/api/render/hal"
	"gitlab.com/swarmfund/api/render/problem"
)

type GetUserFileAction struct {
	Action

	Address string
	Version string

	User     *api.User
	Response map[string]string
}

func (action *GetUserFileAction) JSON() {
	action.Do(
		action.checkAvailable,
		action.loadParams,
		action.checkAllowed,
		action.loadUser,
		action.performRequest,
		func() {
			hal.Render(action.W, action.Response)
		},
	)
}

func (action *GetUserFileAction) checkAvailable() {
	if action.App.Config().Storage().Disabled {
		action.Log.Warn("storage service disabled")
		action.Err = &problem.P{
			Status: http.StatusServiceUnavailable,
		}
		return
	}
}

func (action *GetUserFileAction) loadParams() {
	action.Address = action.GetNonEmptyString("id")
	action.Version = action.GetNonEmptyString("version")
}

func (action *GetUserFileAction) checkAllowed() {
	action.checkSignerConstraints(
		SignerOf(action.Address),
		SignerOf(action.App.CoreInfo.MasterAccountID),
	)
}

func (action *GetUserFileAction) loadUser() {
	action.User, action.Err = action.APIQ().Users().ByAddress(action.Address)
	if action.Err != nil {
		action.Log.WithError(action.Err).Error("Failed to load user")
		action.Err = &problem.ServerError
		return
	}

	if action.User == nil {
		action.Err = &problem.NotFound
		return
	}
}

func (action *GetUserFileAction) performRequest() {
	document := action.User.Documents.Get(func(doc *api.Document) bool {
		return doc.Version == action.Version
	})

	if document == nil {
		action.Err = &problem.NotFound
		return
	}
}
