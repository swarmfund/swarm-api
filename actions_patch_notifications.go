package api

import (
	"gitlab.com/swarmfund/api/db2/api"
	"gitlab.com/swarmfund/api/render/hal"
	"gitlab.com/swarmfund/api/render/problem"
)

type PatchNotificationsRequest struct {
	Email     *string `json:"email"`
	KYCReview *bool   `json:"kyc_review"`
}

type PatchNotificationsAction struct {
	Action

	AccountID string
	Request   PatchNotificationsRequest
}

func (action *PatchNotificationsAction) JSON() {
	action.Do(
		action.ValidateBodyType,
		action.loadParams,
		action.checkAllowed,
		action.performRequest,
		func() {
			hal.Render(action.W, problem.Success)
		},
	)
}

func (action *PatchNotificationsAction) loadParams() {
	action.UnmarshalBody(&action.Request)
	action.AccountID = action.GetNonEmptyString("id")
}

func (action *PatchNotificationsAction) checkAllowed() {
	action.checkSignerConstraints(
		SignerOf(action.App.CoreInfo.MasterAccountID),
	)
}

func (action *PatchNotificationsAction) performRequest() {
	if action.Request.KYCReview != nil {
		var err error
		if *action.Request.KYCReview {
			err = action.APIQ().Notifications().Enable(action.AccountID, api.NotificationTypeKYC)
		} else {
			err = action.APIQ().Notifications().Disable(action.AccountID, api.NotificationTypeKYC)
		}
		if err != nil {
			action.Log.WithError(err).Error("failed to update kyc notification")
			action.Err = &problem.ServerError
			return
		}
	}

	if action.Request.Email != nil {
		if err := action.APIQ().Notifications().SetEmail(action.AccountID, *action.Request.Email); err != nil {
			action.Log.WithError(err).Error(err)
			action.Err = &problem.ServerError
			return
		}
	}
}
