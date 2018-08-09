package data

import (
	"gitlab.com/swarmfund/api/internal/types"
)

//go:generate mockery -case underscore -name Blobs
type Blobs interface {
	New() Blobs
	Transaction(fn func(Blobs) error) error
	// Delete is used for some sorcery, let's stay away from that
	Delete(...types.Blob) error
	// MarkDeleted make blob hidden by default
	MarkDeleted(id string) error
	Create(address *types.Address, blob *types.Blob) error
	Get(id string) (*types.Blob, error)

	// filter
	ByOwner(types.Address) Blobs
	ByType(types.BlobType) Blobs
	ExcludeDeleted() Blobs
	ByRelationships(map[string]string) Blobs
	Select() ([]types.Blob, error)
}
