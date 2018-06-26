package handlers

import (
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"gitlab.com/swarmfund/api/db2/api"
	"gitlab.com/swarmfund/api/resource/base"
	"io/ioutil"
	"net/http"
	"gitlab.com/tokend/go/doorman"
	"gitlab.com/swarmfund/api/internal/api/movetoape"
	"encoding/json"
	"gitlab.com/swarmfund/api/internal/types"
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

	//check allowed
	if err := Doorman(r,
		doorman.SignerOf(CoreInfo(r).GetMasterAccountID()),
	); err != nil {
		movetoape.RenderDoormanErr(w, err)
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		Log(r).WithError(err).Error("Incorrect body")
		ape.RenderErr(w, problems.BadRequest(err)...)
	}

	var res Participants

	err = json.Unmarshal(body, &res.Request)
	if err != nil{
		Log(r).WithError(err).Error("Can't unmarshal request")
		ape.RenderErr(w, problems.BadRequest(err)...)
	}


	//TODO single query
	users := make([]api.User, 0)
	for _, op := range res.Request.Participants{
		for pi := range op{
			user, err := UsersQ(r).ByAddress(string(op[pi].AccountID))
			if err != nil{
				Log(r).WithError(err).Error("failed to get users")
				ape.RenderErr(w, problems.InternalError())
				return
			}
			users = append(users, *user)
		}
	}

	usersMap := map[types.Address]api.User{}
	for _, user := range users {
		usersMap[user.Address] = user
	}

	for _, op := range res.Request.Participants {
		for pi := range op {
			participant := op[pi]
			if user, ok := usersMap[participant.AccountID]; ok {
				participant.Email = &user.Email
				op[pi] = participant
			}
		}
	}

	res.Resource = map[int64][]base.Participant{}
	for op, participants := range res.Request.Participants {
		for _, participant := range participants {
			var r base.Participant
			r.Populate(&participant)
			res.Resource[op] = append(res.Resource[op], r)
		}
	}

	json.NewEncoder(w).Encode(&res.Resource)
}
