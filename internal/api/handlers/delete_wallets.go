package handlers

import (
	"net/http"

	"github.com/go-chi/chi"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"gitlab.com/swarmfund/api/internal/api/movetoape"
	"gitlab.com/swarmfund/go/doorman"
)

func DeleteWallets(w http.ResponseWriter, r *http.Request) {
	if err := Doorman(r, doorman.SignerOf(CoreInfo(r).GetMasterAccountID())); err != nil {
		movetoape.RenderDoormanErr(w, err)
		return
	}

	walletID := chi.URLParam(r, "wallet-id")

	if err := WalletQ(r).Delete(walletID); err != nil {
		Log(r).WithError(err).Error("failed to delete wallets")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
