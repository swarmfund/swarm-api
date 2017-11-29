package handlers

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-ozzo/ozzo-validation"
	"github.com/pkg/errors"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
)

func RequestVerification(w http.ResponseWriter, r *http.Request) {
	walletID := chi.URLParam(r, "wallet-id")
	token, err := EmailTokensQ(r).Get(walletID)
	if err != nil {
		Log(r).WithError(err).Error("failed to get token")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	if token == nil {
		// token is expected to exist, probably messed up database state
		Log(r).WithField("wallet", walletID).Error("token expected to exist")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	if token.Confirmed {
		ape.RenderErr(w, problems.BadRequest(validation.Errors{
			"": errors.New("email already confirmed"),
		})...)
		return
	}

	w.WriteHeader(204)
}
