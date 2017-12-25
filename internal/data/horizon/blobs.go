package horizon

import (
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
