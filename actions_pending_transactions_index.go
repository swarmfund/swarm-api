package api

import (
	"strconv"

	"gitlab.com/swarmfund/api/db2"
	"gitlab.com/swarmfund/api/db2/api"
	"gitlab.com/swarmfund/api/render/hal"
	"gitlab.com/swarmfund/api/resource"
)

// PendingTransactionIndexAction returns a paged slice of admin transactions based upon the provided
// filters
type PendingTransactionIndexAction struct {
	Action

	AccountID string

	SignedByFilter    string
	NotSignedByFilter string
	StateFilter       int32
	PagingParams      db2.PageQuery
	Records           []api.PendingTransaction
	Page              hal.Page
}

// JSON is a method for actions.JSON
func (action *PendingTransactionIndexAction) JSON() {
	action.Do(
		action.ValidateCursorAsDefault,
		action.loadParams,
		action.checkAllowed,
		action.loadRecords,
		action.loadPage,
	)
	action.Do(func() {
		hal.Render(action.W, action.Page)
	})
}

func (action *PendingTransactionIndexAction) loadParams() {
	action.AccountID = action.GetNonEmptyString("id")
	action.StateFilter = action.GetInt32("state")
	action.SignedByFilter = action.GetString("signed_by")
	action.NotSignedByFilter = action.GetString("not_signed_by")
	action.PagingParams = action.GetPageQuery()
	action.Page.Filters = map[string]string{
		"state":         strconv.FormatInt(int64(action.StateFilter), 10),
		"signed_by":     action.SignedByFilter,
		"not_signed_by": action.NotSignedByFilter,
		"pending":       "true",
	}
}

func (action *PendingTransactionIndexAction) checkAllowed() {
	action.checkSignerConstraints(
		SignerOf(action.AccountID),
	)
}

func (action *PendingTransactionIndexAction) loadRecords() {
	q := action.APIQ().PendingTransactions()

	q = q.ForSource(action.AccountID)

	if action.StateFilter > 0 {
		q = q.ForState(action.StateFilter)
	}

	switch {
	case action.SignedByFilter != "":
		q.SignedBy(action.SignedByFilter)
	case action.NotSignedByFilter != "":
		q.NotSignedBy(action.NotSignedByFilter)
	}

	action.Err = q.Page(action.PagingParams).Select(&action.Records)
}

func (action *PendingTransactionIndexAction) loadPage() {
	for i := range action.Records {
		// FIXME woah! transaction already have signers
		var pendingTransaction resource.PendingTransaction
		var signers []api.PendingTransactionSigner
		action.Err = action.APIQ().PendingTransactionSigners().ForTransaction(action.Records[i].ID).Select(&signers)
		pendingTransaction.Populate(&action.Records[i])
		pendingTransaction.PopulateWithSigners(signers)
		action.Page.Add(&pendingTransaction)
	}

	action.Page.BaseURL = action.BaseURL()
	action.Page.BasePath = action.Path()
	action.Page.Limit = action.PagingParams.Limit
	action.Page.Cursor = action.PagingParams.Cursor
	action.Page.Order = action.PagingParams.Order
	action.Page.PopulateLinks()
}
