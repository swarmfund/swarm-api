package api

import (
	"gitlab.com/swarmfund/api/db2/api"
	"gitlab.com/swarmfund/api/render/hal"
	"gitlab.com/swarmfund/api/render/problem"
	"gitlab.com/swarmfund/api/resource"
)

type DetailsRequest struct {
	Addresses []string `json:"addresses"`
}

type DetailsAction struct {
	Action

	Request DetailsRequest

	Records  []api.User
	Resource resource.ShortenUsersDetails
}

func (action *DetailsAction) JSON() {
	action.Do(
		action.ValidateBodyType,
		action.loadParams,
		action.loadRecords,
		action.loadResource,
		func() {
			hal.Render(action.W, action.Resource)
		},
	)
}

func (action *DetailsAction) loadParams() {
	action.UnmarshalBody(&action.Request)
}

func (action *DetailsAction) loadRecords() {
	users, err := action.APIQ().Users().ByAddresses(action.Request.Addresses)
	if err != nil {
		action.Log.WithError(err).Error("failed to get users")
		action.Err = &problem.ServerError
		return
	}
	action.Records = users
}

func (action *DetailsAction) loadResource() {
	action.Resource.Populate(action.Records)
}
