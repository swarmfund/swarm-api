package resource

import (
	"gitlab.com/swarmfund/api/db2/api"
)

type Notifications struct {
	Email string `json:"email"`
	KYC   bool   `json:"kyc_review"`
}

func (n *Notifications) Populate(records *api.Notifications) {
	for _, record := range records.Records {
		n.Email = record.Email
		switch record.Type {
		case api.NotificationTypeKYC:
			n.KYC = true
		}
	}
}
