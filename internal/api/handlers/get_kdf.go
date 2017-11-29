package handlers

import (
	"math"
	"net/http"

	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"gitlab.com/swarmfund/api/internal/api/resources"
)

func GetKDF(w http.ResponseWriter, r *http.Request) {
	// TODO read from db
	kdf := &resources.KDF{
		Version:   1,
		Algorithm: "scrypt",
		Bits:      256,
		N:         math.Pow(2, 12),
		R:         8,
		P:         1,
	}

	email := r.URL.Query().Get("email")
	if email != "" {
		wallet, err := WalletQ(r).ByEmail(email)
		if err != nil {
			Log(r).WithError(err).Error("failed to get wallet")
			ape.RenderErr(w, problems.InternalError())
			return
		}

		if wallet != nil {
			if wallet.KDF != 1 {
				Log(r).WithField("wallet", wallet.Id).Error("non-default kdf is not implemented")
				ape.RenderErr(w, problems.InternalError())
				return
			}

			kdf.Salt = wallet.Salt
		}
	}

	ape.Render(w, kdf)
}
