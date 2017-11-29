package handlers

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/google/jsonapi"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"gitlab.com/swarmfund/api/internal/api/resources"
	"gitlab.com/swarmfund/api/internal/secondfactor"
	"gitlab.com/swarmfund/api/internal/types"
)

func GetWallet(w http.ResponseWriter, r *http.Request) {
	walletID := chi.URLParam(r, "wallet-id")

	wallet, err := WalletQ(r).ByWalletID(walletID)
	if err != nil {
		Log(r).WithError(err).Error("failed to get wallet")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	if wallet == nil {
		ape.RenderErr(w, problems.NotFound())
		return
	}

	if !wallet.Verified {
		ape.RenderErr(w, &jsonapi.ErrorObject{
			Title:  http.StatusText(http.StatusForbidden),
			Status: fmt.Sprintf("%d", http.StatusForbidden),
			Detail: "Email should be verified before login",
		})
		return
	}

	// TODO check block state

	if err := secondfactor.NewConsumer(TFAQ(r)).WithBackendType(types.WalletFactorTOTP).Consume(r, wallet); err != nil {
		RenderFactorConsumeError(w, r, err)
		return
	}

	ape.Render(w, resources.NewWallet(wallet, nil))
}
