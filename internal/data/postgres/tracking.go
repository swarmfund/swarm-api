package postgres

import (
	"encoding/json"

	"database/sql"

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

func (t *Tracking) Last(query *track.Event) (*track.Event, error) {
	clauses := map[string]interface{}{}
	if query.Address != "" {
		clauses["address"] = query.Address
	}
	stmt := squirrel.Select("address", "signer", "details").
		From("tracking").
		Where(clauses).
		OrderBy("id desc").
		Limit(1)
	event := track.Event{}
	err := t.Get(&event, stmt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &event, err
}
