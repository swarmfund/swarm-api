package resources

import (
	"gitlab.com/swarmfund/api/db2/api"
	"gitlab.com/swarmfund/api/internal/types"
)

type (
	User struct {
		Type       string         `json:"type"`
		ID         types.Address  `json:"id"`
		Attributes UserAttributes `json:"attributes"`
	}
	UserAttributes struct {
		Email string `json:"email"`
	}
)

func NewUser(user *api.User) User {
	return User{
		Type: string(user.UserType),
		ID:   user.Address,
		Attributes: UserAttributes{
			Email: user.Email,
		},
	}
}
