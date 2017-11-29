package handlers

import (
	"net/http"

	"github.com/go-chi/chi"
	. "github.com/go-ozzo/ozzo-validation"
	"github.com/google/jsonapi"
	"github.com/pkg/errors"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
)

type VerifyWalletRequest struct {
	WalletID string `json:"-"`
	Token    string `json:"token" jsonapi:"attr,token"`
}

func NewVerifyWalletRequest(r *http.Request) (VerifyWalletRequest, error) {
	request := VerifyWalletRequest{
		WalletID: chi.URLParam(r, "wallet-id"),
	}
	if err := jsonapi.UnmarshalPayload(r.Body, &request); err != nil {
		return request, errors.Wrap(err, "failed to unmarshal")
	}
	return request, request.Validate()
}

func (r *VerifyWalletRequest) Validate() error {
	return ValidateStruct(r,
		Field(&r.Token, Required),
	)
}

func VerifyWallet(w http.ResponseWriter, r *http.Request) {
	request, err := NewVerifyWalletRequest(r)
	if err != nil {
		ape.RenderErr(w, problems.BadRequest(err)...)
	}

	ok, err := EmailTokensQ(r).Verify(request.WalletID, request.Token)
	if err != nil {
		Log(r).WithError(err).Error("failed to verify token")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	if !ok {
		ape.RenderErr(w, problems.BadRequest(Errors{
			"/data/attributes/token": errors.New("invalid token"),
		})...)
		return
	}

	w.WriteHeader(204)
}
