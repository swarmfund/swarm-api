package api

import (
	"encoding/json"

	"fmt"

	"github.com/pkg/errors"
	"gitlab.com/swarmfund/api/internal/types"
)

type RequiredDocument struct {
	Type     DocumentType
	EntityID int64
}

func (user *User) CheckState() types.UserState {
	switch user.UserType {
	case types.UserTypeNotVerified:
		return types.UserStateNil
	case types.UserTypeGeneral:
		// TODO unstub after proper KYC deployed
		return types.UserStateWaitingForApproval

		// should have individual details
		for _, record := range user.KYCEntities {
			if record.Entity.Type == types.KYCEntityTypeIndividual {
				if err := record.Entity.Individual.Validate(); err != nil {
					return types.UserStateRejected
				}
				if user.State == types.UserStateApproved {
					return types.UserStateApproved
				}
				return types.UserStateWaitingForApproval
			}
		}
	case types.UserTypeSyndicate:
		// TODO should have syndicate details
		if user.State == types.UserStateApproved {
			return types.UserStateApproved
		}
		return types.UserStateWaitingForApproval

	default:
		panic(fmt.Errorf("unknown user type %s", user.UserType))
	}

	return types.UserStateWaitingForApproval
}

type CorporationDetails struct {
	EntityName          string `json:"entity_name" valid:"required"`
	DateOfEstablishment string `json:"date_of_establishment" valid:"required"`
	Email               string `json:"email" valid:"required"`
	Landline            string `json:"landline" valid:"required"`
	Mobile              string `json:"mobile"`
	Fax                 string `json:"fax"`
	BusinessNature      string `json:"business_nature" valid:"required"`
	BusinessType        string `json:"business_type" valid:"required"`
	RegistrationNumber  string `json:"registration_number" valid:"required"`
	RegistrationExpiry  string `json:"registration_expiry" valid:"required"`
	EntityType          string `json:"entity_type" valid:"required"`
	ExchangeListed      string `json:"exchange_listed"`
}

func (rr *CorporationDetails) Populate(entity *KYCEntity) {
	if entity == nil {
		return
	}
	err := json.Unmarshal(entity.Data, &rr)
	if err != nil {
		panic(err)
	}
}

func (d CorporationDetails) KYCEntity() (KYCEntity, error) {
	data, err := json.Marshal(&d)
	return KYCEntity{
		Data: data,
		Type: KYCEntityTypeCorporationDetails,
	}, errors.Wrap(err, "failed to marshal")
}
