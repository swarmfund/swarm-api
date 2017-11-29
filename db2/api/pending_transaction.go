package api

import (
	"time"
)

// PendingTransaction is a row of data from the `pending_transactions` table
type PendingTransaction struct {
	ID            int64                     `db:"id"`
	TxHash        string                    `db:"tx_hash"`
	TxEnvelope    string                    `db:"tx_envelope"`
	OperationType int32                     `db:"operation_type"`
	State         int32                     `db:"state"`
	CreatedAt     time.Time                 `db:"created_at"`
	UpdatedAt     time.Time                 `db:"updated_at"`
	OperationKey  *string                   `db:"operation_key"`
	Signers       PendingTransactionSigners `db:"signers"`
	Source        string                    `db:"source"`
}

func NewPendingTransaction(opType int32, envelope string, hash string, source string) *PendingTransaction {
	tx := &PendingTransaction{
		OperationType: opType,
		TxEnvelope:    envelope,
		TxHash:        hash,
		Source:        source,
	}
	return tx
}
