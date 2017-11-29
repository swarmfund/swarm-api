package api

import (
	"encoding/json"

	"github.com/pkg/errors"
)

type RequiredDocument struct {
	Type     DocumentType
	EntityID int64
}

func (user *User) CheckState() UserState {
	if err := user.Details().Validate(); err != nil {
		return UserRejected
	}

	if !user.RejectReasons().Empty() {
		return UserRejected
	}

	requiredDocs := user.Details().RequiredDocuments()
	for _, check := range requiredDocs {
		if ok := user.HaveDocument(check.EntityID, check.Type); !ok {
			return UserNeedDocs
		}
	}

	return UserWaitingForApproval
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
