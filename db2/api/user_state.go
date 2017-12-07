package api

type UserRecoveryState int
type UserLimitReviewState int

const (
	UserRecoveryStateNil UserRecoveryState = iota
	UserRecoveryStatePending
)

const (
	UserLimitReviewNil UserLimitReviewState = iota
	UserLimitReviewPending
)

//
//var (
//	userStates = map[UserState]bool{
//		UserNeedDocs:           true,
//		UserWaitingForApproval: true,
//		UserApproved:           true,
//		UserRejected:           true,
//		UserWalletPending:      true,
//	}
//)
//
//func IsUserState(state string) bool {
//	_, ok := userStates[UserState(state)]
//	return ok
//}
//
//func (us UserState) Value() (driver.Value, error) {
//	return string(us), nil
//}
//
//func (us *UserState) Scan(value interface{}) error {
//	// if value is nil, empty string
//	if value == nil {
//		// set the value of the pointer us to UserType(false)
//		*us = UserState("")
//		return nil
//	}
//	if sv, err := driver.String.ConvertValue(value); err == nil {
//		// if this is a string type
//		switch v := sv.(type) {
//		case string:
//			*us = UserState(v)
//			return nil
//		case []byte:
//			*us = UserState(string(v))
//			return nil
//		}
//	}
//	// otherwise, return an error
//	return errors.New("failed to scan UserState")
//}
