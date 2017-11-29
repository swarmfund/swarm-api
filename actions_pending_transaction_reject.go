package api

import (
	"gitlab.com/swarmfund/api/db2/api"
	"gitlab.com/swarmfund/api/render/hal"
	"gitlab.com/swarmfund/api/render/problem"
)

type PendingTransactionRejectAction struct {
	Action

	TX *api.PendingTransaction

	TxHash string
	Result *problem.P
}

// JSON format action handler
func (action *PendingTransactionRejectAction) JSON() {
	action.ValidateBodyType()
	action.Do(
		action.loadParams,
		action.loadTX,
		action.checkIsAllowed,
		action.performUpdate,
		func() {
			hal.Render(action.W, action.Result)
		})
}

func (action *PendingTransactionRejectAction) loadParams() {
	action.TxHash = action.GetNonEmptyString("tx_hash")
}

func (action *PendingTransactionRejectAction) loadTX() {
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

	if tx.State != api.PendingTxStatusPending {
		action.Err = &problem.AdminOperationRejectRestrictionViolated
		return
	}

	action.TX = tx
}

func (action *PendingTransactionRejectAction) checkIsAllowed() {
	action.checkSignerConstraints(
		SignerOf(action.TX.Source),
	)
}

func (action *PendingTransactionRejectAction) performUpdate() {
	action.TX.State = api.PendingTxStatusRejected
	err := action.APIQ().PendingTransactions().Update(action.TX)
	if err != nil {
		action.Log.WithError(err).Error("Failed to update pending tx")
		action.Err = &problem.ServerError
		return
	}
	action.Result = &problem.Success
}
