package handlers

import (
	"encoding/json"
	"net/http"

	. "github.com/go-ozzo/ozzo-validation"
	"github.com/pkg/errors"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"gitlab.com/swarmfund/api/db2/api"
	"gitlab.com/swarmfund/api/internal/api/resources"
	"gitlab.com/swarmfund/api/internal/data"
	"gitlab.com/swarmfund/api/internal/lorem"
	"gitlab.com/swarmfund/api/tfa"
)

type (
	CreateWalletRequest struct {
		resources.Wallet
	}
)

func NewCreateWalletRequest(r *http.Request) (CreateWalletRequest, error) {
	request := CreateWalletRequest{}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return request, errors.Wrap(err, "failed to unmarshal")
	}
	return request, request.Validate()
}

func (r *CreateWalletRequest) Validate() error {
	errs := Errors{
		"/data/":                       Validate(r.Data, Required),
		"/data/relationships/kdf":      Validate(r.Data.Relationships.KDF, Required),
		"/data/relationships/factor":   Validate(r.Data.Relationships.Factor, Required),
		"/data/relationships/recovery": Validate(r.Data.Relationships.Recovery, Required),
	}
	if r.Data.Relationships.Recovery != nil {
		errs["/data/relationships/recovery/account_id"] = Validate(
			r.Data.Relationships.Recovery.Data.Attributes.AccountID, Required)
	}
	return errs.Filter()
}

func CreateWallet(w http.ResponseWriter, r *http.Request) {
	request, err := NewCreateWalletRequest(r)
	if err != nil {
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	wallet := &api.Wallet{
		Username:         request.Data.Attributes.Email,
		AccountID:        request.Data.Attributes.AccountID,
		CurrentAccountID: request.Data.Attributes.AccountID,
		WalletId:         request.Data.ID,
		KeychainData:     request.Data.Attributes.KeychainData,
	}

	walletKDF := data.WalletKDF{
		Wallet:  request.Data.Attributes.Email,
		Version: request.Data.Relationships.KDF.Data.ID,
		Salt:    request.Data.Attributes.Salt,
	}

	factor := tfa.NewPasswordBackend(tfa.PasswordDetails{
		Salt:         request.Data.Relationships.Factor.Data.Attributes.Salt,
		AccountID:    request.Data.Relationships.Factor.Data.Attributes.AccountID,
		KeychainData: request.Data.Relationships.Factor.Data.Attributes.KeychainData,
	})

	err = WalletQ(r).Transaction(func(q api.WalletQI) error {
		existing, err := q.ByEmail(wallet.Username)
		if err != nil {
			return errors.Wrap(err, "failed to get wallet")
		}

		if existing != nil {
			return api.ErrWalletsConflict
		}

		if err := q.Create(wallet); err != nil {
			return errors.Wrap(err, "failed to create wallet")
		}

		if err = q.CreatePasswordFactor(wallet.WalletId, factor); err != nil {
			return errors.Wrap(err, "failed to create password factor")
		}

		if err = q.CreateRecovery(api.RecoveryKeychain{
			Email:    request.Data.Attributes.Email,
			Salt:     request.Data.Relationships.Recovery.Data.Attributes.Salt,
			Keychain: request.Data.Relationships.Recovery.Data.Attributes.KeychainData,
			WalletID: request.Data.Relationships.Recovery.Data.ID,
			Address:  request.Data.Relationships.Recovery.Data.Attributes.AccountID,
		}); err != nil {
			return errors.Wrap(err, "failed to create recovery")
		}

		if err := q.CreateWalletKDF(walletKDF); err != nil {
			return errors.Wrap(err, "failed to create wallet kdf")
		}
		return nil
	})
	if err != nil {
		cause := errors.Cause(err)

		if cause == api.ErrWalletsConflict || cause == api.ErrWalletsWalletIDViolated || cause == api.ErrRecoveriesConflict {
			ape.RenderErr(w, problems.Conflict())
			return
		}

		if cause == api.ErrWalletsKDFViolated {
			ape.RenderErr(w, problems.BadRequest(Errors{
				"/data/relationships/kdf/data/id": errors.New("invalid kdf version"),
			})...)
			return
		}

		Log(r).WithError(err).Error("failed to save wallet")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	{
		resource := resources.NewWallet(wallet)
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(&resource)
	}

	// wallet has been saved, so technically request has succeeded
	// no errors should be rendered from now on
	// TODO move token create to transaction
	if err := EmailTokensQ(r).Create(wallet.WalletId, lorem.Token()); err != nil {
		Log(r).WithError(err).Error("failed to save token")
		return
	}
}
