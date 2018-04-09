package types

//go:generate jsonenums -type=UserType -tprefix=false -transform=snake
type UserType int32

var (
	DefaultUserType = UserTypeNotVerified
)

const (
	UserTypeNotVerified UserType = 1 << iota
	UserTypeSyndicate
	UserTypeGeneral
)
