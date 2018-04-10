package types

//go:generate jsonenums -type=UserState -tprefix=false -transform=snake
type UserState int32

var (
	DefaultUserState = UserStateNil
)

const (
	UserStateUndefined          UserState = 0
	UserStateNil                UserState = 1
	UserStateWaitingForApproval UserState = 2
	UserStateApproved           UserState = 4
	UserStateRejected           UserState = 8
)
