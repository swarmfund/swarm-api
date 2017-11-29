package resource

import (
	"fmt"

	"gitlab.com/swarmfund/api/db2/api"
)

type Wallet struct {
	PT           string `json:"paging_token"`
	AccountID    string `json:"accountId"`
	Username     string `json:"username"`
	WalletId     string `json:"walletId"`
	KeychainData string `json:"keychainData"`
	Verified     bool   `json:"verified"`
	LockVersion  int    `json:"lockVersion"`

	Detached     bool    `json:"detached"`
	Organization *string `json:"organization,omitempty"`
}

func (p *Wallet) Populate(w *api.Wallet) {
	p.PT = fmt.Sprintf("%d", w.Id)
	p.AccountID = w.AccountID
	p.Username = w.Username
	p.WalletId = w.WalletId
	p.KeychainData = w.KeychainData
	p.Verified = w.Verified
	p.LockVersion = 0
	p.Detached = w.Detached
	p.Organization = w.OrganizationAddress
}

func (p Wallet) PagingToken() string {
	return p.PT
}

type TfaKeychain struct {
	TfaSalt     string `json:"tfaSalt"`
	TfaKeychain string `json:"tfaKeychain"`
}
