package api

import (
	"math"

	"gitlab.com/swarmfund/api/render/hal"
)

type KdfParams struct {
	Algorithm string  `json:"algorithm"`
	Bits      uint    `json:"bits"`
	N         float64 `json:"n"`
	R         uint    `json:"r"`
	P         uint    `json:"p"`
}

func (p *KdfParams) Populate() {
	p.Algorithm = "scrypt"
	p.Bits = 256
	p.N = math.Pow(2, 12)
	p.R = 8
	p.P = 1
}

type KdfParamsAction struct {
	Action
}

func (action *KdfParamsAction) JSON() {
	action.ValidateBodyType()
	action.Do(func() {
		var response KdfParams
		response.Populate()
		hal.Render(action.W, response)
	})
}
