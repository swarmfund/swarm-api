package horizon

import (
	"database/sql"

	"github.com/lann/squirrel"
	"github.com/lib/pq"
	"github.com/pkg/errors"
	"gitlab.com/swarmfund/api/db2"
	"gitlab.com/swarmfund/api/internal/types"
)

const (
	blobsTable        = "blobs"
	blobsPKConstraint = "blobs_pkey"
)

type Blobs struct {
	*db2.Repo
}

func (q *Blobs) Create(address types.Address, blob *types.Blob) error {
	stmt := squirrel.Insert(blobsTable).SetMap(map[string]interface{}{
		"owner_address": address,
		"id":            blob.ID,
		"type":          blob.Type,
		"value":         blob.Value,
		"relationships": blob.Relationships,
	})

	_, err := q.Exec(stmt)
	if err != nil {
		cause := errors.Cause(err)
		pqerr, ok := cause.(*pq.Error)
		if ok {
			if pqerr.Constraint == blobsPKConstraint {
				// id is deterministic based on blob,
				// silencing error makes request idempotent
				return nil
			}
		}
	}
	return err
}

func (q *Blobs) Filter(owner string, filters map[string]string) ([]types.Blob, error) {
	var result []types.Blob
	stmt := squirrel.
		Select("id", "value", "type", "relationships").
		From(blobsTable).
		Where("owner_address = ?", owner)

	for k, v := range filters {
		stmt = stmt.Where("relationships->>? = ?", k, v)
	}

	err := q.Repo.Select(&result, stmt)
	return result, err
}

func (q *Blobs) Get(id string) (*types.Blob, error) {
	var result types.Blob
	stmt := squirrel.
		Select("id", "value", "type", "relationships").
		From(blobsTable).
		Where("id = ?", id)

	err := q.Repo.Get(&result, stmt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &result, nil
}
