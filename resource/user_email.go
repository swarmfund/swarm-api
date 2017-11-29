package resource

import (
	"gitlab.com/swarmfund/api/db2/api"
)

type ShortenUserDetails struct {
	UserType  string  `json:"user_type"`
	UserState string  `json:"user_state"`
	Email     string  `json:"email"`
	FullName  *string `json:"full_name,omitempty"`
}

type ShortenUsersDetails struct {
	Users map[string]ShortenUserDetails `json:"users"`
}

func (d *ShortenUsersDetails) Populate(records []api.User) {
	d.Users = map[string]ShortenUserDetails{}
	for _, record := range records {
		d.Users[string(record.Address)] = ShortenUserDetails{
			UserType:  string(record.UserType),
			UserState: string(record.State),
			Email:     record.Email,
			FullName:  record.Details().DisplayName(),
		}
	}
}
