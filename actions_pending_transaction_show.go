package api

import (
	"gitlab.com/swarmfund/api/db2/api"
	"gitlab.com/swarmfund/api/render/hal"
	"gitlab.com/swarmfund/api/resource"
)

type PendingTransactionShowAction struct {
	Action
	ID       int64
	Record   api.PendingTransaction
	Resource resource.PendingTransaction
}

func (action *PendingTransactionShowAction) loadParams() {
	action.ID = action.GetInt64("id")
}

func (action *PendingTransactionShowAction) loadRecord() {
	action.Err = action.APIQ().PendingTransactionByID(&action.Record, action.ID)
}

func (action *PendingTransactionShowAction) loadResource() {
	var signers []api.PendingTransactionSigner
	action.Err = action.APIQ().PendingTransactionSigners().ForTransaction(action.ID).Select(&signers)
	action.Resource.Populate(&action.Record)
	action.Resource.PopulateWithSigners(signers)
}

// JSON is a method for actions.JSON
func (action *PendingTransactionShowAction) JSON() {
	action.Do(
		action.loadParams,
		action.checkAllowed,
		action.loadRecord,
		action.loadResource,
		func() { hal.Render(action.W, action.Resource) },
	)
}

func (action *PendingTransactionShowAction) checkAllowed() {
	action.checkSignerConstraints(
		SignerOf(action.App.CoreInfo.MasterAccountID),
	)
}
