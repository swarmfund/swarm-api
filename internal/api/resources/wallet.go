package resources

import (
	"github.com/go-ozzo/ozzo-validation"
	"gitlab.com/swarmfund/api/db2/api"
	"gitlab.com/swarmfund/api/internal/types"
	"gitlab.com/swarmfund/api/tfa"
)

type Wallet struct {
	Data WalletData `json:"data"`
}

func (r Wallet) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.Data, validation.Required),
	)
}

type WalletData struct {
	Type          string              `json:"type"`
	ID            string              `json:"id"`
	Attributes    WalletAttributes    `json:"attributes"`
	Relationships WalletRelationships `json:"relationships"`
}

func (r WalletData) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.Type, validation.Required),
		validation.Field(&r.ID, validation.Required),
		validation.Field(&r.Attributes, validation.Required),
	)
}

type WalletRelationships struct {
	KDF         *KDFPlain       `json:"kdf,omitempty"`
	Factor      *PasswordFactor `json:"factor,omitempty"`
	Transaction *Transaction    `json:"transaction,omitempty"`
	Recovery    *PasswordFactor `json:"recovery,omitempty"`
}

func (r WalletRelationships) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.KDF),
		validation.Field(&r.Factor),
		validation.Field(&r.Transaction),
		validation.Field(&r.Recovery),
	)
}

type PasswordFactor struct {
	Data PasswordFactorData `json:"data"`
}

func (r PasswordFactor) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.Data, validation.Required),
	)
}

type PasswordFactorData struct {
	ID         string                   `json:"id"`
	Type       string                   `json:"type"`
	Attributes PasswordFactorAttributes `json:"attributes"`
}

func (r PasswordFactorData) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.Type, validation.Required, validation.In("password", "recovery")),
		validation.Field(&r.Attributes, validation.Required),
	)
}

type PasswordFactorAttributes struct {
	AccountID    types.Address `json:"account_id"`
	KeychainData string        `json:"keychain_data"`
	Salt         string        `json:"salt,omitempty"`
}

func (r PasswordFactorAttributes) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.KeychainData, validation.Required),
		validation.Field(&r.Salt, validation.Required),
	)
}

type WalletAttributes struct {
	AccountID    types.Address `json:"account_id"`
	Email        string        `json:"email,omitempty"`
	KeychainData string        `json:"keychain_data"`
	Verified     bool          `json:"verified"`
	Salt         string        `json:"salt,omitempty"`
}

func (r WalletAttributes) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.AccountID, validation.Required),
		validation.Field(&r.Email, validation.Required),
		validation.Field(&r.KeychainData, validation.Required),
		validation.Field(&r.Salt, validation.Required),
	)
}

func NewWallet(w *api.Wallet) Wallet {
	r := Wallet{
		WalletData{
			Type: "wallet",
			ID:   w.WalletId,
			Attributes: WalletAttributes{
				AccountID:    w.AccountID,
				Email:        w.Username,
				KeychainData: w.KeychainData,
				Verified:     w.Verified,
			},
		},
	}
	return r
}

func NewPasswordFactor(factor *tfa.Password) *PasswordFactor {
	return &PasswordFactor{
		Data: PasswordFactorData{
			Type: "password",
			Attributes: PasswordFactorAttributes{
				Salt:         factor.Details.Salt,
				AccountID:    factor.Details.AccountID,
				KeychainData: factor.Details.KeychainData,
			},
		},
	}
}
