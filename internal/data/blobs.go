package data

import "gitlab.com/swarmfund/api/internal/types"

//go:generate mockery -case underscore -name Blobs
type Blobs interface {
	Create(address types.Address, blob *types.Blob) error
	Get(id string) (*types.Blob, error)
	Filter(owner string, filters map[string]string) ([]types.Blob, error)
}
