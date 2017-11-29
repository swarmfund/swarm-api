package api

import (
	sq "github.com/lann/squirrel"
)

var selectPendingTransactionSigner = sq.Select("ps.*").From("pending_transaction_signers ps")
var insertPendingTransactionSigner = sq.Insert("pending_transaction_signers").Columns(
	"pending_transaction_id",
	"signer_identity",
	"signer_public_key",
	"signer_name",
)

// AccountsQ is a helper struct to aid in configuring queries that loads
// slices of account structs.
type PendingTransactionSignersQ struct {
	Err    error
	parent *Q
	sql    sq.SelectBuilder
}

type PendingTransactionSignersQI interface {
	Select(dest interface{}) error
	Create(transactionSigner *PendingTransactionSigner) error
	ForTransaction(transactionID int64) PendingTransactionSignersQI
}

// PendingTransactions provides a helper to filter rows from the `pending_transactions` table
// with pre-defined filters.  See `PendingTransactionSignersQI` methods for the available filters.
func (q *Q) PendingTransactionSigners() PendingTransactionSignersQI {
	return &PendingTransactionSignersQ{
		parent: q,
		sql:    selectPendingTransactionSigner,
	}
}

func (q *PendingTransactionSignersQ) ForTransaction(transactionID int64) PendingTransactionSignersQI {
	q.sql = q.sql.Where("ps.pending_transaction_id = ?", transactionID)
	return q
}

func (q *PendingTransactionSignersQ) Create(transactionSigner *PendingTransactionSigner) error {
	sql := insertPendingTransactionSigner.Values(
		transactionSigner.PendingTransactionID,
		transactionSigner.SignerIdentity,
		transactionSigner.SignerPublicKey,
		transactionSigner.SignerName,
	).RunWith(q.parent.DB)

	_, err := q.parent.Exec(sql)
	return err
}

// Select loads the results of the query specified by `q` into `dest`.
func (q *PendingTransactionSignersQ) Select(dest interface{}) error {
	if q.Err != nil {
		return q.Err
	}

	q.Err = q.parent.Select(dest, q.sql)
	return q.Err
}
