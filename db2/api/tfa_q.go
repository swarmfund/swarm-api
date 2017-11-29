package api

import (
	"strings"

	"database/sql"

	"encoding/json"

	sq "github.com/lann/squirrel"
	"github.com/lib/pq"
	"github.com/pkg/errors"
	"gitlab.com/swarmfund/api/internal/types"
	"gitlab.com/swarmfund/api/tfa"
)

const (
	tfaTable        = "tfa"
	tfaBackendTable = "tfa_backends"

	tfaWalletBackendConstraint        = "tfa_backends_totp_constraint"
	tfaWalletFactorPasswordConstraint = "tfa_backends_password_constraint"
)

var (
	ErrWalletBackendConflict = errors.New("backend violated constraint")
)

//go:generate mockery -case underscore -name TFAQI
type TFAQI interface {
	New() TFAQI

	Get(token string) (*TFA, error)
	Create(tfa *TFA) error
	Consume(token string) (bool, error)
	Verify(token string) error

	Backend(id int64) (*Backend, error)
	Backends(walletID string) ([]Backend, error)
	DeleteBackends(wallet *Wallet) error
	DeleteBackend(id int64) error
	CreateBackend(walletID string, backend tfa.Backend) (*int64, error)
	SetBackendPriority(id int64, priority int) error
}

type TFAQ struct {
	parent *Q
}

func (q *Q) TFA() TFAQI {
	return &TFAQ{
		parent: q,
	}
}

func (q *TFAQ) New() TFAQI {
	return &TFAQ{
		parent: q.parent,
	}
}

func (q *TFAQ) Backends(walletID string) ([]Backend, error) {
	var records []Backend
	stmt := sq.Select("*").From(tfaBackendTable).
		Where("wallet_id = ?", walletID).
		OrderBy("priority desc")
	err := q.parent.Select(&records, stmt)
	return records, err
}

func (q *TFAQ) Backend(id int64) (*Backend, error) {
	var record Backend
	stmt := sq.Select("*").From(tfaBackendTable).
		Where("id = ?", id)
	err := q.parent.Get(&record, stmt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &record, err
}

func (q *TFAQ) CreateBackend(walletID string, backend tfa.Backend) (*int64, error) {
	model := Backend{
		WalletID: walletID,
		Priority: 0,
	}
	switch backend.(type) {
	case tfa.GoogleTOPT:
		model.BackendType = types.WalletFactorTOTP
		details, err := json.Marshal(&backend)
		if err != nil {
			return nil, errors.Wrap(err, "failed to marshal details")
		}
		model.Details = details
	}
	var id int64
	stmt := sq.Insert(tfaBackendTable).
		SetMap(map[string]interface{}{
			"wallet_id": model.WalletID,
			"details":   model.Details,
			"priority":  model.Priority,
			"backend":   model.BackendType,
		}).
		Suffix(`returning id`)

	err := q.parent.Get(&id, stmt)
	if err != nil {
		cause := errors.Cause(err)
		pqerr, ok := cause.(*pq.Error)
		if ok {
			if pqerr.Constraint == tfaWalletBackendConstraint {
				return nil, ErrWalletBackendConflict
			}
		}
	}

	return &id, err
}

func (q *TFAQ) Create(tfa *TFA) error {
	sql := sq.Insert(tfaTable).SetMap(map[string]interface{}{
		"otp_data": `{}`,
		"token":    tfa.Token,
		"backend":  tfa.BackendID,
	})

	_, err := q.parent.Exec(sql)
	if err != nil && strings.Contains(err.Error(), `"tfa_token_unique"`) {
		return nil
	}
	return err
}

func (q *TFAQ) Consume(token string) (bool, error) {
	sql := sq.Delete(tfaTable).Where(sq.Eq{
		"token":    token,
		"verified": true,
	})

	result, err := q.parent.Exec(sql)
	if err != nil {
		return false, err
	}

	affected, err := result.RowsAffected()
	return affected > 0, err
}

func (q *TFAQ) Verify(token string) error {
	sql := sq.Update(tfaTable).
		Where(sq.Eq{
			"token": token,
		}).
		Set("verified", true)

	_, err := q.parent.Exec(sql)
	return err
}

func (q *TFAQ) Get(token string) (*TFA, error) {
	var t TFA
	stmt := sq.Select("*").From(tfaTable).Where(sq.Eq{
		"token": token,
	})

	err := q.parent.Get(&t, stmt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
	}

	return &t, err
}

func (q *TFAQ) SetBackendPriority(id int64, priority int) error {
	stmt := sq.Update(tfaBackendTable).
		Where("id = ?", id).
		Set("priority", priority)
	_, err := q.parent.Exec(stmt)
	return err
}

func (q *TFAQ) DeleteBackend(id int64) error {
	stmt := sq.Delete(tfaBackendTable).Where("id = ?", id)
	_, err := q.parent.Exec(stmt)
	return err
}

func (q *TFAQ) DeleteBackends(wallet *Wallet) error {
	sql := sq.Delete(tfaBackendTable).Where(sq.Eq{
		"wallet_id": wallet.Id,
	})

	_, err := q.parent.Exec(sql)
	return err
}
