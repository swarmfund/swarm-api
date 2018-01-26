package handlers

import (
	"encoding/json"
	"net/http"

	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
)

type detailsRequest struct {
	Addresses []string `json:"addresses"`
}

func GetDetails(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")

	request := detailsRequest{}

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&request); err != nil {
		Log(r).WithError(err).Error("failed decode request")
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	users, err := UsersQ(r).ByAddresses(request.Addresses)
	if err != nil {
		Log(r).WithError(err).Error("failed to get addresses")
	}

	for _, user := range users {
		request.Addresses = append(request.Addresses, string(user.Address))
	}

	json.NewEncoder(w).Encode(request)
}
