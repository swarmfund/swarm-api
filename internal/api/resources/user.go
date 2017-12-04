package resources

import "gitlab.com/swarmfund/api/internal/types"

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
