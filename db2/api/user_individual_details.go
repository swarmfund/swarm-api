package api

import (
	"github.com/asaskevich/govalidator"
)

type IndividualDetails struct {
	PersonalDetails   PersonDetails     `json:"personal_details" valid:"required"`
	Address           Address           `json:"address" valid:"required"`
	EmploymentDetails EmploymentDetails `json:"employment_details" valid:"required"`
	BankDetails       BankDetails       `json:"bank_details" valid:"required"`
}

func (d *IndividualDetails) Validate() error {
	ok, err := govalidator.ValidateStruct(d)
	if !ok {
		return err
	}
	return nil
}

func (d *IndividualDetails) RequiredDocuments() []RequiredDocument {
	return []RequiredDocument{
		{DocumentTypeBankProof, 0},
		{DocumentTypeSelfie, 0},
		{DocumentTypeIDProof, 0},
		{DocumentTypeSignedForm, 0},
	}
}

func (d *IndividualDetails) DisplayName() *string {
	return d.PersonalDetails.DisplayName()
}
