package types

//go:generate jsonenums -type=UserState -tprefix=false -transform=snake
type UserState int32

var (
	DefaultUserState = UserStateNil
)

const (
	UserStateUndefined UserState = 1<<iota - 1
	UserStateNil
	UserStateWaitingForApproval
	UserStateApproved
	UserStateRejected
)
