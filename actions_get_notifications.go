package api

import (
	"gitlab.com/swarmfund/api/db2/api"
	"gitlab.com/swarmfund/api/render/hal"
	"gitlab.com/swarmfund/api/render/problem"
	"gitlab.com/swarmfund/api/resource"
)

type GetNotificationsAction struct {
	Action

	AccountID string

	Records *api.Notifications

	Resource resource.Notifications
}

func (action *GetNotificationsAction) JSON() {
	action.Do(
		action.loadParams,
		action.checkAllowed,
		action.loadRecords,
		action.populateResource,
		func() {
			hal.Render(action.W, action.Resource)
		},
	)
}

func (action *GetNotificationsAction) loadParams() {
	action.AccountID = action.GetNonEmptyString("id")
}

func (action *GetNotificationsAction) checkAllowed() {
	action.checkSignerConstraints(
		SignerOf(action.App.CoreInfo.MasterAccountID),
	)
}

func (action *GetNotificationsAction) loadRecords() {
	records, err := action.APIQ().Notifications().Get(action.AccountID)
	if err != nil {
		action.Log.WithError(err).Error("failed to get notifications")
		action.Err = &problem.ServerError
		return
	}
	action.Records = records
}

func (action *GetNotificationsAction) populateResource() {
	action.Resource.Populate(action.Records)
}
