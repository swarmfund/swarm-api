package types

import (
	"encoding/json"
	"fmt"

	"database/sql/driver"

	"github.com/pkg/errors"
)

type UserType int32

const (
	UserTypeNotVerified UserType = 1 << iota
	UserTypeSyndicate
)

var (
	userTypeReverseMap = map[string]UserType{
		"not_verified": UserTypeNotVerified,
		"syndicate":    UserTypeSyndicate,
	}

	userTypeMap = map[UserType]string{
		UserTypeNotVerified: "not_verified",
		UserTypeSyndicate:   "syndicate",
	}
	ErrUserTypeInvalid = errors.New("user type is invalid")
)

func (t UserType) Validate() error {
	if t < UserTypeNotVerified || t > UserTypeSyndicate {
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
		*t = userTypeReverseMap[v]
	default:
		return fmt.Errorf("invalid value for %T: %v", t, v)
	}
	return nil
}

func (t UserType) MarshalJSON() ([]byte, error) {
	return json.Marshal(userTypeMap[t])
}

func (t UserType) Value() (driver.Value, error) {
	return int64(t), nil
}

func (t *UserType) Scan(src interface{}) error {
	i, ok := src.(int64)
	if !ok {
		return fmt.Errorf("can't scan from %T", src)
	}
	*t = UserType(i)
	return nil
}
