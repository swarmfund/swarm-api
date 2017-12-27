// generated by jsonenums -type=UserState -tprefix=false -transform=snake; DO NOT EDIT
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

var ErrUserStateInvalid = errors.New("UserState is invalid")

func init() {
	var v UserState
	if _, ok := interface{}(v).(fmt.Stringer); ok {
		_UserStateNameToValue = map[string]UserState{
			interface{}(UserStateNil).(fmt.Stringer).String():                UserStateNil,
			interface{}(UserStateWaitingForApproval).(fmt.Stringer).String(): UserStateWaitingForApproval,
			interface{}(UserStateApproved).(fmt.Stringer).String():           UserStateApproved,
			interface{}(UserStateRejected).(fmt.Stringer).String():           UserStateRejected,
		}
	}
}

var _UserStateNameToValue = map[string]UserState{
	"nil": UserStateNil,
	"waiting_for_approval": UserStateWaitingForApproval,
	"approved":             UserStateApproved,
	"rejected":             UserStateRejected,
}

var _UserStateValueToName = map[UserState]string{
	UserStateNil:                "nil",
	UserStateWaitingForApproval: "waiting_for_approval",
	UserStateApproved:           "approved",
	UserStateRejected:           "rejected",
}

// String is generated so UserState satisfies fmt.Stringer.
func (r UserState) String() string {
	s, ok := _UserStateValueToName[r]
	if !ok {
		return fmt.Sprintf("UserState(%d)", r)
	}
	return s
}

// Validate verifies that value is predefined for UserState.
func (r UserState) Validate() error {
	_, ok := _UserStateValueToName[r]
	if !ok {
		return ErrUserStateInvalid
	}
	return nil
}

// MarshalJSON is generated so UserState satisfies json.Marshaler.
func (r UserState) MarshalJSON() ([]byte, error) {
	if s, ok := interface{}(r).(fmt.Stringer); ok {
		return json.Marshal(s.String())
	}
	s, ok := _UserStateValueToName[r]
	if !ok {
		return nil, fmt.Errorf("invalid UserState: %d", r)
	}
	return json.Marshal(s)
}

// UnmarshalJSON is generated so UserState satisfies json.Unmarshaler.
func (r *UserState) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return fmt.Errorf("UserState should be a string, got %s", data)
	}
	v, ok := _UserStateNameToValue[s]
	if !ok {
		return fmt.Errorf("invalid UserState %q", s)
	}
	*r = v
	return nil
}

func (t *UserState) Scan(src interface{}) error {
	i, ok := src.(int64)
	if !ok {
		return fmt.Errorf("can't scan from %T", src)
	}
	*t = UserState(i)
	return nil
}

func (t UserState) Value() (driver.Value, error) {
	return int64(t), nil
}