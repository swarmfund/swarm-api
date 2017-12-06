package api

import (
	"encoding/json"

	"fmt"

	"gitlab.com/swarmfund/api/db2/api"
	"gitlab.com/swarmfund/api/render/hal"
	"gitlab.com/swarmfund/api/render/problem"
	"gitlab.com/swarmfund/go/xdr"
	horizon "gitlab.com/swarmfund/horizon-connector"
)

type PatchKYCEntityRequest struct {
	PersonalDetails   *api.PersonDetails     `json:"personal_details" valid:"optional"`
	EmploymentDetails *api.EmploymentDetails `json:"employment_details" valid:"optional"`
	Address           *api.Address           `json:"address" valid:"optional"`
	BankDetails       *api.BankDetails       `json:"bank_details" valid:"optional"`
}

type PatchKYCEntityAction struct {
	Action

	AccountID string
	EntityID  int64
	Request   PatchKYCEntityRequest

	User   *api.User
	Entity *api.KYCEntity
}

func (action *PatchKYCEntityAction) JSON() {
	action.Do(
		action.ValidateBodyType,
		action.loadParams,
		action.checkAllowed,
		action.loadUser,
		action.getEntity,
		action.performRequest,
		action.blockAccount,
		action.updateState,
		func() {
			hal.Render(action.W, problem.Success)
		},
	)
}

func (action *PatchKYCEntityAction) loadParams() {
	action.UnmarshalBody(&action.Request)
	action.AccountID = action.GetNonEmptyString("user")
	action.EntityID = action.GetInt64("entity")
}

func (action *PatchKYCEntityAction) checkAllowed() {
	action.checkSignerConstraints(
		SignerOf(action.AccountID),
	)
}

func (action *PatchKYCEntityAction) loadUser() {
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

func (action *PatchKYCEntityAction) getEntity() {
	for _, entity := range action.User.KYCEntities {
		if entity.ID == action.EntityID {
			action.Entity = &entity
			return
		}
	}
	action.Err = &problem.Forbidden
	return
}

func (action *PatchKYCEntityAction) performRequest() {
	var rejectReasonsEntity *api.KYCEntity
	eid := fmt.Sprintf("%d", action.Entity.ID)
	switch action.Entity.Type {
	case api.KYCEntityTypeJointIdentity:
		var rejectReasons api.JointRejectReasons
		identity, err := action.Entity.JointIdentity()
		if err != nil {
			action.Log.
				WithError(err).
				WithField("entity", action.Entity.ID).
				Error("failed to unmarshal")
			action.Err = &problem.ServerError
			return
		}
		rejectReasonsEntity, rejectReasons = action.User.KYCEntities.JointRejectReasons()
		if action.Request.PersonalDetails != nil {
			identity.PersonalDetails = *action.Request.PersonalDetails
			identityRR, ok := rejectReasons.IdentityDetails[eid]
			if ok {
				identityRR.PersonalDetails = api.PersonDetails{}
				rejectReasons.IdentityDetails[eid] = identityRR
			}
		}
		if action.Request.EmploymentDetails != nil {
			identity.EmploymentDetails = *action.Request.EmploymentDetails
			identityRR, ok := rejectReasons.IdentityDetails[eid]
			if ok {
				identityRR.EmploymentDetails = api.EmploymentDetails{}
				rejectReasons.IdentityDetails[eid] = identityRR
			}
		}
		if action.Request.Address != nil {
			identity.Address = *action.Request.Address
			identityRR, ok := rejectReasons.IdentityDetails[eid]
			if ok {
				identityRR.Address = api.Address{}
				rejectReasons.IdentityDetails[eid] = identityRR
			}
		}
		if action.Request.BankDetails != nil {
			identity.BankDetails = *action.Request.BankDetails
			identityRR, ok := rejectReasons.IdentityDetails[eid]
			if ok {
				identityRR.BankDetails = api.BankDetails{}
				rejectReasons.IdentityDetails[eid] = identityRR
			}
		}
		data, err := json.Marshal(&identity)
		if err != nil {
			action.Log.WithError(err).Error("failed to marshal entity")
			action.Err = &problem.ServerError
			return
		}
		action.Entity.Data = data

		if rejectReasonsEntity != nil {
			data, err = json.Marshal(&rejectReasons)
			if err != nil {
				action.Log.WithError(err).Error("failed to marshal entity")
				action.Err = &problem.ServerError
				return
			}
			rejectReasonsEntity.Data = data
		}
	case api.KYCEntityTypeBusinessSignatory:
		var rejectReasons api.BusinessRejectReasons
		person, err := action.Entity.BusinessPerson()
		if err != nil {
			action.Log.
				WithError(err).
				WithField("entity", action.Entity.ID).
				Error("failed to unmarshal")
			action.Err = &problem.ServerError
			return
		}
		rejectReasonsEntity, rejectReasons = action.User.KYCEntities.BusinessRejectReasons()
		if action.Request.PersonalDetails != nil {
			person.PersonDetails = *action.Request.PersonalDetails
			ownerRR, ok := rejectReasons.Signatories[eid]
			if ok {
				ownerRR.PersonDetails = api.PersonDetails{}
				rejectReasons.Signatories[eid] = ownerRR
			}
		}
		if action.Request.Address != nil {
			person.Address = *action.Request.Address
			ownerRR, ok := rejectReasons.Signatories[eid]
			if ok {
				ownerRR.Address = api.Address{}
				rejectReasons.Signatories[eid] = ownerRR
			}
		}
		data, err := json.Marshal(&person)
		if err != nil {
			action.Log.WithError(err).Error("failed to marshal entity")
			action.Err = &problem.ServerError
			return
		}
		action.Entity.Data = data

		if rejectReasonsEntity != nil {
			data, err = json.Marshal(&rejectReasons)
			if err != nil {
				action.Log.WithError(err).Error("failed to marshal entity")
				action.Err = &problem.ServerError
				return
			}
			rejectReasonsEntity.Data = data
		}
	case api.KYCEntityTypeBusinessOwner:
		var rejectReasons api.BusinessRejectReasons
		person, err := action.Entity.BusinessPerson()
		if err != nil {
			action.Log.
				WithError(err).
				WithField("entity", action.Entity.ID).
				Error("failed to unmarshal")
			action.Err = &problem.ServerError
			return
		}
		rejectReasonsEntity, rejectReasons = action.User.KYCEntities.BusinessRejectReasons()
		if action.Request.PersonalDetails != nil {
			person.PersonDetails = *action.Request.PersonalDetails
			ownerRR, ok := rejectReasons.Owners[eid]
			if ok {
				ownerRR.PersonDetails = api.PersonDetails{}
				rejectReasons.Owners[eid] = ownerRR
			}
		}

		if action.Request.Address != nil {
			person.Address = *action.Request.Address
			ownerRR, ok := rejectReasons.Owners[eid]
			if ok {
				ownerRR.Address = api.Address{}
				rejectReasons.Owners[eid] = ownerRR
			}
		}
		data, err := json.Marshal(&person)
		if err != nil {
			action.Log.WithError(err).Error("failed to marshal entity")
			action.Err = &problem.ServerError
			return
		}
		action.Entity.Data = data

		if rejectReasonsEntity != nil {
			data, err = json.Marshal(&rejectReasons)
			if err != nil {
				action.Log.WithError(err).Error("failed to marshal entity")
				action.Err = &problem.ServerError
				return
			}
			rejectReasonsEntity.Data = data
		}
	default:
		panic("wrong state")
	}

	err := action.APIQ().Users().KYC().Update(action.Entity.ID, action.Entity.Data)
	if err != nil {
		action.Log.WithError(err).Error("failed to update entity")
		action.Err = &problem.ServerError
		return
	}

	if rejectReasonsEntity != nil {
		err = action.APIQ().Users().KYC().Update(rejectReasonsEntity.ID, rejectReasonsEntity.Data)
		if err != nil {
			action.Log.WithError(err).Error("failed to update entity")
			action.Err = &problem.ServerError
			return
		}
	}
}

func (action *PatchKYCEntityAction) updateState() {
	state := action.User.CheckState()
	if state != action.User.State {
		//err := action.APIQ().Users().ChangeState(string(action.User.Address), state)
		//if err != nil {
		//	action.Log.WithError(err).Error("failed to change user state")
		//	action.Err = &problem.ServerError
		//	return
		//}
	}
}

func (action *PatchKYCEntityAction) blockAccount() {
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
