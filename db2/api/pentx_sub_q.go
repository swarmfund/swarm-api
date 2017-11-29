package api

import (
	"database/sql"
	"fmt"

	"github.com/lann/squirrel"
)

const (
	pendingTransactionsTable        = "pending_transactions"
	pendingTransactionsSignersTable = "pending_transaction_signers"
	PENDING                         = 1
)

type PenTXSubQI interface {
	TransactionSigners(txID int64) ([]PendingTransactionSigner, error)
	CreateTransactionSigner(signer *PendingTransactionSigner) error
	TransactionByHash(hash string) (*PendingTransaction, error)
	UpdateTransaction(tx *PendingTransaction) error
	CreateTransaction(tx *PendingTransaction) (int64, error)
	DeleteTransaction(hash string) error
}

type PenTXSubQ struct {
	parent *Q
}

func (q *Q) PenTXSub() PenTXSubQI {
	return &PenTXSubQ{
		parent: q,
	}
}

func (q PenTXSubQ) TransactionSigners(txID int64) ([]PendingTransactionSigner, error) {
	var result []PendingTransactionSigner
	query := fmt.Sprintf(`select * from %s where pending_transaction_id = $1`, pendingTransactionsSignersTable)

	err := q.parent.DB.Select(&result, query, txID)
	return result, err
}

func (q PenTXSubQ) CreateTransactionSigner(signer *PendingTransactionSigner) error {
	query := fmt.Sprintf(`
		insert into %s (pending_transaction_id, signer_identity, signer_public_key, signer_name)
		values ($1, $2, $3, $4)`, pendingTransactionsSignersTable)

	_, err := q.parent.DB.Exec(query, signer.PendingTransactionID, signer.SignerIdentity, signer.SignerPublicKey, signer.SignerName)
	return err
}

func (q PenTXSubQ) TransactionByHash(hash string) (*PendingTransaction, error) {
	var tx PendingTransaction
	query := fmt.Sprintf(`select * from %s where tx_hash = $1 limit 1`, pendingTransactionsTable)
	err := q.parent.DB.Get(&tx, query, hash)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &tx, err
}

func (q PenTXSubQ) UpdateTransaction(tx *PendingTransaction) error {
	stmt, err := q.parent.DB.Prepare(fmt.Sprintf(`
		update %s
		set tx_envelope = $1,
		    state = $2,
		    updated_at = timestamp 'now'
		where id = $3`, pendingTransactionsTable))
	if err != nil {
		return err
	}
	_, err = stmt.Exec(tx.TxEnvelope, tx.State, tx.ID)
	return err
}

func (q PenTXSubQ) CreateTransaction(tx *PendingTransaction) (int64, error) {
	query := fmt.Sprintf(`
		insert into %s (
		  "tx_hash", "tx_envelope", "operation_type", "state", "created_at",
		  "updated_at", "source")
		values ($1,$2,$3,$4, timestamp 'now', timestamp 'now', $5) returning id
	`, pendingTransactionsTable)
	var id int64
	err := q.parent.DB.Get(
		&id, query, tx.TxHash, tx.TxEnvelope, tx.OperationType, PENDING, tx.Source)
	return id, err
}

func (q PenTXSubQ) DeleteTransaction(hash string) error {
	sql := squirrel.Delete("pending_transactions").
		Where("tx_hash = ?", hash)
	_, err := q.parent.Exec(sql)
	return err
}
