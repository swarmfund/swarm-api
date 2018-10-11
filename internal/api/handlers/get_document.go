package handlers

import (
	"net/http"

	"gitlab.com/swarmfund/api/internal/storage"

	"encoding/json"

	"github.com/go-chi/chi"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"gitlab.com/swarmfund/api/internal/api/movetoape"
	"gitlab.com/tokend/go/doorman"
)

func GetDocument(w http.ResponseWriter, r *http.Request) {
	document := chi.URLParam(r, "document")

	key := storage.Key{}
	if err := key.UnmarshalText([]byte(document)); err != nil {
		Log(r).WithError(err).Error("failed to unmarshal key")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	user, err := UsersQ(r).ByID(key.UserID)
	if err != nil {
		Log(r).WithError(err).Error("failed to get user by id")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	constrains := []doorman.SignerConstraint{doorman.SignerOf(CoreInfo(r).GetMasterAccountID())}
	if user != nil {
		constrains = append(constrains, doorman.SignerOf(string(user.Address)))
	}

	if err := Doorman(r, constrains...); err != nil {
		movetoape.RenderDoormanErr(w, err)
		return
	}

	url, err := Storage(r).SignedObjectURL(document)
	if err != nil {
		Log(r).WithError(err).Error("failed get document url")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	json.NewEncoder(w).Encode(map[string]string{
		"url": url.String(),
	})
}
