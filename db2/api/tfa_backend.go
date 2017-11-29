package api

import (
	"fmt"

	"gitlab.com/swarmfund/api/internal/types"
	"gitlab.com/swarmfund/api/tfa"
)

type Backend struct {
	ID          int64                  `db:"id"`
	WalletID    string                 `db:"wallet_id"`
	Details     []byte                 `db:"details"`
	BackendType types.WalletFactorType `db:"backend"`
	Priority    uint                   `db:"priority"`
}

func (b *Backend) Backend() (tfa.Backend, error) {
	switch b.BackendType {
	case types.WalletFactorTOTP:
		return tfa.NewTOTPFromDB(b.ID, b.Details)
	case types.WalletFactorPassword:
		return tfa.NewPasswordFromDB(b.ID, b.Details)
	default:
		return nil, fmt.Errorf("unknown backend %s", b.BackendType)
	}
}
