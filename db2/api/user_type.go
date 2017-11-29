package api

import "errors"

// UserType represents the type of user.
type UserType string

const (
	// UserTypeNil is internal only, should not go anywhere
	UserTypeNil        UserType = ""
	UserTypeIndividual UserType = "individual"
	UserTypeJoint      UserType = "joint"
	UserTypeBusiness   UserType = "business"
)

var (
	userTypes = map[UserType]bool{
		UserTypeIndividual: true,
		UserTypeJoint:      true,
		UserTypeBusiness:   true,
	}
	ErrUnknownUserType = errors.New("unknown user type")
)

func (t UserType) Validate() error {
	if _, ok := userTypes[t]; !ok {
		return ErrUnknownUserType
	}
	return nil
}
