package api

import (
	"gitlab.com/swarmfund/api/db2"
	"gitlab.com/swarmfund/api/db2/api"
	"gitlab.com/swarmfund/api/render/hal"
	"gitlab.com/swarmfund/api/resource"
)

type GetUnverifiedWalletsAction struct {
	Action

	PagingParams db2.PageQuery

	Records []api.Wallet
	Page    hal.Page
}

func (action *GetUnverifiedWalletsAction) JSON() {
	action.Do(
		action.ValidateBodyType,
		action.loadParams,
		action.checkAllowed,
		action.loadRecords,
		action.loadPage,
		func() {
			hal.Render(action.W, action.Page)
		})
}

func (action *GetUnverifiedWalletsAction) loadParams() {
	action.PagingParams = action.GetPageQuery()
}

func (action *GetUnverifiedWalletsAction) checkAllowed() {
	action.checkSignerConstraints(
		SignerOf(action.App.CoreInfo.MasterAccountID),
	)
}

func (action *GetUnverifiedWalletsAction) loadRecords() {
	//wallets, err := action.APIQ().Wallet().Page(action.PagingParams).Unverified().Select()
	//if err != nil {
	//	action.Log.WithError(err).Error("failed to get wallets")
	//	action.Err = &problem.ServerError
	//	return
	//}

	//action.Records = wallets
}

func (action *GetUnverifiedWalletsAction) loadPage() {
	for _, record := range action.Records {
		var r resource.Wallet
		ohaigo := record
		r.Populate(&ohaigo)
		action.Page.Add(r)
	}

	action.Page.BaseURL = action.BaseURL()
	action.Page.BasePath = action.Path()
	action.Page.Limit = action.PagingParams.Limit
	action.Page.Cursor = action.PagingParams.Cursor
	action.Page.Order = action.PagingParams.Order
	action.Page.PopulateLinks()
}
