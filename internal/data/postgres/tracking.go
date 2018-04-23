package postgres

import (
	"encoding/json"

	"github.com/lann/squirrel"
	"github.com/pkg/errors"
	"gitlab.com/swarmfund/api/db2"
	"gitlab.com/swarmfund/api/internal/track"
)

type Tracking struct {
	*db2.Repo
}

func NewTracking(repo *db2.Repo) *Tracking {
	return &Tracking{
		repo.Clone(),
	}
}

func (t *Tracking) Track(event track.Event) error {
	details, err := json.Marshal(event.Details)
	if err != nil {
		return errors.Wrap(err, "failed to marshal details")
	}
	stmt := squirrel.Insert("tracking").SetMap(map[string]interface{}{
		"address": event.Address,
		"signer":  event.Signer,
		"details": details,
	})
	_, err = t.Exec(stmt)
	return err
}
