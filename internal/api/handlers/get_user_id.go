package handlers

import (
	"net/http"

	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"

	"encoding/json"
)

func GetUserId(w http.ResponseWriter, r *http.Request) {
	email := r.URL.Query().Get("email")

	user, err := UsersQ(r).ByEmail(email)
	if err != nil {
		Log(r).WithError(err).Error("failed to get user")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	if user == nil {
		ape.RenderErr(w, problems.NotFound())
		return
	}

	json.NewEncoder(w).Encode(map[string]string{
		"account_id": string(user.Address),
	})
}
