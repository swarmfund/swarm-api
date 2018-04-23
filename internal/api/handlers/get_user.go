package handlers

import (
	"net/http"

	"encoding/json"

	"github.com/go-chi/chi"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"gitlab.com/swarmfund/api/internal/api/movetoape"
	"gitlab.com/swarmfund/api/internal/api/resources"
	"gitlab.com/swarmfund/api/internal/track"
	"gitlab.com/tokend/go/doorman"
)

type GetUserResponse struct {
	Data resources.User `json:"data"`
}

func GetUser(w http.ResponseWriter, r *http.Request) {
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

	response := GetUserResponse{
		Data: resources.NewUser(user),
	}

	defer json.NewEncoder(w).Encode(response)

	// FIXME
	{
		if err := Doorman(r, doorman.SignerOf(CoreInfo(r).GetMasterAccountID())); err != nil {
			return
		}

		event, err := Tracker(r).GetLast(track.Event{
			Address: string(user.Address),
		})
		if err != nil {
			Log(r).WithError(err).Error("failed to get event")
		}

		if event == nil {
			return
		}
		response.Data.Attributes.LastIPAddress = event.Details.Request.IP
	}

}
