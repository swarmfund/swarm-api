package api

import (
	"encoding/json"
	"fmt"

	"github.com/pkg/errors"
)

type KYCEntities []KYCEntity

func (kyc *KYCEntities) Scan(value interface{}) error {
	if value == nil {
		*kyc = make(KYCEntities, 0)
		return nil
	}
	if bytes, ok := value.([]byte); ok {
		if err := json.Unmarshal(bytes, &kyc); err != nil {
			return errors.Wrap(err, "failed to scan KYCEntries")
		}
		return nil
	}
	return errors.New("failed to scan KYCEntries")
}

func (kyc KYCEntities) Exists(entityType KYCEntityType) (int64, bool) {
	for _, entity := range kyc {
		if entity.Type == entityType {
			return entity.ID, true
		}
	}
	return 0, false
}

func (kyc KYCEntities) Get(entityType KYCEntityType) KYCEntities {
	result := make(KYCEntities, 0)
	for _, entity := range kyc {
		if entity.Type == entityType {
			result = append(result, entity)
		}
	}
	return result
}

func (kyc KYCEntities) GetSingle(entityType KYCEntityType) *KYCEntity {
	entities := kyc.Get(entityType)
	if len(entities) > 1 {
		panic(fmt.Sprintf("too many %d", entityType))
	}
	if len(entities) == 1 {
		return &entities[0]
	}
	return nil
}

func (kyc KYCEntities) JointRejectReasons() (*KYCEntity, JointRejectReasons) {
	rejectReasons := JointRejectReasons{}
	entity := kyc.GetSingle(KYCEntityTypeJointRejectReasons)
	rejectReasons.Populate(entity)
	return entity, rejectReasons
}

func (kyc KYCEntities) BusinessRejectReasons() (entity *KYCEntity, rr BusinessRejectReasons) {
	entity = kyc.GetSingle(KYCEntityTypeBusinessRejectReasons)
	rr.Populate(entity)
	return entity, rr
}

func (kyc KYCEntities) DocumentsRejectReasons() (entity *KYCEntity, rr DocumentsRejectReasons) {
	entity = kyc.GetSingle(KYCEntityTypeDocumentsRejectReasons)
	rr.Populate(entity)
	return entity, rr
}
