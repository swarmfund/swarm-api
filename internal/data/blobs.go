package data

import "gitlab.com/swarmfund/api/internal/types"

//go:generate mockery -case underscore -name Blobs
type Blobs interface {
	Create(address types.Address, blob *types.Blob) error
}
