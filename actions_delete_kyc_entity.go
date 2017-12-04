package api

import (
	"gitlab.com/swarmfund/api/db2/api"
	"gitlab.com/swarmfund/api/render/hal"
	"gitlab.com/swarmfund/api/render/problem"
	"gitlab.com/swarmfund/go/xdr"
	horizon "gitlab.com/swarmfund/horizon-connector"
)

type DeleteKYCEntityAction struct {
	Action

	AccountID string
	EntityID  int64

	User *api.User
}

func (action *DeleteKYCEntityAction) JSON() {
	action.Do(
		action.loadParams,
		action.checkAllowed,
		action.loadUser,
		action.performRequest,
		action.blockAccount,
		action.updateState,
		func() {
			hal.Render(action.W, problem.Success)
		},
	)
}

func (action *DeleteKYCEntityAction) loadParams() {
	action.AccountID = action.GetNonEmptyString("user")
	action.EntityID = action.GetInt64("entity")
}

func (action *DeleteKYCEntityAction) checkAllowed() {
	action.checkSignerConstraints(
		SignerOf(action.AccountID),
	)
}

func (action *DeleteKYCEntityAction) loadUser() {
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

func (action *DeleteKYCEntityAction) performRequest() {
	for _, entity := range action.User.KYCEntities {
		if entity.ID == action.EntityID {
			if err := action.APIQ().Users().KYC().Delete(entity.ID); err != nil {
				action.Log.WithError(err).Error("failed to delete entity")
				action.Err = &problem.ServerError
				return
			}
			return
		}
	}
	action.Err = &problem.Forbidden
	return
}

func (action *DeleteKYCEntityAction) blockAccount() {
	account, err := action.App.horizon.AccountSigned(action.App.AccountManagerKP(), action.AccountID)
	if err != nil {
		action.Log.WithError(err).Error("failed to get user account")
		action.Err = &problem.ServerError
		return
	}

	if account == nil {
		action.Log.WithField("id", action.AccountID).Warn("account does not exist, but it should")
		action.Err = &problem.ServerError
		return
	}

	if xdr.AccountType(account.AccountType) == xdr.AccountTypeGeneral {
		// block user
		err = action.App.horizon.Transaction(&horizon.TransactionBuilder{Source: action.App.MasterKP()}).
			Op(&horizon.ManageAccountOp{
				AccountType: xdr.AccountType(account.AccountType),
				AccountID:   action.AccountID,
				AddReasons:  xdr.BlockReasonsKycUpdate,
			}).Sign(action.App.AccountManagerKP()).Submit()
		if err != nil {
			action.Log.WithError(err).Error("failed to submit block user tx")
			action.Err = &problem.ServerError
			return
		}
	}
}

func (action *DeleteKYCEntityAction) updateState() {
	state := action.User.CheckState()
	if state != action.User.State {
		err := action.APIQ().Users().ChangeState(string(action.User.Address), state)
		if err != nil {
			action.Log.WithError(err).Error("failed to update user state")
			action.Err = &problem.ServerError
			return
		}
	}
}
