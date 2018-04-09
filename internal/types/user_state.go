package types

//go:generate jsonenums -type=UserState -tprefix=false -transform=snake
type UserState int32

var (
	DefaultUserState = UserStateNil
)

const (
	UserStateNil UserState = 1 << iota
	UserStateWaitingForApproval
	UserStateApproved
	UserStateRejected
)
