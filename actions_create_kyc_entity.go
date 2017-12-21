package api

import (
	"github.com/pkg/errors"
	"gitlab.com/swarmfund/api/db2/api"
	"gitlab.com/swarmfund/api/render/hal"
	"gitlab.com/swarmfund/api/render/problem"
	"gitlab.com/swarmfund/api/resource"
	"gitlab.com/swarmfund/go/xdr"
	horizon "gitlab.com/swarmfund/horizon-connector"
)

type CreateKYCEntityRequest struct {
	Type api.KYCEntityType `json:"type" valid:"required"`
}

type CreateKYCEntityAction struct {
	Action

	AccountID string
	Request   CreateKYCEntityRequest

	User *api.User

	Resource resource.KYCEntity
}

func (action *CreateKYCEntityAction) JSON() {
	action.Do(
		action.ValidateBodyType,
		action.loadParams,
		action.checkAllowed,
		action.loadUser,
		action.performRequest,
		action.blockAccount,
		action.updateState,
		func() {
			hal.Render(action.W, action.Resource)
		},
	)
}

func (action *CreateKYCEntityAction) loadParams() {
	action.UnmarshalBody(&action.Request)
	action.AccountID = action.GetNonEmptyString("user")
}

func (action *CreateKYCEntityAction) checkAllowed() {
	action.checkSignerConstraints(
		SignerOf(action.AccountID),
	)
}

func (action *CreateKYCEntityAction) loadUser() {
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

func (action *CreateKYCEntityAction) performRequest() {
	//var record interface{}
	switch action.Request.Type {
	//case api.KYCEntityTypeJointIdentity:
	//	if action.User.UserType != api.UserTypeJoint {
	//		action.Err = &problem.Forbidden
	//		return
	//	}
	//	record = api.IdentityDetails{}
	//case api.KYCEntityTypeBusinessOwner, api.KYCEntityTypeBusinessSignatory:
	//	if action.User.UserType != api.UserTypeBusiness {
	//		action.Err = &problem.Forbidden
	//		return
	//	}
	//	record = api.BusinessPerson{}
	default:
		action.SetInvalidField("type", errors.New("invalid"))
		return
	}
	//data, err := json.Marshal(&record)
	//if err != nil {
	//	action.Log.WithError(err).WithField("type", action.Request.Type).Error("failed to marshal")
	//	action.Err = &problem.ServerError
	//	return
	//}
	//entity := api.KYCEntity{
	//	Data:   data,
	//	UserID: action.User.ID,
	//	Type:   action.Request.Type,
	//}
	//eid, err := action.APIQ().Users().KYC().Create(entity)
	//if err != nil {
	//	action.Log.WithError(err).Error("failed to save entity")
	//	action.Err = &problem.ServerError
	//	return
	//}
	//action.Resource = resource.KYCEntity{
	//	ID: eid,
	//}
}

func (action *CreateKYCEntityAction) blockAccount() {
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

func (action *CreateKYCEntityAction) updateState() {
	state := action.User.CheckState()
	if state != action.User.State {
		//err := action.APIQ().Users().ChangeState(string(action.User.Address), state)
		//if err != nil {
		//	action.Log.WithError(err).Error("failed to update user state")
		//	action.Err = &problem.ServerError
		//	return
		//}
	}
}
