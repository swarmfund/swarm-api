package api

import (
	"time"

	"database/sql"

	"github.com/lann/squirrel"
	"github.com/pkg/errors"
	"gitlab.com/swarmfund/api/db2"
	"gitlab.com/swarmfund/api/internal/data"
)

const (
	emailTokenTable        = "email_tokens"
	emailTokenTableAliased = "email_tokens et"
)

var (
	emailTokensSelect = squirrel.
		Select("et.*", "w.email").
		From(emailTokenTableAliased).
		Join("wallets w on w.wallet_id=et.wallet_id")
)

type EmailTokensQ struct {
	*db2.Repo
}

func NewEmailTokensQ(repo *db2.Repo) *EmailTokensQ {
	return &EmailTokensQ{
		Repo: repo,
	}
}

func (q *EmailTokensQ) New() data.EmailTokensQ {
	return NewEmailTokensQ(q.Repo)
}

func (q *EmailTokensQ) Create(wid, token string) error {
	stmt := squirrel.Insert(emailTokenTable).SetMap(map[string]interface{}{
		"wallet_id": wid,
		"token":     token,
	})
	_, err := q.Exec(stmt)
	return err
}

func (q *EmailTokensQ) Get(walletID string) (*data.EmailToken, error) {
	stmt := emailTokensSelect.
		Where("et.wallet_id = ?", walletID)
	result := data.EmailToken{}
	err := q.Repo.Get(&result, stmt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
	}
	return &result, err
}

func (q *EmailTokensQ) MarkUnsent(tid int64) error {
	stmt := squirrel.
		Update(emailTokenTable).
		Set("last_sent_at", nil).
		Where("id = ?", tid)
	_, err := q.Exec(stmt)
	return err
}

func (q *EmailTokensQ) Verify(walletID, token string) (bool, error) {
	stmt := squirrel.
		Update(emailTokenTable).
		Set("confirmed", true).
		Where("wallet_id = ? and token = ?", walletID, token)

	result, err := q.Exec(stmt)
	if err != nil {
		return false, errors.Wrap(err, "update failed")
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return false, errors.Wrap(err, "failed to get rows affected")
	}
	return rows > 0, nil
}

func (q *EmailTokensQ) GetUnsent() ([]data.EmailToken, error) {
	stmt := emailTokensSelect.
		Where("confirmed = false and last_sent_at is null")

	result := []data.EmailToken{}
	err := q.Select(&result, stmt)
	return result, err
}

func (q *EmailTokensQ) GetUnconfirmed() ([]data.EmailToken, error) {
	stmt := emailTokensSelect.
		Where("confirmed = false")

	result := []data.EmailToken{}
	err := q.Select(&result, stmt)
	return result, err
}

func (q *EmailTokensQ) MarkSent(tid int64) error {
	stmt := squirrel.
		Update(emailTokenTable).
		Set("last_sent_at", time.Now().UTC()).
		Where("id = ?", tid)
	_, err := q.Exec(stmt)
	return err
}
