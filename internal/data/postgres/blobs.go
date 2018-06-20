package postgres

import (
	"database/sql"

	"time"

	"github.com/lann/squirrel"
	"github.com/lib/pq"
	"github.com/pkg/errors"
	"gitlab.com/swarmfund/api/db2"
	"gitlab.com/swarmfund/api/internal/data"
	"gitlab.com/swarmfund/api/internal/types"
)

const (
	blobsTable        = "blobs"
	blobsPKConstraint = "blobs_pkey"
)

var (
	ErrBlobsConflict = errors.New("blobs primary key conflict")
	blobsSelect      = squirrel.
				Select("id", "owner_address", "value", "type", "relationships", "deleted_at").
				From(blobsTable)
)

type Blobs struct {
	*db2.Repo
	stmt squirrel.SelectBuilder
}

func NewBlobs(repo *db2.Repo) *Blobs {
	return &Blobs{
		repo.Clone(), blobsSelect,
	}
}

func (q *Blobs) New() data.Blobs {
	return NewBlobs(q.Repo.Clone())
}

func (q *Blobs) Transaction(fn func(data.Blobs) error) error {
	return q.Repo.Transaction(func() error {
		return fn(q)
	})
}

func (q *Blobs) Delete(blobs ...types.Blob) error {
	if len(blobs) == 0 {
		return nil
	}
	ids := make([]string, 0, len(blobs))
	for _, blob := range blobs {
		ids = append(ids, blob.ID)
	}
	stmt := squirrel.Delete(blobsTable).Where(squirrel.Eq{"id": ids})
	_, err := q.Exec(stmt)
	return err
}

func (q *Blobs) MarkDeleted(id string) error {
	stmt := squirrel.
		Update(blobsTable).
		Where("id = ?", id).
		SetMap(map[string]interface{}{
			"deleted_at": time.Now().UTC(),
		})
	_, err := q.Exec(stmt)
	return err
}

func (q *Blobs) ByOwner(owner types.Address) data.Blobs {
	return &Blobs{
		q.Repo,
		q.stmt.Where("owner_address = ?", owner),
	}
}

func (q *Blobs) ByType(tpe types.BlobType) data.Blobs {
	return &Blobs{
		q.Repo,
		q.stmt.Where("type & ? != 0", tpe),
	}
}

func (q *Blobs) ExcludeDeleted() data.Blobs {
	return &Blobs{
		q.Repo,
		q.stmt.Where("deleted_at is null"),
	}
}

func (q *Blobs) ByRelationships(filters map[string]string) data.Blobs {
	builder := q.stmt
	for k, v := range filters {
		builder = builder.Where("relationships->>? = ?", k, v)
	}
	return &Blobs{
		q.Repo,
		builder,
	}
}

func (q *Blobs) Select() (result []types.Blob, err error) {
	err = q.Repo.Select(&result, q.stmt)
	return result, err
}

func (q *Blobs) Create(address *types.Address, blob *types.Blob) error {
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
				return ErrBlobsConflict
			}
		}
	}
	return err
}

func (q *Blobs) Get(id string) (*types.Blob, error) {
	var result types.Blob
	stmt := q.stmt.Where("id = ?", id)

	err := q.Repo.Get(&result, stmt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &result, nil
}
