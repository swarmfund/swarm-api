package skrill

import (
	"strings"
	"time"
)

type TransactionType int
type TransactionStatus int

const (
	_                             = iota
	TxTypeReceive TransactionType = iota
	TxTypeSend
	TxTypeCancellation
)

const (
	TxStatusFailed    TransactionStatus = -2
	TxStatusCancelled TransactionStatus = -1
	TxStatusProcessed TransactionStatus = 2
	TxStatusScheduled TransactionStatus = 1
	TxStatusPending   TransactionStatus = 0
)

var (
	SkrillTypeTransactionType map[string]TransactionType = map[string]TransactionType{
		"Receive Money":           TxTypeReceive,
		"Send Money":              TxTypeSend,
		"Send Money Cancellation": TxTypeCancellation,
	}
	SkrillStatusTransactionStatus map[string]TransactionStatus = map[string]TransactionStatus{
		"-2":        TxStatusFailed,
		"-1":        TxStatusCancelled,
		"0":         TxStatusPending,
		"2":         TxStatusProcessed,
		"processed": TxStatusProcessed,
		"scheduled": TxStatusScheduled,
		"cancelled": TxStatusCancelled,
	}
)

type HistoryTransaction struct {
	ID            string
	Time          time.Time
	Type          TransactionType
	Details       string
	SendUSD       int64
	ReceivedUSD   int64
	Status        TransactionStatus
	Balance       int64
	Reference     string // transaction_id on session create
	AmountSent    int64
	CurrencySent  string
	Info          string
	TransactionID string
	PaymentType   string
}

func (tx *HistoryTransaction) From() *string {
	if !strings.HasPrefix(tx.Details, "from ") {
		return nil
	}
	split := strings.Split(tx.Details, " ")
	return &split[1]
}

func (tx *HistoryTransaction) IsFee() bool {
	return tx.Details == "Fee" || tx.Details == "Per Transaction Fee"
}

func (tx *HistoryTransaction) Cursor() string {
	return tx.Time.Format(time.RFC822)
}
