package api

import "gitlab.com/swarmfund/api/internal/types"

type Participant struct {
	OperationID int64         `db:"history_operation_id"`
	AccountID   types.Address `db:"account_id" json:"account_id"`
	BalanceID   string        `db:"balance_id" json:"balance_id"`
	Email       *string
	Effects     string
}
