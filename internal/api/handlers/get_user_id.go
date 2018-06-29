package handlers

import (
	"net/http"

	"github.com/go-chi/chi"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"

	"encoding/json"

	"gitlab.com/swarmfund/api/db2/api"
)

func GetUserId(w http.ResponseWriter, r *http.Request) {
	email := chi.URLParam(r, "email")

	q := UsersQ(r).EmailMatches(email)
	var user api.User
	if err := q.Select(&user); err != nil {
		Log(r).WithError(err).Error("failed to get user")
		ape.RenderErr(w, problems.NotFound())
		return
	}

	json.NewEncoder(w).Encode(map[string]string{
		"account_id": string(user.Address),
	})
	return
}
