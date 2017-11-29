package resource

import "gitlab.com/swarmfund/api/db2/api"

type Documents struct {
	Version   int64         `json:"version"`
	Documents api.Documents `json:"documents"`
}

func (d *Documents) Populate(user *api.User) {
	d.Version = user.DocumentsVersion
	d.Documents = user.Documents
}
