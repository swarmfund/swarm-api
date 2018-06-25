package handlers

import (
	"github.com/go-chi/chi"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"net/http"

	"encoding/json"
	"gitlab.com/swarmfund/api/internal/api/movetoape"
	"gitlab.com/swarmfund/api/internal/api/resources"
	"gitlab.com/tokend/go/doorman"
)

func GetUserId(w http.ResponseWriter, r *http.Request) {
	address := chi.URLParam(r, "address")

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

	if err := Doorman(r,
		doorman.SignerOf(address),
		doorman.SignerOf(CoreInfo(r).GetMasterAccountID()),
	); err != nil {
		movetoape.RenderDoormanErr(w, err)
		return
	}

	response := map[string]string{
		"account_id": string(resources.NewUser(user).ID),
	}

	defer json.NewEncoder(w).Encode(&response)

}
