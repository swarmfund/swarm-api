package resource

import (
	"gitlab.com/swarmfund/api/db2/api"
	"gitlab.com/swarmfund/api/internal/types"
)

type TFABackends struct {
	Backends []TFABackend `json:"backends"`
}

type TFABackend struct {
	ID       int64                  `json:"id"`
	Type     types.WalletFactorType `json:"type"`
	Priority uint                   `json:"priority"`
}

func (r *TFABackend) Populate(b *api.Backend) {
	r.ID = b.ID
	r.Type = b.BackendType
	r.Priority = b.Priority
}
