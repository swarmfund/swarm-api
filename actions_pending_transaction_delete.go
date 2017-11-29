package api

import (
	"gitlab.com/swarmfund/api/db2/api"
	"gitlab.com/swarmfund/api/render/hal"
	"gitlab.com/swarmfund/api/render/problem"
)

type PendingTransactionDeleteAction struct {
	Action
	TxHash      string
	Transaction *api.PendingTransaction
}

// JSON format action handler
func (action *PendingTransactionDeleteAction) JSON() {
	action.ValidateBodyType()
	action.Do(
		action.loadParams,
		action.loadRecord,
		action.checkIsAllowed,
		action.performUpdate,
		func() {
			hal.Render(action.W, &problem.Success)
		})
}

func (action *PendingTransactionDeleteAction) loadParams() {
	action.TxHash = action.GetNonEmptyString("tx_hash")
}

func (action *PendingTransactionDeleteAction) loadRecord() {
	tx, err := action.APIQ().PendingTransactionByHash(action.TxHash)
	if err != nil {
		action.Log.WithError(err).Error("Failed to get pending tx")
		action.Err = &problem.ServerError
		return
	}

	if tx == nil {
		action.Err = &problem.NotFound
		return
	}

	action.Transaction = tx
}

func (action *PendingTransactionDeleteAction) checkIsAllowed() {
	action.checkSignerConstraints(
		SignerOf(action.Transaction.Source),
	)
}

func (action *PendingTransactionDeleteAction) performUpdate() {
	err := action.APIQ().PendingTransactions().Delete(action.Transaction)
	if err != nil {
		action.Log.WithError(err).Error("Failed to delete pending tx")
		action.Err = &problem.ServerError
		return
	}
}
