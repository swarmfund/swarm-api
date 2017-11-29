package api

import (
	"gitlab.com/swarmfund/api/db2"
	"gitlab.com/swarmfund/api/db2/api"
	"gitlab.com/swarmfund/api/render/hal"
	"gitlab.com/swarmfund/api/render/problem"
)

type GetRecoveryRequestsAction struct {
	Action
	Records      []api.RecoveryRequest
	PagingParams db2.PageQuery
	Page         hal.Page
}

func (action *GetRecoveryRequestsAction) JSON() {
	action.Do(
		action.loadParams,
		action.checkAllowed,
		action.loadResource,
		action.loadPage,
		func() {
			hal.Render(action.W, action.Page)
		},
	)
}

func (action *GetRecoveryRequestsAction) loadParams() {
	action.PagingParams = action.GetPageQuery()
}

func (action *GetRecoveryRequestsAction) checkAllowed() {
	action.checkSignerConstraints(
		SignerOf(action.App.CoreInfo.MasterAccountID),
	)
}

func (action *GetRecoveryRequestsAction) loadResource() {
	var err error
	action.Records, err = action.APIQ().Recoveries().Page(action.PagingParams).Uploaded().Select()
	if err != nil {
		action.Log.WithError(err).Error("failed to get recovery requests")
		action.Err = &problem.ServerError
		return
	}
}

func (action *GetRecoveryRequestsAction) loadPage() {
	for _, record := range action.Records {
		action.Page.Add(record)
	}

	action.Page.BaseURL = action.BaseURL()
	action.Page.BasePath = action.Path()
	action.Page.Limit = action.PagingParams.Limit
	action.Page.Cursor = action.PagingParams.Cursor
	action.Page.Order = action.PagingParams.Order
	action.Page.PopulateLinks()
}
