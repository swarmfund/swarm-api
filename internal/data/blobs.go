package data

import (
	"gitlab.com/swarmfund/api/internal/types"
)

//go:generate mockery -case underscore -name Blobs
type Blobs interface {
	New() Blobs
	Transaction(fn func(Blobs) error) error
	Delete(...types.Blob) error
	Create(address types.Address, blob *types.Blob) error
	Get(id string) (*types.Blob, error)

	// filter
	ByOwner(types.Address) Blobs
	ByType(types.BlobType) Blobs
	ByRelationships(map[string]string) Blobs
	Select() ([]types.Blob, error)
}
