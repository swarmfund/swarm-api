package api

import (
	"gitlab.com/swarmfund/api/db2/api"
	"gitlab.com/swarmfund/api/internal/types"
	"gitlab.com/swarmfund/api/render/hal"
	"gitlab.com/swarmfund/api/render/problem"
	"gitlab.com/swarmfund/go/xdr"
	"gitlab.com/swarmfund/horizon-connector"
)

type CreateUserRequest struct {
	Address  string       `json:"address" valid:"required"`
	UserType api.UserType `json:"type" valid:"required"`
}

// CreateUserAction posts registration request to registration_requests table
type CreateUserAction struct {
	Action

	Request CreateUserRequest

	Wallet *api.Wallet
}

// JSON format action handler
func (action *CreateUserAction) JSON() {
	action.Do(
		action.ValidateBodyType,
		action.loadParams,
		action.checkAllowed,
		action.checkWallet,
		action.checkExists,
		action.ensureAccount,
		action.createUser,
		func() {
			hal.Render(action.W, problem.Success)
		})
}

func (action *CreateUserAction) loadParams() {
	action.UnmarshalBody(&action.Request)
	if err := action.Request.UserType.Validate(); err != nil {
		action.SetInvalidField("type", err)
	}
}

func (action *CreateUserAction) checkAllowed() {
	action.checkSignerConstraints(
		SignerOf(action.Request.Address),
	)
}

func (action *CreateUserAction) checkWallet() {
	// wallet should exists and be verified when creating user
	wallet, err := action.APIQ().Wallet().ByAccountID(types.Address(action.Request.Address))
	if err != nil {
		action.Log.WithError(err).Error("failed to get wallet")
		action.Err = &problem.ServerError
		return
	}

	if wallet == nil {
		action.Err = &problem.BadRequest
		return
	}

	if !wallet.Verified {
		action.Err = &problem.Forbidden
		return
	}

	action.Wallet = wallet
}

func (action *CreateUserAction) checkExists() {
	user, err := action.APIQ().Users().ByAddress(action.Request.Address)
	if err != nil {
		action.Log.WithError(err).Error("failed to get user")
		action.Err = &problem.ServerError
		return
	}

	if user != nil {
		action.Err = &problem.Conflict
		return
	}
}

func (action *CreateUserAction) ensureAccount() {
	tx, err := action.App.horizon.Transaction(&horizon.TransactionBuilder{Source: action.App.MasterKP()}).
		Op(&horizon.CreateAccountOp{
			AccountID:   action.Request.Address,
			AccountType: xdr.AccountTypeNotVerified,
		}).Sign(action.App.AccountManagerKP()).Marshal64()

	if err != nil {
		action.Log.WithError(err).Error("failed to marshal tx")
		action.Err = &problem.ServerError
		return
	}

	result, err := action.App.horizon.SubmitTXVerbose(*tx)
	if err != nil {
		action.Log.WithError(err).Error("failed to submit tx")
		action.Err = &problem.ServerError
		return
	}

	var xdrResult xdr.TransactionResult
	err = xdr.SafeUnmarshalBase64(result.Result, &xdrResult)
	if err != nil {
		action.Log.WithError(err).Error("failed to unmarshal tx result")
		action.Err = &problem.ServerError
		return
	}
}

func (action *CreateUserAction) createUser() {
	var err error
	user := api.User{
		Address:  types.Address(action.Request.Address),
		Email:    action.Wallet.Username,
		UserType: action.Request.UserType,
		State:    api.UserNeedDocs,
	}

	err = action.APIQ().Users().Create(&user)
	if err == api.ErrUsersConflict {
		action.Err = &problem.Conflict
		return
	}
	if err != nil {
		action.Log.WithError(err).Error("Failed to put user into db")
		action.Err = &problem.ServerError
		return
	}
}
