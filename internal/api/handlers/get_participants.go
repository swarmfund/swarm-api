package handlers

import (
	"net/http"

	"encoding/json"

	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"gitlab.com/swarmfund/api/db2/api"
)

type participantsRequest struct {
	ForAccount   string                      `json:"for_account" valid:"required"`
	Participants map[int64][]api.Participant `json:"participants"`
}

func GetParticipants(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")

	var request participantsRequest

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&request); err != nil {
		Log(r).WithError(err).Error("failed to get participant details")
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	err := UsersQ(r).Participants(request.Participants)
	if err != nil {
		Log(r).WithError(err).Error("failed to get participant details")
		ape.RenderErr(w, problems.InternalError()) //server error
		return
	}

	json.NewEncoder(w).Encode(request.Participants)
}
