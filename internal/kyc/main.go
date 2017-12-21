package kyc

import (
	"encoding/json"
	"fmt"

	"database/sql/driver"

	"github.com/go-ozzo/ozzo-validation"
	"github.com/pkg/errors"
	"gitlab.com/swarmfund/api/internal/types"
)

type Individual struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

type Entity struct {
	Type       types.KYCEntityType `json:"type"`
	Individual *Individual         `json:"individual,omitempty"`
}

func (e Entity) Validate() error {
	return validation.ValidateStruct(&e,
		validation.Field(&e.Type, validation.Required),
	)
}

func (e *Entity) UnmarshalJSON(data []byte) error {
	var t struct {
		Type types.KYCEntityType `json:"type"`
		Raw  json.RawMessage     `json:"attributes"`
	}
	if err := json.Unmarshal(data, &t); err != nil {
		return errors.Wrap(err, "failed to unmarshal")
	}

	e.Type = t.Type

	var specific interface{}
	switch t.Type {
	case types.KYCEntityTypeIndividual:
		specific = &e.Individual
	default:
		panic(fmt.Errorf("unknown user type %s", t.Type))
	}

	if err := json.Unmarshal(t.Raw, &specific); err != nil {
		return errors.Wrap(err, "failed to unmarshal details")
	}

	return nil
}

func (e *Entity) Scan(src interface{}) error {
	bytes, ok := src.([]byte)
	if !ok {
		return fmt.Errorf("can't scan from %T", src)
	}
	return json.Unmarshal(bytes, &e)
}

func (e Entity) Value() (driver.Value, error) {
	return json.Marshal(e)
}
