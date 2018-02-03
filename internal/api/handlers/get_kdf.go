package handlers

import (
	"math"
	"net/http"

	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"gitlab.com/swarmfund/api/internal/api/resources"
	"strconv"
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
	// if invalid just ignore it
	isRecovery, _ := strconv.ParseBool(r.URL.Query().Get("is_recovery"))

	if email != "" {
		wallet, err := WalletQ(r).ByEmail(email)
		if err != nil {
			Log(r).WithError(err).Error("failed to get wallet")
			ape.RenderErr(w, problems.InternalError())
			return
		}

		if wallet == nil {
			ape.RenderErr(w, problems.NotFound())
			return
		}

		if wallet.KDF != 1 {
			Log(r).WithField("wallet", wallet.Id).Error("non-default kdf is not implemented")
			ape.RenderErr(w, problems.InternalError())
			return
		}

		kdf.Salt = wallet.Salt
		if isRecovery {
			if wallet.RecoverySalt == nil {
				Log(r).WithField("wallet", wallet.Id).Error("does not have recovery salt")
				ape.RenderErr(w, problems.InternalError())
				return
			}

			kdf.Salt = *wallet.RecoverySalt
		}
	}

	ape.Render(w, kdf)
}
