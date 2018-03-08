package handlers

import (
	"net/http"

	"strconv"

	"encoding/json"

	"strings"

	"github.com/pkg/errors"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/swarmfund/api/internal/api/resources"
	"gitlab.com/swarmfund/api/internal/data"
)

func getCaseInsensitiveKDF(r *http.Request, email string) (kdf *data.KDF, emailFound string, err error) {
	// first try as-is
	kdf, err = WalletQ(r).KDFByEmail(email)
	if err != nil {
		return nil, email, errors.Wrap(err, "failed to get wallet kdf")
	}
	if kdf == nil {
		// if kdf is not found let's try lower-case
		email = strings.ToLower(email)
		kdf, err = WalletQ(r).KDFByEmail(email)
		if err != nil {
			return nil, email, errors.Wrap(err, "failed to get wallet kdf")
		}
	}
	return kdf, email, nil
}

func GetKDF(w http.ResponseWriter, r *http.Request) {

	email := r.URL.Query().Get("email")
	// if invalid just ignore it
	isRecovery, _ := strconv.ParseBool(r.URL.Query().Get("is_recovery"))

	var err error
	var kdf *data.KDF

	switch {
	case email == "": // load default KDF
		// TODO move version to config
		kdf, err = WalletQ(r).KDFByVersion(2)
		if err != nil {
			Log(r).WithError(err).Error("failed to get kdf")
			ape.RenderErr(w, problems.InternalError())
			return
		}
	case email != "" && !isRecovery: // load wallet KDF
		kdf, _, err = getCaseInsensitiveKDF(r, email)
		if err != nil {
			Log(r).WithError(err).Error("failed to get wallet kdf")
			ape.RenderErr(w, problems.InternalError())
			return
		}
	case email != "" && isRecovery: // load recovery KDF for wallet
		// FIXME it's bad, you need to move recoveries to kdf_wallets
		// load wallet to get recovery salt
		kdf, email, err = getCaseInsensitiveKDF(r, email)
		if err != nil {
			Log(r).WithError(err).Error("failed to get wallet kdf")
			ape.RenderErr(w, problems.InternalError())
			return
		}
		if kdf == nil {
			// 404 will be rendered below
			break
		}
		wallet, err := WalletQ(r).ByEmail(email)
		if err != nil {
			Log(r).WithError(err).Error("failed to get wallet")
			ape.RenderErr(w, problems.InternalError())
			return
		}
		if wallet == nil {
			// 404 will be rendered below
			kdf = nil
			break
		}
		kdf.Salt = wallet.RecoverySalt
	default:
		Log(r).WithFields(logan.F{
			"is_recovery": isRecovery,
			"email":       email,
		}).Error("undefined state")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	if kdf == nil {
		ape.RenderErr(w, problems.NotFound())
		return
	}

	json.NewEncoder(w).Encode(struct {
		Data resources.KDF `json:"data"`
	}{resources.NewKDF(*kdf)})
}
