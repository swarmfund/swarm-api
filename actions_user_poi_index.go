package api

import (
	"gitlab.com/swarmfund/api/db2/api"
	"gitlab.com/swarmfund/api/render/hal"
	"gitlab.com/swarmfund/api/render/problem"
)

type UserProofOfIncomeIndexAction struct {
	Action

	AccountID string
	User      *api.User
	Records   []interface{}
}

func (action *UserProofOfIncomeIndexAction) JSON() {
	action.Do(
		action.ValidateBodyType,
		action.loadParams,
		action.checkAllowed,
		action.loadUser,
		action.loadRecords,
		func() {
			hal.Render(action.W, action.Records)
		},
	)
}

func (action *UserProofOfIncomeIndexAction) checkAllowed() {
	action.checkSignerConstraints(
		SignerOf(action.App.CoreInfo.MasterAccountID),
	)
}

func (action *UserProofOfIncomeIndexAction) loadParams() {
	action.AccountID = action.GetNonEmptyString("id")
}

func (action *UserProofOfIncomeIndexAction) loadUser() {
	user, err := action.APIQ().Users().ByAddress(action.AccountID)
	if err != nil {
		action.Log.WithError(err).Error("failed to get user")
		action.Err = &problem.ServerError
		return
	}

	if user == nil {
		action.Err = &problem.NotFound
		return
	}

	action.User = user
}

func (action *UserProofOfIncomeIndexAction) loadRecords() {
	action.Records = []interface{}{}

	documents, ok := action.User.Documents[api.DocumentTypeProofOfIncome]
	if !ok {
		return
	}

	for _, document := range documents {
		if raw, ok := document.Meta["reviewed"]; ok {
			if reviewed, ok := raw.(bool); ok {
				if !reviewed {
					action.Records = append(action.Records, document)
				}
			}
		}
	}

}
