package api

import (
	"gitlab.com/swarmfund/api/db2/api"
	"gitlab.com/swarmfund/api/render/hal"
	"gitlab.com/swarmfund/api/render/problem"
	"gitlab.com/swarmfund/api/resource/base"
)

type ParticipantsRequest struct {
	ForAccount   string                      `json:"for_account" valid:"required"`
	Participants map[int64][]api.Participant `json:"participants"`
}

type ParticipantsAction struct {
	Action

	Request  ParticipantsRequest
	Resource map[int64][]base.Participant
}

func (action *ParticipantsAction) JSON() {
	action.Do(
		action.ValidateBodyType,
		action.loadParams,
		action.checkAllowed,
		action.loadParticipants,
		action.loadResource,
		func() {
			hal.Render(action.W, action.Resource)
		},
	)
}

func (action *ParticipantsAction) loadParams() {
	action.UnmarshalBody(&action.Request)
}

func (action *ParticipantsAction) checkAllowed() {
	// TODO master
}

func (action *ParticipantsAction) loadParticipants() {
	err := action.APIQ().Users().Participants(action.Request.Participants)
	if err != nil {
		action.Log.WithError(err).Error("failed to get participant details")
		action.Err = &problem.ServerError
		return
	}
}

func (action *ParticipantsAction) loadResource() {
	action.Resource = map[int64][]base.Participant{}
	for op, participants := range action.Request.Participants {
		for _, participant := range participants {
			var r base.Participant
			r.Populate(&participant)
			action.Resource[op] = append(action.Resource[op], r)
		}
	}
}
