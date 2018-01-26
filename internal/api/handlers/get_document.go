package handlers

import (
	"net/http"

	"encoding/json"

	"github.com/go-chi/chi"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"gitlab.com/swarmfund/api/internal/api/movetoape"
	"gitlab.com/swarmfund/go/doorman"
)

func GetDocument(w http.ResponseWriter, r *http.Request) {
	address := chi.URLParam(r, "address")
	document := chi.URLParam(r, "document")

	if err := Doorman(r,
		doorman.SignerOf(address),
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
