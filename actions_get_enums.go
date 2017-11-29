package api

import (
	"gitlab.com/swarmfund/api/assets"
	"gitlab.com/swarmfund/api/render/hal"
)

type GetEnumsAction struct {
	Action
	Enums map[string][]string
}

// JSON format action handler
func (action *GetEnumsAction) JSON() {
	action.Do(
		action.ValidateBodyType,
		func() {
			hal.Render(action.W, assets.Enums.Data())
		})
}
