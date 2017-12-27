// generated by jsonenums -tprefix=false -transform=snake -type=KYCEntityType; DO NOT EDIT
package types

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
)

func init() {
	// stubs for imports
	_ = json.Delim('s')
	_ = driver.Int32

}

var ErrKYCEntityTypeInvalid = errors.New("KYCEntityType is invalid")

func init() {
	var v KYCEntityType
	if _, ok := interface{}(v).(fmt.Stringer); ok {
		_KYCEntityTypeNameToValue = map[string]KYCEntityType{
			interface{}(KYCEntityTypeIndividual).(fmt.Stringer).String(): KYCEntityTypeIndividual,
			interface{}(KYCEntityTypeSyndicate).(fmt.Stringer).String():  KYCEntityTypeSyndicate,
		}
	}
}

var _KYCEntityTypeNameToValue = map[string]KYCEntityType{
	"individual": KYCEntityTypeIndividual,
	"syndicate":  KYCEntityTypeSyndicate,
}

var _KYCEntityTypeValueToName = map[KYCEntityType]string{
	KYCEntityTypeIndividual: "individual",
	KYCEntityTypeSyndicate:  "syndicate",
}

// String is generated so KYCEntityType satisfies fmt.Stringer.
func (r KYCEntityType) String() string {
	s, ok := _KYCEntityTypeValueToName[r]
	if !ok {
		return fmt.Sprintf("KYCEntityType(%d)", r)
	}
	return s
}

// Validate verifies that value is predefined for KYCEntityType.
func (r KYCEntityType) Validate() error {
	_, ok := _KYCEntityTypeValueToName[r]
	if !ok {
		return ErrKYCEntityTypeInvalid
	}
	return nil
}

// MarshalJSON is generated so KYCEntityType satisfies json.Marshaler.
func (r KYCEntityType) MarshalJSON() ([]byte, error) {
	if s, ok := interface{}(r).(fmt.Stringer); ok {
		return json.Marshal(s.String())
	}
	s, ok := _KYCEntityTypeValueToName[r]
	if !ok {
		return nil, fmt.Errorf("invalid KYCEntityType: %d", r)
	}
	return json.Marshal(s)
}

// UnmarshalJSON is generated so KYCEntityType satisfies json.Unmarshaler.
func (r *KYCEntityType) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return fmt.Errorf("KYCEntityType should be a string, got %s", data)
	}
	v, ok := _KYCEntityTypeNameToValue[s]
	if !ok {
		return fmt.Errorf("invalid KYCEntityType %q", s)
	}
	*r = v
	return nil
}

func (t *KYCEntityType) Scan(src interface{}) error {
	i, ok := src.(int64)
	if !ok {
		return fmt.Errorf("can't scan from %T", src)
	}
	*t = KYCEntityType(i)
	return nil
}

func (t KYCEntityType) Value() (driver.Value, error) {
	return int64(t), nil
}