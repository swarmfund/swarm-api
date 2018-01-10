package kyc

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/pkg/errors"
	"gitlab.com/swarmfund/api/internal/types"
)

type Entity struct {
	Type       types.KYCEntityType `json:"type"`
	Individual *Individual         `json:"individual,omitempty"`
	Syndicate  *Syndicate          `json:"syndicate,omitempty"`
}

func (e Entity) Attributes() interface{} {
	switch e.Type {
	case types.KYCEntityTypeIndividual:
		return e.Individual
	case types.KYCEntityTypeSyndicate:
		return e.Syndicate
	default:
		panic(fmt.Errorf("unknown user type %s", e.Type))
	}
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
	case types.KYCEntityTypeSyndicate:
		specific = &e.Syndicate
	default:
		panic(fmt.Errorf("unknown user type %s", t.Type))
	}

	if err := json.Unmarshal(t.Raw, &specific); err != nil {
		return errors.Wrap(err, "failed to unmarshal details")
	}

	return nil
}

func (e Entity) MarshalJSON() ([]byte, error) {
	var t struct {
		Type types.KYCEntityType `json:"type"`
		Raw  interface{}         `json:"attributes"`
	}

	t.Type = e.Type
	t.Raw = e.Attributes()

	return json.Marshal(t)
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
