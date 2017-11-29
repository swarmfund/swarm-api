package handlers

import (
	"net/http"

	. "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"github.com/google/jsonapi"
	"github.com/pkg/errors"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"gitlab.com/swarmfund/api/db2/api"
	"gitlab.com/swarmfund/api/internal/api/resources"
	"gitlab.com/swarmfund/api/internal/lorem"
	"gitlab.com/swarmfund/api/tfa"
)

type CreateWalletRequest struct {
	WalletID       string                    `json:"wallet_id" jsonapi:"primary,wallet"`
	Email          string                    `json:"email" jsonapi:"attr,email"`
	AccountID      string                    `json:"account_id" jsonapi:"attr,account_id"`
	Salt           string                    `json:"salt" jsonapi:"attr,salt"`
	KeychainData   string                    `json:"keychain_data" jsonapi:"attr,keychain_data"`
	KDF            *resources.KDFVersion     `json:"kdf" jsonapi:"relation,kdf"`
	PasswordFactor *resources.PasswordFactor `json:"factor" jsonapi:"relation,factor"`
}

func NewCreateWalletRequest(r *http.Request) (CreateWalletRequest, error) {
	request := CreateWalletRequest{}
	if err := jsonapi.UnmarshalPayload(r.Body, &request); err != nil {
		return request, errors.Wrap(err, "failed to unmarshal")
	}
	return request, request.Validate()
}

func (r *CreateWalletRequest) Validate() error {
	return ValidateStruct(r,
		Field(&r.Email, Required, is.Email),
		// TODO account id validation
		Field(&r.AccountID, Required),
		Field(&r.WalletID, Required),
		Field(&r.Salt, Required),
		Field(&r.KeychainData, Required),
		Field(&r.KDF, Required),
		Field(&r.PasswordFactor, Required),
	)
}

func CreateWallet(w http.ResponseWriter, r *http.Request) {
	request, err := NewCreateWalletRequest(r)
	if err != nil {
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	wallet := &api.Wallet{
		Username:         request.Email,
		AccountID:        request.AccountID,
		CurrentAccountID: request.AccountID,
		WalletId:         request.WalletID,
		Salt:             request.Salt,
		KDF:              request.KDF.Version,
		KeychainData:     request.KeychainData,
	}

	factor := tfa.NewPasswordBackend(tfa.PasswordDetails{
		Salt:         request.PasswordFactor.Salt,
		AccountID:    request.PasswordFactor.AccountID,
		KeychainData: request.PasswordFactor.KeychainData,
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

		return nil
	})
	if err != nil {
		cause := errors.Cause(err)

		if cause == api.ErrWalletsConflict {
			ape.RenderErr(w, problems.Conflict())
			return
		}

		if cause == api.ErrWalletsWalletIDViolated {
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

	// wallet has been saved, so technically request has succeeded
	// no errors should be rendered from now on
	ape.Render(w, resources.NewWallet(wallet, factor))

	if err := EmailTokensQ(r).Create(wallet.WalletId, lorem.Token()); err != nil {
		Log(r).WithError(err).Error("failed to save token")
		return
	}
}
