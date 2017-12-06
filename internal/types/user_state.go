package types

type UserState int32

const (
	UserStateNil UserState = 1 << iota
	UserStateWaitingForApproval
	UserStateApproved
	UserStateRejected
)

var (
	userStateMap = map[UserType]string{}
)
