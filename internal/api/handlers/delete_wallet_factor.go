package handlers

import (
	"net/http"

	"strconv"

	"github.com/go-chi/chi"
	. "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"gitlab.com/swarmfund/api/internal/api/movetoape"
	"gitlab.com/swarmfund/api/internal/secondfactor"
	"gitlab.com/swarmfund/api/internal/types"
	"gitlab.com/swarmfund/go/doorman"
)

type DeleteWalletFactorRequest struct {
	WalletID  string `json:"-"`
	BackendID string `json:"-"`
}

func NewDeleteWalletFactorRequest(r *http.Request) (DeleteWalletFactorRequest, error) {
	request := DeleteWalletFactorRequest{
		WalletID:  chi.URLParam(r, "wallet-id"),
		BackendID: chi.URLParam(r, "backend"),
	}
	return request, request.Validate()
}

func (r DeleteWalletFactorRequest) Validate() error {
	return ValidateStruct(&r,
		Field(&r.WalletID, Required),
		Field(&r.BackendID, Required, is.Int),
	)
}

func DeleteWalletFactor(w http.ResponseWriter, r *http.Request) {
	request, err := NewDeleteWalletFactorRequest(r)
	if err != nil {
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	wallet, err := WalletQ(r).ByWalletID(request.WalletID)
	if err != nil {
		Log(r).WithError(err).Error("failed to get wallet")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	if wallet == nil {
		ape.RenderErr(w, problems.NotFound())
		return
	}

	// check allowed
	if err := Doorman(r, doorman.SignerOf(wallet.AccountID)); err != nil {
		movetoape.RenderDoormanErr(w, err)
		return
	}

	// ask password before writing anything
	if err := secondfactor.NewConsumer(TFAQ(r)).WithBackendType(types.WalletFactorPassword).Consume(r, wallet); err != nil {
		RenderFactorConsumeError(w, r, err)
		return
	}

	bid, err := strconv.ParseInt(request.BackendID, 10, 64)
	if err != nil {
		Log(r).WithError(err).Error("failed to parse backend id")
		ape.RenderErr(w, problems.InternalError())
		return
	}
	if err := TFAQ(r).DeleteBackend(bid); err != nil {
		Log(r).WithError(err).Error("failed to delete backend")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	w.WriteHeader(204)
}
