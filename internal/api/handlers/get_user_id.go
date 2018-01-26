package handlers

import (
	"net/http"

	"encoding/json"

	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
)

func GetUserID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")

	email := r.URL.Query().Get("email")

	record, err := UsersQ(r).ByEmail(email)
	if err != nil {
		Log(r).WithError(err).Error("failed to get user by email")
		return
	}

	if record == nil {
		Log(r).WithError(err).Error("No record for this user")
		ape.RenderErr(w, problems.NotFound())
		return
	}

	res := map[string]string{
		"account_id": string(record.Address),
	}

	json.NewEncoder(w).Encode(res)
}
