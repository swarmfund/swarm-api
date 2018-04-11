package types

//go:generate jsonenums -type=UserType -tprefix=false -transform=snake
type UserType int32

var (
	DefaultUserType = UserTypeNotVerified
)

const (
	UserTypeUndefined   UserType = 0
	UserTypeNotVerified UserType = 1
	UserTypeSyndicate   UserType = 2
	UserTypeGeneral     UserType = 4
)
