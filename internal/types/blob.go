package types

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"
)

type Blob struct {
	ID            string
	Type          BlobType
	Value         string
	Relationships BlobRelationships
	OwnerAddress  *Address   `db:"owner_address"`
	DeletedAt     *time.Time `db:"deleted_at"`
}

type BlobRelationships map[string]string

func (r *BlobRelationships) Scan(value interface{}) error {
	if value == nil {
		return nil
	}
	switch v := value.(type) {
	case []byte:
		return json.Unmarshal(v, &r)
	default:
		return fmt.Errorf("unsupported Scan from type %T", v)
	}
}

func (r BlobRelationships) Value() (driver.Value, error) {
	return json.Marshal(r)
}
