package api

import (
	"database/sql"
	"time"

	"gitlab.com/swarmfund/api/db2"
	sq "github.com/lann/squirrel"
)

var selectPendingTransaction = sq.
	Select("pt.*", "json_agg(pts order by pts.id) as signers").
	From("pending_transactions pt").
	LeftJoin("pending_transaction_signers pts on pt.id = pts.pending_transaction_id").
	GroupBy("pt.id")

var updatePendingTransaction = sq.Update("pending_transactions")

const (
	_ = iota
	PendingTxStatusPending
	_
	PendingTxStatusRejected
)

// AccountsQ is a helper struct to aid in configuring queries that loads
// slices of account structs.
type PendingTransactionsQ struct {
	Err    error
	parent *Q
	sql    sq.SelectBuilder
}

type PendingTransactionsQI interface {
	Page(page db2.PageQuery) PendingTransactionsQI
	Select(dest interface{}) error
	Update(transaction *PendingTransaction) error
	Delete(transaction *PendingTransaction) error
	ForState(state int32) PendingTransactionsQI
	ForSource(address string) PendingTransactionsQI
	NotSignedBy(accountID string) PendingTransactionsQI
	SignedBy(accountID string) PendingTransactionsQI
}

// PendingTransactions provides a helper to filter rows from the `pending_transactions` table
// with pre-defined filters.  See `PendingTransactionsQI` methods for the available filters.
func (q *Q) PendingTransactions() PendingTransactionsQI {
	return &PendingTransactionsQ{
		parent: q,
		sql:    selectPendingTransaction,
	}
}

func (q *Q) PendingTransactionByID(dest interface{}, id int64) error {
	sql := selectPendingTransaction.Limit(1).Where("pt.id = ?", id)
	return q.Get(dest, sql)
}

func (q *Q) PendingTransactionByHash(hash string) (*PendingTransaction, error) {
	tx := &PendingTransaction{}
	stmt := selectPendingTransaction.Limit(1).Where("pt.tx_hash = ?", hash)
	err := q.Get(tx, stmt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return tx, err
}

func (q *PendingTransactionsQ) ForSource(address string) PendingTransactionsQI {
	q.sql = q.sql.Where("source = ?", address)
	return q
}

func (q *PendingTransactionsQ) ForState(state int32) PendingTransactionsQI {
	q.sql = q.sql.Where("pt.state = ?", state)
	return q
}

func (q *PendingTransactionsQ) NotSignedBy(accountID string) PendingTransactionsQI {
	q.sql = q.sql.Where("pts.signer_public_key != ?", accountID)
	return q
}

func (q *PendingTransactionsQ) SignedBy(accountID string) PendingTransactionsQI {
	q.sql = q.sql.Where("pts.signer_public_key = ?", accountID)
	return q
}

func (q *PendingTransactionsQ) Update(transaction *PendingTransaction) error {
	sql := updatePendingTransaction.
		Set("tx_envelope", transaction.TxEnvelope).
		Set("state", transaction.State).
		Set("updated_at", time.Now().UTC()).
		Where("id = ?", transaction.ID)
	_, err := q.parent.Exec(sql)
	return err
}

func (q *PendingTransactionsQ) Delete(transaction *PendingTransaction) error {
	sql := sq.Delete("pending_transactions").
		Where("id = ?", transaction.ID)
	_, err := q.parent.Exec(sql)
	return err
}

// Page specifies the paging constraints for the query being built by `q`.
func (q *PendingTransactionsQ) Page(page db2.PageQuery) PendingTransactionsQI {
	if q.Err != nil {
		return q
	}

	q.sql, q.Err = page.ApplyTo(q.sql, "pt.id")
	return q
}

// Select loads the results of the query specified by `q` into `dest`.
func (q *PendingTransactionsQ) Select(dest interface{}) error {
	if q.Err != nil {
		return q.Err
	}

	q.Err = q.parent.Select(dest, q.sql)
	return q.Err
}
