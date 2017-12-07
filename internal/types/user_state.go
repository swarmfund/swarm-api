package types

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"

	"github.com/pkg/errors"
)

type UserState int32

const (
	UserStateNil UserState = 1 << iota
	UserStateWaitingForApproval
	UserStateApproved
	UserStateRejected
)

var (
	userStateMap = map[UserState]string{
		UserStateNil:                "nil",
		UserStateWaitingForApproval: "waiting_for_approval",
		UserStateApproved:           "approved",
		UserStateRejected:           "rejected",
	}
	userStateReverseMap = map[string]UserState{
		"waiting_for_approval": UserStateWaitingForApproval,
		"approved":             UserStateApproved,
		"nil":                  UserStateNil,
		"rejected":             UserStateRejected,
	}
	ErrUserStateInvalid = errors.New("user state invalid")
)

func (t UserState) Validate() error {
	_, ok := userStateMap[t]
	if !ok {
		return ErrUserStateInvalid
	}
	return nil
}

func (t *UserState) UnmarshalJSON(data []byte) error {
	var tmp interface{}

	if err := json.Unmarshal(data, &tmp); err != nil {
		return errors.Wrap(err, "failed to unmarshal")
	}

	switch v := tmp.(type) {
	case string:
		*t = userStateReverseMap[v]
	default:
		return fmt.Errorf("invalid value for %T: %v", t, v)
	}
	return nil
}

func (t UserState) MarshalJSON() ([]byte, error) {
	return json.Marshal(userStateMap[t])
}

func (t UserState) Value() (driver.Value, error) {
	return int64(t), nil
}

func (t *UserState) Scan(src interface{}) error {
	i, ok := src.(int64)
	if !ok {
		return fmt.Errorf("can't scan from %T", src)
	}
	*t = UserState(i)
	return nil
}
