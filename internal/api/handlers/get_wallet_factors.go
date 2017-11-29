package handlers

import (
	"net/http"

	"encoding/json"

	"github.com/go-chi/chi"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"gitlab.com/swarmfund/api/internal/api/resources"
)

type GetWalletFactorsResponse struct {
	Data []resources.WalletFactor `json:"data"`
}

func GetWalletFactors(w http.ResponseWriter, r *http.Request) {
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

	// todo check allowed

	records, err := TFAQ(r).Backends(wallet.WalletId)
	if err != nil {
		Log(r).WithError(err).Error("failed to get backends")
		ape.RenderErr(w, problems.InternalError())
		return
	}
	response := GetWalletFactorsResponse{
		Data: make([]resources.WalletFactor, 0, len(records)),
	}
	for _, record := range records {
		response.Data = append(response.Data, resources.WalletFactor{
			Type: record.BackendType,
			ID:   record.ID,
			Attributes: resources.WalletFactorAttributes{
				Priority: record.Priority,
			},
		})
	}

	json.NewEncoder(w).Encode(&response)
}
