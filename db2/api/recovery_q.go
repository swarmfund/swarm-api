package api

import (
	"database/sql"
	"strings"
	"time"

	"gitlab.com/swarmfund/api/db2"
	"github.com/lann/squirrel"
)

var tableRecovery = "recovery_requests"
var insertRecovery = squirrel.Insert(tableRecovery)
var selectRecovery = squirrel.Select("*").From(tableRecovery)
var updateRecovery = squirrel.Update(tableRecovery)
var deleteRecovery = squirrel.Delete(tableRecovery)

type RecoveryQ struct {
	Err error
	sql squirrel.SelectBuilder
	db  *Q
}

type RecoveryQI interface {
	Create(request *RecoveryRequest) error
	Get(token string) (*RecoveryRequest, error)
	MarkSent(token string) error
	MarkCodeShown(token string) error
	MarkRejected(accountID string) error
	ByID(id int64) (*RecoveryRequest, error)
	ByWalletID(id int) (*RecoveryRequest, error)
	ByAccountID(accountID string) (*RecoveryRequest, error)
	ByUsername(username string) (*RecoveryRequest, error)
	Delete(accountID string) error
	SetRecoveryWalletID(id int64, walletID string) error

	Page(page db2.PageQuery) RecoveryQI
	Uploaded() RecoveryQI
	Select() ([]RecoveryRequest, error)
}

func (q *Q) Recoveries() RecoveryQI {
	return &RecoveryQ{
		db: q,
	}
}

func (q *RecoveryQ) Create(request *RecoveryRequest) error {
	stmt := insertRecovery. //Columns("wallet_ids", "email_token", "code").
				SetMap(map[string]interface{}{
			"account_id":  request.AccountID,
			"wallet_id":   request.WalletID,
			"email_token": request.EmailToken,
			"code":        request.Code,
			"username":    request.Username,
		})

	_, err := q.db.Exec(stmt)

	if err != nil {
		if strings.Contains(err.Error(), "recovery_requests_account_id_fkey") {
			// suppressing error so requester does not know of account existence
			return nil
		}
	}
	return err
}

func (q *RecoveryQ) SetRecoveryWalletID(id int64, walletID string) error {
	stmt := updateRecovery.
		Where("id = ?", id).
		Set("recovery_wallet_id", walletID).
		Set("uploaded_at", time.Now().UTC())
	_, err := q.db.Exec(stmt)
	return err
}

func (q *RecoveryQ) Get(token string) (*RecoveryRequest, error) {
	result := &RecoveryRequest{}
	stmt := selectRecovery.Where("email_token = ?", token)
	err := q.db.Get(result, stmt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
	}
	return result, err
}

func (q *RecoveryQ) ByID(id int64) (*RecoveryRequest, error) {
	result := &RecoveryRequest{}
	stmt := selectRecovery.Where("id = ?", id)
	err := q.db.Get(result, stmt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
	}
	return result, err
}

func (q *RecoveryQ) ByWalletID(id int) (*RecoveryRequest, error) {
	result := &RecoveryRequest{}
	stmt := selectRecovery.Where("wallet_id = ?", id)
	err := q.db.Get(result, stmt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
	}
	return result, err
}

func (q *RecoveryQ) ByUsername(username string) (*RecoveryRequest, error) {
	result := &RecoveryRequest{}
	stmt := selectRecovery.Where("username = ?", username)
	err := q.db.Get(result, stmt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
	}
	return result, err
}

func (q *RecoveryQ) ByAccountID(accountID string) (*RecoveryRequest, error) {
	result := &RecoveryRequest{}
	stmt := selectRecovery.Where("account_id = ?", accountID)
	err := q.db.Get(result, stmt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
	}
	return result, err
}

func (q *RecoveryQ) MarkSent(token string) error {
	stmt := updateRecovery.Set("sent_at", time.Now())
	_, err := q.db.Exec(stmt)
	return err
}

func (q *RecoveryQ) MarkCodeShown(token string) error {
	stmt := updateRecovery.Set("code_shown_at", time.Now())
	_, err := q.db.Exec(stmt)
	return err
}

func (q *RecoveryQ) MarkRejected(accountID string) error {
	stmt := updateRecovery.
		Set("uploaded_at", nil).
		Set("recovery_wallet_id", nil).
		Where("account_id = ?", accountID)
	_, err := q.db.Exec(stmt)
	return err
}

func (q *RecoveryQ) Delete(accountID string) error {
	stmt := deleteRecovery.Where("account_id = ?", accountID)
	_, err := q.db.Exec(stmt)
	return err

}

func (q *RecoveryQ) Page(page db2.PageQuery) RecoveryQI {
	if q.Err != nil {
		return q
	}

	q.sql, q.Err = page.ApplyTo(selectRecovery, "id")

	return q
}

func (q *RecoveryQ) Uploaded() RecoveryQI {
	if q.Err != nil {
		return q
	}

	q.sql = q.sql.Where("uploaded_at is not null")

	return q
}

func (q *RecoveryQ) Select() ([]RecoveryRequest, error) {
	if q.Err != nil {
		return nil, q.Err
	}

	result := []RecoveryRequest{}
	q.Err = q.db.Select(&result, q.sql)

	return result, q.Err
}
