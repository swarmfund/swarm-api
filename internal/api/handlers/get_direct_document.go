package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"gitlab.com/swarmfund/api/internal/api/movetoape"
	"gitlab.com/swarmfund/go/doorman"
)

func GetDirectDocument(w http.ResponseWriter, r *http.Request) {
	document := chi.URLParam(r, "document")

	if err := Doorman(r,
		doorman.SignerOf(CoreInfo(r).GetMasterAccountID()),
	); err != nil {
		movetoape.RenderDoormanErr(w, err)
		return
	}

	url, err := Storage(r).DocumentURL(document)
	if err != nil {
		Log(r).WithError(err).Error("failed get document url")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	json.NewEncoder(w).Encode(map[string]string{
		"url": url.String(),
	})
}
