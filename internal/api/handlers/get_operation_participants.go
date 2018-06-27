package handlers

import (
	"encoding/json"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"gitlab.com/swarmfund/api/db2/api"
	"gitlab.com/swarmfund/api/internal/types"
	"gitlab.com/swarmfund/api/resource/base"
	"io/ioutil"
	"net/http"
)

type (
	ParticipantsRequest struct {
		ForAccount   string                      `json:"for_account" valid:"required"`
		Participants map[int64][]api.Participant `json:"participants"`
	}
	Participants struct {
		Request  ParticipantsRequest
		Resource map[int64][]base.Participant
	}
)

func GetParticipants(w http.ResponseWriter, r *http.Request) {

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		Log(r).WithError(err).Error("Incorrect body")
		ape.RenderErr(w, problems.BadRequest(err)...)
	}

	var res Participants

	err = json.Unmarshal(body, &res.Request)
	if err != nil {
		Log(r).WithError(err).Error("Can't unmarshal request")
		ape.RenderErr(w, problems.BadRequest(err)...)
	}

	accountIDs := res.Request.fetchIDs()

	users, err := UsersQ(r).ByAddresses(accountIDs)
	if err != nil {
		Log(r).WithError(err).Error("Can't find users")
		ape.RenderErr(w, problems.InternalError())
	}
	res.Request = res.Request.update(users)

	res.Resource = res.Request.populate()

	json.NewEncoder(w).Encode(&res.Resource)
}

func (request ParticipantsRequest) populate() (response map[int64][]base.Participant) {
	for op, participants := range request.Participants {
		for _, participant := range participants {
			var r base.Participant
			r.Populate(&participant)
			response[op] = append(response[op], r)
		}
	}

	return
}

func (request ParticipantsRequest) fetchIDs() (accountIDs []string) {
	for _, op := range request.Participants {
		for pi := range op {
			accountIDs = append(accountIDs, string(op[pi].AccountID))
		}
	}
	return
}

func (request ParticipantsRequest) update(users []api.User) ParticipantsRequest {

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
	return request
}
