package handlers

import (
	"encoding/json"
	"net/http"

	. "github.com/go-ozzo/ozzo-validation"
	"github.com/pkg/errors"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"gitlab.com/swarmfund/api/db2/api"
	"gitlab.com/swarmfund/api/internal/types"
	"gitlab.com/swarmfund/api/resource/base"
)

type (
	GetParticipantsRequest struct {
		ForAccount   string                      `json:"for_account"`
		Participants map[int64][]api.Participant `json:"participants"`
	}
)

func NewGetParticipantsRequest(r *http.Request) (GetParticipantsRequest, error) {
	var request GetParticipantsRequest

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return request, errors.Wrap(err, "failed to unmarshal")
	}

	return request, request.Validate()
}

func (r GetParticipantsRequest) Validate() error {
	return ValidateStruct(&r,
		Field(&r.ForAccount, Required),
		Field(&r.Participants, Required),
	)
}

func GetParticipants(w http.ResponseWriter, r *http.Request) {

	request, err := NewGetParticipantsRequest(r)

	if err != nil {
		ape.RenderErr(w, problems.BadRequest(err)...)
	}

	var accountIDs []string

	for _, op := range request.Participants {
		for pi := range op {
			accountIDs = append(accountIDs, string(op[pi].AccountID))
		}
	}

	users, err := UsersQ(r).ByAddresses(accountIDs)
	if err != nil {
		Log(r).WithError(err).Error("Can't find users")
		ape.RenderErr(w, problems.InternalError())
	}

	usersMap := map[types.Address]api.User{}
	for _, user := range users {
		usersMap[user.Address] = user
	}

	for _, op := range request.Participants {
		for pi := range op {
			participant := op[pi]
			if user, ok := usersMap[participant.AccountID]; ok {
				participant.Email = &user.Email
				op[pi] = participant
			}
		}
	}

	response := make(map[int64][]base.Participant)
	for op, participants := range request.Participants {
		for _, participant := range participants {
			var r base.Participant
			r.Populate(&participant)
			response[op] = append(response[op], r)
		}
	}

	json.NewEncoder(w).Encode(response)
}
