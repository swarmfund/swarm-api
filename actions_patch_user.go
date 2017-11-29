package api

import (
	"encoding/json"

	"gitlab.com/swarmfund/api/db2/api"
	"gitlab.com/swarmfund/api/render/hal"
	"gitlab.com/swarmfund/api/render/problem"
	"gitlab.com/swarmfund/go/xdr"
	horizon "gitlab.com/swarmfund/horizon-connector"
)

type PatchUserRequest struct {
	Details *struct {
		PersonalDetails       *api.PersonDetails      `json:"personal_details" valid:"optional"`
		EmploymentDetails     *api.EmploymentDetails  `json:"employment_details" valid:"optional"`
		BankDetails           *api.BankDetails        `json:"bank_details" valid:"optional"`
		CorporationDetails    *api.CorporationDetails `json:"corporation_details" valid:"optional"`
		FinancialDetails      *api.FinancialDetails   `json:"financial_details" valid:"optional"`
		CorrespondenceAddress *api.Address            `json:"correspondence_address" valid:"optional"`
		RegisteredAddress     *api.Address            `json:"registered_address" valid:"optional"`
		Address               *api.Address            `json:"address" valid:"optional"`
	} `json:"details" valid:"optional"`
}

type PatchUserAction struct {
	Action

	AccountID string
	Request   PatchUserRequest

	User *api.User
}

func (action *PatchUserAction) JSON() {
	action.Do(
		action.ValidateBodyType,
		action.loadParams,
		action.checkAllowed,
		action.loadUser,
		action.patchDetails,
		action.updateState,
		func() {
			hal.Render(action.W, &problem.Success)
		},
	)
}

func (action *PatchUserAction) loadParams() {
	action.UnmarshalBody(&action.Request)
	action.AccountID = action.GetNonEmptyString("id")
}

func (action *PatchUserAction) checkAllowed() {
	action.checkSignerConstraints(
		SignerOf(action.AccountID),
	)
}

func (action *PatchUserAction) loadUser() {
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

func (action *PatchUserAction) patchDetails() {
	if action.Request.Details == nil {
		return
	}

	// TODO begin tx

	details := action.Request.Details
	update := map[api.KYCEntityType]interface{}{}

	switch action.User.UserType {
	case api.UserTypeIndividual:
		rejectReasons := api.IndividualRejectReasons{}
		entity := action.User.KYCEntities.GetSingle(api.KYCEntityTypeIndividualRejectReasons)
		rejectReasons.Populate(entity)
		if details.BankDetails != nil {
			update[api.KYCEntityTypeBankDetails] = details.BankDetails
			rejectReasons.BankDetails = api.BankDetails{}
		}
		if details.EmploymentDetails != nil {
			update[api.KYCEntityTypeEmploymentDetails] = details.EmploymentDetails
			rejectReasons.EmploymentDetails = api.EmploymentDetails{}
		}
		if details.PersonalDetails != nil {
			update[api.KYCEntityTypePersonalDetails] = details.PersonalDetails
			rejectReasons.PersonalDetails = api.PersonDetails{}
		}
		if details.Address != nil {
			update[api.KYCEntityTypeAddress] = details.Address
			rejectReasons.Address = api.Address{}
		}
		update[api.KYCEntityTypeIndividualRejectReasons] = rejectReasons
	case api.UserTypeJoint:
	case api.UserTypeBusiness:
		rejectReasons := api.BusinessRejectReasons{}
		entity := action.User.KYCEntities.GetSingle(api.KYCEntityTypeBusinessRejectReasons)
		rejectReasons.Populate(entity)
		if details.FinancialDetails != nil {
			update[api.KYCEntityTypeFinancialDetails] = details.FinancialDetails
			rejectReasons.FinancialDetails = api.FinancialDetails{}
		}
		if details.RegisteredAddress != nil {
			update[api.KYCEntityTypeRegisteredAddress] = details.RegisteredAddress
			rejectReasons.RegisteredAddress = api.Address{}
		}
		if details.CorrespondenceAddress != nil {
			update[api.KYCEntityTypeCorrespondenceAddress] = details.CorrespondenceAddress
			rejectReasons.CorrespondenceAddress = api.Address{}
		}
		if details.CorporationDetails != nil {
			update[api.KYCEntityTypeCorporationDetails] = details.CorporationDetails
			rejectReasons.CorporationDetails = api.CorporationDetails{}
		}
		update[api.KYCEntityTypeBusinessRejectReasons] = rejectReasons
	}

	if len(update) == 0 {
		return
	}

	for entityType, strct := range update {
		data, err := json.Marshal(&strct)
		if err != nil {
			action.Log.WithError(err).WithField("type", entityType).Error("failed to marshal")
			action.Err = &problem.ServerError
			return
		}
		if eid, ok := action.User.KYCEntities.Exists(entityType); ok {
			if err = action.APIQ().Users().KYC().Update(eid, data); err != nil {
				action.Log.WithError(err).Error("failed to update entity")
				action.Err = &problem.ServerError
				return
			}
		} else {
			entity := api.KYCEntity{
				Data:   data,
				UserID: action.User.ID,
				Type:   entityType,
			}
			if _, err = action.APIQ().Users().KYC().Create(entity); err != nil {
				action.Log.WithError(err).Error("failed to save entity")
				action.Err = &problem.ServerError
				return
			}
		}
	}
}

func (action *PatchUserAction) blockAccount() {
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

	if account.AccountType == xdr.AccountTypeGeneral {
		// block user
		err = action.App.horizon.Transaction(&horizon.TransactionBuilder{Source: action.App.MasterKP()}).
			Op(&horizon.ManageAccountOp{
				AccountType: account.AccountType,
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

func (action *PatchUserAction) updateState() {
	state := action.User.CheckState()
	if state != action.User.State {
		err := action.APIQ().Users().ChangeState(string(action.User.Address), state)
		if err != nil {
			action.Log.WithError(err).Error("failed to update user state")
			action.Err = &problem.ServerError
			return
		}
	}

	// TODO commit tx
}
