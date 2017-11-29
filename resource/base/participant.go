package base

import (
	"gitlab.com/swarmfund/api/db2/api"
	"gitlab.com/swarmfund/api/internal/types"
)

type Participant struct {
	AccountID types.Address `json:"account_id,omitempty"`
	BalanceID string        `json:"balance_id,omitempty"`
	Email     *string       `json:"email,omitempty"`
	Effects   string        `json:"effects,omitempty"`
}

func (f *Participant) Populate(p *api.Participant) {
	f.AccountID = p.AccountID
	f.BalanceID = p.BalanceID
	f.Email = p.Email
	f.Effects = p.Effects
}
