package api

import (
	"github.com/pkg/errors"
	"gitlab.com/swarmfund/api/db2"
	"gitlab.com/swarmfund/api/db2/api"
	"gitlab.com/swarmfund/api/render/hal"
	"gitlab.com/swarmfund/api/render/problem"
	"gitlab.com/swarmfund/api/resource"
)

// This file contains the actions:
//
// UserIndexAction: pages of users
// UserShowAction: single user by accountId
// GetUserIdAction: account_id by email.

// UserIndexAction renders a page of operations resources, identified by
// a normal page query and optionally filtered by an account, ledger, or
// transaction.
type UserIndexAction struct {
	Action
	Status       string
	PagingParams db2.PageQuery
	Records      []api.User
	Page         hal.Page
}

// JSON is a method for actions.JSON
func (action *UserIndexAction) JSON() {
	action.Do(
		action.ValidateCursorAsDefault,
		action.loadParams,
		action.checkAllowed,
		action.loadRecords,
		action.loadPage,
		func() {
			hal.Render(action.W, action.Page)
		},
	)
}

func (action *UserIndexAction) loadParams() {
	action.Status = action.GetString("status")
	action.PagingParams = action.GetPageQuery()
	action.Page.Filters = map[string]string{
		"status": action.Status,
	}
}

func (action *UserIndexAction) loadRecords() {
	u := action.APIQ().Users()

	if api.IsUserState(action.Status) {
		u.ByState(api.UserState(action.Status))
	} else {
		switch action.Status {
		case "poi_review":
			u.LimitReviewRequests()
		case "recovery":
			u.RecoveryPending()
		default:
			action.SetInvalidField("status", errors.New("invalid"))
		}
	}

	action.Err = u.Page(action.PagingParams).Select(&action.Records)
}

func (action *UserIndexAction) loadPage() {
	for _, record := range action.Records {
		var res resource.User
		res.Populate(&record)
		action.Page.Add(res)
	}

	action.Page.BaseURL = action.BaseURL()
	action.Page.BasePath = action.Path()
	action.Page.Limit = action.PagingParams.Limit
	action.Page.Cursor = action.PagingParams.Cursor
	action.Page.Order = action.PagingParams.Order
	action.Page.PopulateLinks()
}

func (action *UserIndexAction) checkAllowed() {
	action.checkSignerConstraints(
		SignerOf(action.App.CoreInfo.MasterAccountID),
	)
}

// UserShowAction renders a ledger found by its sequence number.
type UserShowAction struct {
	Action
	AccountId string
	Record    *api.User
}

// JSON is a method for actions.JSON
func (action *UserShowAction) JSON() {
	action.Do(
		action.loadParams,
		action.checkAllowed,
		action.loadRecord,
		func() {
			var res resource.User
			res.Populate(action.Record)
			hal.Render(action.W, res)
		})
}
func (action *UserShowAction) loadParams() {
	action.AccountId = action.GetNonEmptyString("id")
}

func (action *UserShowAction) checkAllowed() {
	action.checkSignerConstraints(
		SignerOf(action.AccountId),
		SignerOf(action.App.CoreInfo.MasterAccountID),
	)
}

func (action *UserShowAction) loadRecord() {
	action.Record, action.Err = action.APIQ().Users().ByAddress(action.AccountId)
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
