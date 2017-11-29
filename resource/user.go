package resource

import (
	"encoding/json"

	"fmt"

	"gitlab.com/swarmfund/api/db2/api"
	"gitlab.com/swarmfund/api/internal/types"
)

type User struct {
	PT              string          `json:"paging_token"`
	Address         types.Address   `json:"address"`
	Email           string          `json:"email"`
	UserType        api.UserType    `json:"user_type"`
	State           api.UserState   `json:"state"`
	RejectReasons   interface{}     `json:"reject_reasons,omitempty"`
	SignupLink      string          `json:"signup_url,omitempty"`
	IntegrationMeta json.RawMessage `json:"integration_meta,omitempty"`
	Details         interface{}     `json:"details,omitempty"`
}

func (u *User) Populate(hu *api.User) {
	u.PT = fmt.Sprintf("%d", hu.ID)
	u.Email = hu.Email
	u.Address = hu.Address
	u.UserType = hu.UserType
	u.State = hu.State
	u.Details = hu.Details()
	u.RejectReasons = hu.RejectReasons()
	u.IntegrationMeta = hu.IntegrationMeta
}

func (r User) PagingToken() string {
	return r.PT
}
