package api

import (
	"encoding/json"
	"fmt"
)

// PendingTransaction is a row of data from the `pending_transactions` table
type PendingTransactionSigner struct {
	ID                   int64
	PendingTransactionID int64  `db:"pending_transaction_id"`
	SignerIdentity       int32  `db:"signer_identity"`
	SignerPublicKey      string `db:"signer_public_key"`
	SignerName           string `db:"signer_name"`
}

type PendingTransactionSigners []PendingTransactionSigner

func (s *PendingTransactionSigners) Scan(value interface{}) error {
	if value == nil {
		return nil
	}
	switch v := value.(type) {
	case []byte:
		return json.Unmarshal(v, &s)
	default:
		return fmt.Errorf("unsupported Scan from type %T", v)
	}
}
