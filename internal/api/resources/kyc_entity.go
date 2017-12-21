package resources

import (
	"gitlab.com/swarmfund/api/db2/api"
	"gitlab.com/swarmfund/api/internal/types"
)

type KYCEntity struct {
	ID         int64               `json:"id,string"`
	Type       types.KYCEntityType `json:"type"`
	Attributes interface{}         `json:"attributes"`
}

func NewKYCEntity(record api.KYCEntityRecord) KYCEntity {
	return KYCEntity{
		ID:         record.ID,
		Type:       record.Entity.Type,
		Attributes: record.Entity.Attributes(),
	}
}
