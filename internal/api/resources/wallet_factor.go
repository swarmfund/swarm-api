package resources

import "gitlab.com/swarmfund/api/internal/types"

type WalletFactor struct {
	Type       types.WalletFactorType `json:"type"`
	ID         int64                  `json:"id"`
	Attributes WalletFactorAttributes `json:"attributes"`
}

type WalletFactorAttributes struct {
	Priority uint `json:"priority"`
}
