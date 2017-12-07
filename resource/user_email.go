package resource

import (
	"gitlab.com/swarmfund/api/db2/api"
	"gitlab.com/swarmfund/api/internal/types"
)

type ShortenUserDetails struct {
	UserType  types.UserType  `json:"user_type"`
	UserState types.UserState `json:"user_state"`
	Email     string          `json:"email"`
}

type ShortenUsersDetails struct {
	Users map[types.Address]ShortenUserDetails `json:"users"`
}

func (d *ShortenUsersDetails) Populate(records []api.User) {
	d.Users = map[types.Address]ShortenUserDetails{}
	for _, record := range records {
		d.Users[record.Address] = ShortenUserDetails{
			UserType:  record.UserType,
			UserState: record.State,
			Email:     record.Email,
		}
	}
}
