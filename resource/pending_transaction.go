package resource

import (
	"strconv"
	"time"

	"gitlab.com/swarmfund/api/db2/api"
)

type PendingTransaction struct {
	PT             string                     `json:"paging_token"`
	TxHash         string                     `json:"tx_hash"`
	TxEnvelope     string                     `json:"tx_envelope"`
	OperationType  string                     `json:"operation_type"`
	OperationTypeI int32                      `json:"operation_type_i"`
	State          int32                      `json:"state"`
	Signers        []PendingTransactionSigner `json:"signers"`
	CreatedAt      time.Time                  `json:"created_at"`
	UpdatedAt      time.Time                  `json:"updated_at"`
}

type PendingTransactionSigner struct {
	PublicKey string `json:"public_key"`
	Identity  int32  `json:"identity"`
	Name      string `json:"name"`
}

func (transaction *PendingTransaction) PopulateWithSigners(rows []api.PendingTransactionSigner) {
	for _, row := range rows {
		var signer PendingTransactionSigner
		signer.Identity = row.SignerIdentity
		signer.PublicKey = row.SignerPublicKey
		signer.Name = row.SignerName
		transaction.Signers = append(transaction.Signers, signer)
	}
}

// Populate fills out the resource's fields
func (transaction *PendingTransaction) Populate(row *api.PendingTransaction) {
	transaction.PT = strconv.FormatInt(row.ID, 10)
	transaction.TxHash = row.TxHash
	transaction.TxEnvelope = row.TxEnvelope
	transaction.OperationTypeI = row.OperationType
	transaction.State = row.State
	transaction.CreatedAt = row.CreatedAt
	transaction.UpdatedAt = row.UpdatedAt
}

// PagingToken implementation for hal.Pageable
func (request *PendingTransaction) PagingToken() string {
	return request.PT
}
