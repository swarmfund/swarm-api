package handlers

import (
	"net/http"

	"encoding/json"

	"github.com/go-chi/chi"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"gitlab.com/swarmfund/api/internal/api/movetoape"
	"gitlab.com/swarmfund/api/internal/api/resources"
	"gitlab.com/swarmfund/go/doorman"
)

type GetUserResponse struct {
	Data resources.User `json:"data"`
}

func GetUser(w http.ResponseWriter, r *http.Request) {
	address := chi.URLParam(r, "address")

	// TODO unhardcode
	if err := Doorman(r,
		doorman.SignerOf(address),
		doorman.SignerOf("GD7AHJHCDSQI6LVMEJEE2FTNCA2LJQZ4R64GUI3PWANSVEO4GEOWB636"),
	); err != nil {
		movetoape.RenderDoormanErr(w, err)
		return
	}

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

	response := GetUserResponse{
		Data: resources.NewUser(user),
	}
	json.NewEncoder(w).Encode(response)
}
