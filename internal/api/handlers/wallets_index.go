package handlers

import (
	"net/http"

	"encoding/json"

	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"gitlab.com/swarmfund/api/internal/api/movetoape"
	"gitlab.com/swarmfund/api/internal/api/resources"
	"gitlab.com/swarmfund/go/doorman"
)

type (
	WalletsIndexResponse struct {
		Data WalletsIndexData `json:"data"`
	}
	WalletsIndexData []*resources.Wallet
)

func WalletsIndex(w http.ResponseWriter, r *http.Request) {
	if err := Doorman(r, doorman.SignerOf("master-account")); err != nil {
		movetoape.RenderDoormanErr(w, err)
		return
	}

	// load paging params
	wallets, err := WalletQ(r).Select()
	if err != nil {
		Log(r).WithError(err).Error("failed to get wallets")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	wallets = wallets

	response := WalletsIndexResponse{
		Data: make(WalletsIndexData, 0, len(wallets)),
	}
	for _, wallet := range wallets {
		response.Data = append(response.Data, resources.NewWallet(&wallet, nil))
	}
	json.NewEncoder(w).Encode(&response)
}
