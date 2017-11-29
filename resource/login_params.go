package resource

import (
	"gitlab.com/swarmfund/api/db2/api"
)

type LoginParams struct {
	Username string `json:"username"`
	Salt     string `json:"salt"`
	Kdf      string `json:"kdfParams"`
}

func (p *LoginParams) Populate(w *api.Wallet) {
	p.Username = w.Username
	p.Salt = w.Salt
	//p.Kdf = w.KDF
}
