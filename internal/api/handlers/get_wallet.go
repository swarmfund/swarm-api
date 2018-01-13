package handlers

import (
	"net/http"

	"encoding/json"

	"github.com/go-chi/chi"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"gitlab.com/swarmfund/api/internal/api/movetoape"
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
		ape.RenderErr(w, movetoape.Forbidden("verification_required"))
		return
	}

	// TODO check block state

	if err := secondfactor.NewConsumer(TFAQ(r)).WithBackendType(types.WalletFactorTOTP).Consume(r, wallet); err != nil {
		RenderFactorConsumeError(w, r, err)
		return
	}

	{
		resource := resources.NewWallet(wallet)
		json.NewEncoder(w).Encode(&resource)
	}

}
