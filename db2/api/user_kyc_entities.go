package api

import (
	"encoding/json"

	"github.com/pkg/errors"

	"reflect"

	"fmt"
)

type PersonRejectReasons struct {
	PersonDetails
	Selfie  string `json:"selfie"`
	IDProof string `json:"id_proof"`
}

/*
documents: {
	<entity>: {
		<type>: "value"
	}
}
*/

type DocumentsRejectReasons map[int64]map[DocumentType]string

func (rr *DocumentsRejectReasons) Populate(entity *KYCEntity) {
	if entity == nil {
		*rr = DocumentsRejectReasons{}
		return
	}
	err := json.Unmarshal(entity.Data, &rr)
	if err != nil {
		panic(err)
	}
}

func (rr DocumentsRejectReasons) Drop(eid int64, docType DocumentType) {
	if rr == nil || rr[eid] == nil {
		return
	}
	delete(rr[eid], docType)
	if len(rr[eid]) == 0 {
		delete(rr, eid)
	}
}

func (rr *DocumentsRejectReasons) Set(eid int64, docType DocumentType, reason string) {
	if rr == nil {
		*rr = DocumentsRejectReasons{}
	}
	if (*rr)[eid] == nil {
		(*rr)[eid] = map[DocumentType]string{}
	}
	(*rr)[eid][docType] = reason
}

type IndividualRejectReasons struct {
	General           string                 `json:"general"`
	PersonalDetails   PersonDetails          `json:"personal_details"`
	Address           Address                `json:"address"`
	EmploymentDetails EmploymentDetails      `json:"employment_details"`
	BankDetails       BankDetails            `json:"bank_details"`
	Documents         DocumentsRejectReasons `json:"documents"`
}

func (rr *IndividualRejectReasons) Empty() bool {
	if len(rr.Documents) != 0 {
		return false
	}
	rr.Documents = nil
	rr.General = ""
	return reflect.DeepEqual(*rr, IndividualRejectReasons{})
}

func (rr *IndividualRejectReasons) Populate(entity *KYCEntity) {
	if entity == nil {
		return
	}
	err := json.Unmarshal(entity.Data, &rr)
	if err != nil {
		panic(err)
	}
}

type JointRejectReasons struct {
	General         string                     `json:"general"`
	IdentityDetails map[string]IdentityDetails `json:"identities"`
	Documents       DocumentsRejectReasons     `json:"documents"`
}

func (rr *JointRejectReasons) Empty() bool {
	if len(rr.Documents) != 0 {
		return false
	}
	rr.Documents = nil
	rr.General = ""
	for _, identity := range rr.IdentityDetails {
		if !reflect.DeepEqual(identity, IdentityDetails{}) {
			return false
		}
	}
	rr.IdentityDetails = nil
	return reflect.DeepEqual(*rr, JointRejectReasons{})
}

func (rr *JointRejectReasons) Populate(entity *KYCEntity) {
	if entity == nil {
		return
	}
	err := json.Unmarshal(entity.Data, &rr)
	if err != nil {
		panic(err)
	}
}

type BusinessRejectReasons struct {
	General               string                    `json:"general"`
	CorporationDetails    CorporationDetails        `json:"corporation_details" valid:"required"`
	FinancialDetails      FinancialDetails          `json:"financial_details" valid:"required"`
	CorrespondenceAddress Address                   `json:"correspondence_address" valid:"required"`
	RegisteredAddress     Address                   `json:"registered_address" valid:"required"`
	Owners                map[string]BusinessPerson `json:"owners"`
	Signatories           map[string]BusinessPerson `json:"signatories"`
	Documents             DocumentsRejectReasons    `json:"documents"`
}

func (rr *BusinessRejectReasons) Empty() bool {
	if len(rr.Documents) != 0 {
		return false
	}
	rr.Documents = nil
	rr.General = ""
	for _, entity := range rr.Signatories {
		if !reflect.DeepEqual(entity, BusinessPerson{}) {
			return false
		}
	}
	for _, entity := range rr.Owners {
		if !reflect.DeepEqual(entity, BusinessPerson{}) {
			return false
		}
	}
	rr.Signatories = nil
	rr.Owners = nil
	return reflect.DeepEqual(*rr, BusinessRejectReasons{})
}

func (rr *BusinessRejectReasons) Populate(entity *KYCEntity) {
	if entity == nil {
		return
	}
	err := json.Unmarshal(entity.Data, &rr)
	if err != nil {
		panic(err)
	}
}

type IdentityDetails struct {
	PersonalDetails   PersonDetails     `json:"personal_details" valid:"required"`
	Address           Address           `json:"address" valid:"required"`
	EmploymentDetails EmploymentDetails `json:"employment_details" valid:"required"`
	BankDetails       BankDetails       `json:"bank_details" valid:"required"`
}

func (d *IdentityDetails) Populate(entity *KYCEntity) {
	if entity == nil {
		return
	}
	err := json.Unmarshal(entity.Data, &d)
	if err != nil {
		panic(err)
	}
}

func (d IdentityDetails) KYCEntity() (KYCEntity, error) {
	data, err := json.Marshal(&d)
	return KYCEntity{
		Data: data,
		Type: KYCEntityTypeJointIdentity,
	}, errors.Wrap(err, "failed to marshal")
}

type BusinessPerson struct {
	PersonDetails PersonDetails `json:"personal_details"`
	Address       Address       `json:"address"`
}

func (d *BusinessPerson) Populate(entity *KYCEntity) {
	if entity == nil {
		return
	}
	err := json.Unmarshal(entity.Data, &d)
	if err != nil {
		panic(err)
	}
}

// KYCEntity satisfies `KYCEntityGetter` interface
// NOTE: don't forget to set type afterwards
func (d BusinessPerson) KYCEntity() (KYCEntity, error) {
	data, err := json.Marshal(&d)
	return KYCEntity{
		Data: data,
	}, errors.Wrap(err, "failed to marshal")
}

type PersonDetails struct {
	FirstName       string `json:"first_name" valid:"required"`
	LastName        string `json:"last_name" valid:"required"`
	FormerFirstName string `json:"former_first_name"`
	FormerLastName  string `json:"former_last_name"`
	DOB             string `json:"date_of_birth" valid:"required"`
	POB             string `json:"place_of_birth" valid:"required"`
	Nationality     string `json:"nationality" valid:"required"`
	Gender          string `json:"gender" valid:"required"`
	MaritalStatus   string `json:"marital_status" valid:"required"`
	Email           string `json:"email" valid:"required"`
	Mobile          string `json:"mobile" valid:"required"`
	IDDocument      string `json:"id_document" valid:"required"`
	IDNumber        string `json:"id_number" valid:"required"`
}

func (d *PersonDetails) DisplayName() *string {
	name := fmt.Sprintf("%s %s", d.FirstName, d.LastName)
	if len(name) > 1 {
		return &name
	}
	return nil
}

func (d *PersonDetails) Populate(entity *KYCEntity) {
	if entity == nil {
		return
	}
	err := json.Unmarshal(entity.Data, &d)
	if err != nil {
		panic(err)
	}
}

func (d PersonDetails) KYCEntity() (KYCEntity, error) {
	data, err := json.Marshal(&d)
	return KYCEntity{
		Data: data,
		Type: KYCEntityTypePersonalDetails,
	}, errors.Wrap(err, "failed to marshal")
}

type EmploymentDetails struct {
	Education string `json:"education" valid:"required"`
	Industry  string `json:"industry" valid:"required"`
	Status    string `json:"status" valid:"required"`
}

func (d *EmploymentDetails) Populate(entity *KYCEntity) {
	if entity == nil {
		return
	}
	err := json.Unmarshal(entity.Data, &d)
	if err != nil {
		panic(err)
	}
}

func (d EmploymentDetails) KYCEntity() (KYCEntity, error) {
	data, err := json.Marshal(&d)
	return KYCEntity{
		Data: data,
		Type: KYCEntityTypeEmploymentDetails,
	}, errors.Wrap(err, "failed to marshal")
}

type BankDetails struct {
	SourceOfWealth string  `json:"source_of_wealth" valid:"required"`
	IBAN           string  `json:"iban" valid:"required"`
	BankName       string  `json:"bank_name" valid:"required"`
	BankAddress    Address `json:"bank_address" valid:"required"`
	BIC            string  `json:"bic" valid:"required"`
	SourceOfIncome string  `json:"source_of_income" valid:"required"`
}

func (d *BankDetails) Populate(entity *KYCEntity) {
	if entity == nil {
		return
	}
	err := json.Unmarshal(entity.Data, &d)
	if err != nil {
		panic(err)
	}
}

func (d BankDetails) KYCEntity() (KYCEntity, error) {
	data, err := json.Marshal(&d)
	return KYCEntity{
		Data: data,
		Type: KYCEntityTypeBankDetails,
	}, errors.Wrap(err, "failed to marshal")
}

type FinancialDetails struct {
	GroupName       string `json:"group_name"`
	GroupAddress    string `json:"group_address"`
	SourceOfFunds   string `json:"source_of_funds" valid:"required"`
	SourceOfCapital string `json:"source_of_capital" valid:"required"`
	ExternalAudit   string `json:"external_audit"`
}

func (rr *FinancialDetails) Populate(entity *KYCEntity) {
	if entity == nil {
		return
	}
	err := json.Unmarshal(entity.Data, &rr)
	if err != nil {
		panic(err)
	}
}

func (d FinancialDetails) KYCEntity() (KYCEntity, error) {
	data, err := json.Marshal(&d)
	return KYCEntity{
		Data: data,
		Type: KYCEntityTypeFinancialDetails,
	}, errors.Wrap(err, "failed to marshal")
}

type Address struct {
	Line1      string `json:"line_1" valid:"required"`
	Line2      string `json:"line_2" valid:"required"`
	City       string `json:"city" valid:"required"`
	State      string `json:"state" valid:"required"`
	PostalCode string `json:"postal_code"`
	Country    string `json:"country" valid:"required"`
}

func (d *Address) Populate(entity *KYCEntity) {
	if entity == nil {
		return
	}
	err := json.Unmarshal(entity.Data, &d)
	if err != nil {
		panic(err)
	}
}

func (d Address) KYCEntity() (KYCEntity, error) {
	data, err := json.Marshal(&d)
	return KYCEntity{
		Data: data,
		Type: KYCEntityTypeAddress,
	}, errors.Wrap(err, "failed to marshal")
}
