package resources

import (
	"gitlab.com/swarmfund/api/db2/api"
	"gitlab.com/swarmfund/api/tfa"
)

type Wallet struct {
	WalletID     string          `jsonapi:"primary,wallet"`
	AccountID    string          `jsonapi:"attr,account_id"`
	Email        string          `jsonapi:"attr,email"`
	KeychainData string          `jsonapi:"attr,keychain_data"`
	Verified     bool            `jsonapi:"attr,verified"`
	Factor       *PasswordFactor `jsonapi:"relation,factor"`
}

type WalletData struct {
	Type       string           `json:"type"`
	ID         string           `json:"id"`
	Attributes WalletAttributes `json:"attributes"`
}

type WalletAttributes struct {
	AccountID    string `json:"account_id"`
	Email        string `json:"email"`
	KeychainData string `json:"keychain_data"`
	Verified     bool   `json:"verified"`
}

func NewWalletData(w *api.Wallet) WalletData {
	r := WalletData{
		Type: "wallet",
		ID:   w.WalletId,
		Attributes: WalletAttributes{
			AccountID:    w.AccountID,
			Email:        w.Username,
			KeychainData: w.KeychainData,
			Verified:     w.Verified,
		},
	}
	return r
}

func NewWallet(w *api.Wallet, password *tfa.Password) *Wallet {
	r := &Wallet{
		WalletID:     w.WalletId,
		AccountID:    w.AccountID,
		Email:        w.Username,
		KeychainData: w.KeychainData,
		Verified:     w.Verified,
	}
	if password != nil {
		r.Factor = &PasswordFactor{
			ID:           password.ID(),
			Salt:         password.Details.Salt,
			AccountID:    password.Details.AccountID,
			KeychainData: password.Details.KeychainData,
		}
	}
	return r
}
