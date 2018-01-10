package handlers

import (
	"net/http"

	"encoding/json"

	"github.com/go-chi/chi"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"gitlab.com/swarmfund/api/internal/api/resources"
)

type (
	KYCEntitiesIndexResponse struct {
		Data []resources.KYCEntity `json:"data"`
	}
)

func KYCEntitiesIndex(w http.ResponseWriter, r *http.Request) {
	address := chi.URLParam(r, "address")

	// TODO check allowed

	// load user
	user, err := UsersQ(r).ByAddress(address)
	if err != nil {
		Log(r).WithError(err).Error("failed to get user")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	if user == nil {
		ape.RenderErr(w, problems.NotFound())
		return
	}

	// load records
	entities, err := UsersQ(r).KYC().Select(user.ID)
	if err != nil {
		Log(r).WithError(err).Error("failed to get entities")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	// craft response
	response := KYCEntitiesIndexResponse{
		Data: []resources.KYCEntity{},
	}
	for _, entity := range entities {
		response.Data = append(response.Data, resources.NewKYCEntity(entity))
	}

	json.NewEncoder(w).Encode(&response)
}
