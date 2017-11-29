package api

import (
	"encoding/json"
	"time"
)

const (
	_ KYCEntityType = iota
	KYCEntityTypeBankDetails
	KYCEntityTypeFinancialDetails
	_
	KYCEntityTypeRegisteredAddress
	KYCEntityTypeCorrespondenceAddress // 5
	KYCEntityTypeCorporationDetails
	KYCEntityTypeEmploymentDetails
	KYCEntityTypePersonalDetails
	KYCEntityTypeIndividualRejectReasons
	KYCEntityTypeJointIdentity // 10
	KYCEntityTypeJointRejectReasons
	KYCEntityTypeBusinessRejectReasons
	KYCEntityTypeBusinessOwner
	KYCEntityTypeBusinessSignatory
	KYCEntityTypeAddress // 15
	KYCEntityTypeDocumentsRejectReasons
)

type KYCEntityType int

type KYCEntityGetter interface {
	KYCEntity() (KYCEntity, error)
}

type KYCEntity struct {
	ID        int64           `db:"id"`
	UserID    int64           `db:"user_id"`
	Data      json.RawMessage `db:"data"`
	Type      KYCEntityType   `db:"type"`
	CreatedAt time.Time       `db:"created_at"`
	UpdatedAt time.Time       `db:"updated_at"`
}

func (entity *KYCEntity) JointIdentity() (result IdentityDetails, err error) {
	err = json.Unmarshal(entity.Data, &result)
	return
}

func (entity *KYCEntity) BusinessPerson() (result BusinessPerson, err error) {
	err = json.Unmarshal(entity.Data, &result)
	return
}
