package resources

import (
	"time"

	"gitlab.com/swarmfund/api/internal/types"
)

type Blob struct {
	ID         string         `json:"id"`
	Type       types.BlobType `json:"type"`
	Attributes struct {
		Value     string     `json:"value"`
		DeletedAt *time.Time `json:"deleted_at,omitempty"`
	} `json:"attributes"`
}

func NewBlob(blob *types.Blob) Blob {
	b := Blob{
		ID:   blob.ID,
		Type: blob.Type,
	}
	b.Attributes.Value = blob.Value
	b.Attributes.DeletedAt = blob.DeletedAt
	return b
}
