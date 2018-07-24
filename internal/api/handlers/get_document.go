package handlers

import (
	"net/http"

	"encoding/json"

	"github.com/go-chi/chi"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"gitlab.com/swarmfund/api/internal/api/movetoape"
	"gitlab.com/tokend/go/doorman"
)

func GetDocument(w http.ResponseWriter, r *http.Request) {
	address := chi.URLParam(r, "address")
	document := chi.URLParam(r, "document")

	constrains := []doorman.SignerConstraint{doorman.SignerOf(CoreInfo(r).GetMasterAccountID())}
	if address != "" {
		constrains = append(constrains, doorman.SignerOf(address))
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
