package api

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"

	"time"

	"regexp"

	"github.com/guregu/null"
	"github.com/pkg/errors"
	"gitlab.com/swarmfund/api/internal/types"
)

//User is a row of data from the `users` table
// TODO
// * remove DeletedAt
// * UpdatedAt probably useless because of joins
type User struct {
	ID      int64         `db:"id"`
	Address types.Address `db:"address"`
	// TODO join from wallets address->account_id
	Email             string               `db:"email"`
	UserType          types.UserType       `db:"type"`
	State             types.UserState      `db:"state"`
	KYCSequence       int64                `db:"kyc_sequence"`
	RejectReason      string               `db:"reject_reason"`
	Documents         Documents            `db:"documents"`
	DocumentsVersion  int64                `db:"documents_version"`
	LimitReviewStatue UserLimitReviewState `db:"limit_review_state"`
	CreatedAt         time.Time            `db:"created_at"`
	UpdatedAt         string               `db:"updated_at"`
	DeletedAt         sql.NullString       `db:"deleted_at"`

	// Nickname comes from join on contacts table when needed
	Nickname null.String `db:"nickname"`
	// IntegrationMeta comes from join on exchange_integrations
	IntegrationMeta json.RawMessage `db:"integration_meta"`
	// KYCEntities comes from join on kyc_entities
	KYCEntities KYCEntities `db:"kyc_entities"`
	// RecoveryAddress recovery signer, comes from join on recoveries
	RecoveryAddress types.Address `db:"recovery_address"`
	// AirdropState treat nils as valid undefined value, comes from json on airdrops
	AirdropState *types.AirdropState `db:"airdrop_state"`
}

func (user *User) IsAirdropEligible() bool {
	pattern := `(?i).+@(163.com|qq.com|126.com|189.cn|139.com|sina.com|aliyun.com|xinjiyuan99.com)`
	blacklisted, err := regexp.MatchString(pattern, user.Email)
	if err != nil {
		panic(errors.Wrap(err, "blacklist check failed"))
	}
	return !blacklisted
}

// Details will throw panic aggressively instead of returning error
// to allow chain calls.
func (user *User) Details() UserDetails {
	switch user.UserType {
	//case types.UserTypeNotVerified:

	//case UserTypeIndividual:
	//	details := IndividualDetails{}
	//
	//	bankDetails := entities.GetSingle(KYCEntityTypeBankDetails)
	//	details.BankDetails.Populate(bankDetails)
	//
	//	personalDetails := entities.GetSingle(KYCEntityTypePersonalDetails)
	//	details.PersonalDetails.Populate(personalDetails)
	//
	//	employment := entities.GetSingle(KYCEntityTypeEmploymentDetails)
	//	details.EmploymentDetails.Populate(employment)
	//
	//	address := entities.GetSingle(KYCEntityTypeAddress)
	//	details.Address.Populate(address)
	//
	//	return &details
	//case UserTypeJoint:
	//	details := JointDetails{}
	//
	//	identityEntities := entities.Get(KYCEntityTypeJointIdentity)
	//	details.Identities = map[int64]IdentityDetails{}
	//	for _, entity := range identityEntities {
	//		entity := entity
	//		identity := IdentityDetails{}
	//		identity.Populate(&entity)
	//		details.Identities[entity.ID] = identity
	//	}
	//
	//	return &details
	//case UserTypeBusiness:
	//	details := BusinessDetails{}
	//
	//	ownerEntities := entities.Get(KYCEntityTypeBusinessOwner)
	//	details.Owners = map[int64]BusinessPerson{}
	//	for _, entity := range ownerEntities {
	//		entity := entity
	//		person := BusinessPerson{}
	//		person.Populate(&entity)
	//		details.Owners[entity.ID] = person
	//	}
	//
	//	signatoryEntites := entities.Get(KYCEntityTypeBusinessSignatory)
	//	details.Signatories = map[int64]BusinessPerson{}
	//	for _, entity := range signatoryEntites {
	//		entity := entity
	//		person := BusinessPerson{}
	//		person.Populate(&entity)
	//		details.Signatories[entity.ID] = person
	//	}
	//
	//	entity := entities.GetSingle(KYCEntityTypeCorporationDetails)
	//	details.CorporationDetails.Populate(entity)
	//
	//	entity = entities.GetSingle(KYCEntityTypeRegisteredAddress)
	//	details.RegisteredAddress.Populate(entity)
	//
	//	entity = entities.GetSingle(KYCEntityTypeFinancialDetails)
	//	details.FinancialDetails.Populate(entity)
	//
	//	entity = entities.GetSingle(KYCEntityTypeCorrespondenceAddress)
	//	details.CorrespondenceAddress.Populate(entity)
	//
	//	return &details
	default:
		panic("unknown details user type")
	}
}

func (user *User) RejectReasons() RejectReasons {
	switch user.UserType {
	//case UserTypeIndividual:
	//	docrr := DocumentsRejectReasons{}
	//	entity := user.KYCEntities.GetSingle(KYCEntityTypeDocumentsRejectReasons)
	//	docrr.Populate(entity)
	//
	//	reasons := IndividualRejectReasons{}
	//	entity = user.KYCEntities.GetSingle(KYCEntityTypeIndividualRejectReasons)
	//	reasons.Populate(entity)
	//
	//	reasons.Documents = docrr
	//	return &reasons
	//case UserTypeJoint:
	//	docrr := DocumentsRejectReasons{}
	//	entity := user.KYCEntities.GetSingle(KYCEntityTypeDocumentsRejectReasons)
	//	docrr.Populate(entity)
	//
	//	reasons := JointRejectReasons{}
	//	entity = user.KYCEntities.GetSingle(KYCEntityTypeJointRejectReasons)
	//	reasons.Populate(entity)
	//	if reasons.IdentityDetails == nil {
	//		reasons.IdentityDetails = map[string]IdentityDetails{}
	//	}
	//
	//	reasons.Documents = docrr
	//	return &reasons
	//case UserTypeBusiness:
	//	docrr := DocumentsRejectReasons{}
	//	entity := user.KYCEntities.GetSingle(KYCEntityTypeDocumentsRejectReasons)
	//	docrr.Populate(entity)
	//
	//	reasons := BusinessRejectReasons{}
	//	entity = user.KYCEntities.GetSingle(KYCEntityTypeBusinessRejectReasons)
	//	reasons.Populate(entity)
	//	if reasons.Owners == nil {
	//		reasons.Owners = map[string]BusinessPerson{}
	//	}
	//	if reasons.Signatories == nil {
	//		reasons.Signatories = map[string]BusinessPerson{}
	//	}
	//
	//	reasons.Documents = docrr
	//
	//	return &reasons
	default:
		panic("unknown reject reasons user type")
	}
}

// nil document stands for deleted doc
type Documents map[DocumentType]map[string]*Document

func (d Documents) Value() (driver.Value, error) {
	j, err := json.Marshal(d)
	return j, err
}

func (d *Documents) Scan(src interface{}) error {
	source, ok := src.([]byte)
	if !ok {
		return errors.New("type assertion failed ([]byte)")
	}

	return json.Unmarshal(source, d)
}

func (d Documents) Get(fn func(doc *Document) bool) *Document {
	for docType, _ := range d {
		for _, document := range d[docType] {
			if fn(document) {
				return document
			}
		}
	}
	return nil
}

func (d Documents) Latest(docType DocumentType) (latest *Document) {
	for _, document := range d[0] {
		if document == nil {
			continue
		}
		if latest == nil || document.CreatedAt > latest.CreatedAt {
			latest = document
		}
	}
	return latest
}
