package api

import (
	"gitlab.com/swarmfund/api/db2/api"
	"gitlab.com/swarmfund/api/render/hal"
	"gitlab.com/swarmfund/api/render/problem"
	"gitlab.com/swarmfund/api/resource"
	"gitlab.com/swarmfund/go/xdr"
)

type GetUserDocsAction struct {
	Action
	Address  string
	User     *api.User
	Resource resource.Documents
}

func (action *GetUserDocsAction) JSON() {
	action.Do(
		action.loadParams,
		action.checkAllowed,
		action.loadUser,
		func() {
			action.Resource.Populate(action.User)
			hal.Render(action.W, action.Resource)
		},
	)
}

func (action *GetUserDocsAction) loadParams() {
	action.Address = action.GetNonEmptyString("id")
}

func (action *GetUserDocsAction) checkAllowed() {
	action.checkSignerConstraints(
		SignerOf(action.Address),
		SignerType(action.App.CoreInfo.MasterAccountID, xdr.SignerTypeGeneralAccManager),
	)
}

func (action *GetUserDocsAction) loadUser() {
	user, err := action.APIQ().Users().ByAddress(action.Address)
	if err != nil {
		action.Log.WithError(err).Error("failed to get user")
		action.Err = &problem.ServerError
		return
	}

	if user == nil {
		action.Err = &problem.NotFound
		return
	}

	action.User = user
}
