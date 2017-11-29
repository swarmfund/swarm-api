package api

import (
	"database/sql/driver"
	"errors"
)

type UserState string
type UserRecoveryState int
type UserLimitReviewState int

const (
	UserNeedDocs           UserState = "need_documents"
	UserWaitingForApproval UserState = "waiting_for_approval"
	UserApproved           UserState = "approved"
	UserRejected           UserState = "rejected"
	UserWalletPending      UserState = "wallet_pending"
)

const (
	UserRecoveryStateNil UserRecoveryState = iota
	UserRecoveryStatePending
)

const (
	UserLimitReviewNil UserLimitReviewState = iota
	UserLimitReviewPending
)

var (
	userStates = map[UserState]bool{
		UserNeedDocs:           true,
		UserWaitingForApproval: true,
		UserApproved:           true,
		UserRejected:           true,
		UserWalletPending:      true,
	}
)

func IsUserState(state string) bool {
	_, ok := userStates[UserState(state)]
	return ok
}

func (us UserState) Value() (driver.Value, error) {
	return string(us), nil
}

func (us *UserState) Scan(value interface{}) error {
	// if value is nil, empty string
	if value == nil {
		// set the value of the pointer us to UserType(false)
		*us = UserState("")
		return nil
	}
	if sv, err := driver.String.ConvertValue(value); err == nil {
		// if this is a string type
		switch v := sv.(type) {
		case string:
			*us = UserState(v)
			return nil
		case []byte:
			*us = UserState(string(v))
			return nil
		}
	}
	// otherwise, return an error
	return errors.New("failed to scan UserState")
}
