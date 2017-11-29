package api

import (
	"errors"

	"github.com/asaskevich/govalidator"
)

type BusinessDetails struct {
	CorporationDetails    CorporationDetails       `json:"corporation_details" valid:"required"`
	FinancialDetails      FinancialDetails         `json:"financial_details" valid:"required"`
	CorrespondenceAddress Address                  `json:"correspondence_address" valid:"required"`
	RegisteredAddress     Address                  `json:"registered_address" valid:"required"`
	Owners                map[int64]BusinessPerson `json:"owners"`
	Signatories           map[int64]BusinessPerson `json:"signatories"`
}

func (d *BusinessDetails) DisplayName() *string {
	if name := d.CorporationDetails.EntityName; name != "" {
		return &name
	}
	return nil
}

func (d *BusinessDetails) Validate() error {
	if len(d.Owners) == 0 {
		return errors.New("at least single owner required")
	}
	if len(d.Signatories) == 0 {
		return errors.New("at least single signatory required")
	}
	ok, err := govalidator.ValidateStruct(d)
	if !ok {
		return err
	}
	for _, person := range d.Owners {
		ok, err := govalidator.ValidateStruct(person)
		if !ok {
			return err
		}
	}
	for _, person := range d.Signatories {
		ok, err := govalidator.ValidateStruct(person)
		if !ok {
			return err
		}
	}
	return nil
}

func (d *BusinessDetails) RequiredDocuments() []RequiredDocument {
	docs := []RequiredDocument{
		{DocumentTypeRegistrationCertificate, 0},
		{DocumentTypeIncorporationCertificate, 0},
		{DocumentTypeMOA, 0},
		{DocumentTypeAOA, 0},
		{DocumentTypePhoneProof, 0},
		{DocumentTypeSignedForm, 0},
	}
	for entityID, _ := range d.Signatories {
		docs = append(docs, RequiredDocument{DocumentTypeSelfie, entityID})
		docs = append(docs, RequiredDocument{DocumentTypeIDProof, entityID})
		//docs = append(docs, RequiredDocument{DocumentTy})
	}
	for entityID, _ := range d.Owners {
		docs = append(docs, RequiredDocument{DocumentTypeSelfie, entityID})
		docs = append(docs, RequiredDocument{DocumentTypeIDProof, entityID})
	}
	return docs
}
