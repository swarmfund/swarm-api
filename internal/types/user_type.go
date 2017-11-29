package types

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/pkg/errors"
)

const (
	_ UserType = iota
	UserTypeIndividual
)

var (
	ErrUserTypeInvalid = errors.New("user type is invalid")
)

type UserType int

func (t UserType) Validate() error {
	if t < UserTypeIndividual || t > UserTypeIndividual {
		return ErrUserTypeInvalid
	}
	return nil
}

// UnmarshalJSON custom unmarshaler supporting both JSON number and string
func (t *UserType) UnmarshalJSON(data []byte) error {
	var tmp interface{}

	if err := json.Unmarshal(data, &tmp); err != nil {
		return errors.Wrap(err, "failed to unmarshal")
	}

	switch v := tmp.(type) {
	case float64:
		*t = UserType(v)
	case string:
		i, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			return errors.Wrap(err, "failed to parse int")
		}
		*t = UserType(i)
	default:
		return fmt.Errorf("invalid value for %T: %v", t, v)
	}
	return nil
}
