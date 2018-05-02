package api

import (
	"time"

	"gitlab.com/swarmfund/api/internal/types"
)

type UserStateUpdate struct {
	Address   types.Address
	Timestamp time.Time
	Type      *types.UserType
	State     *types.UserState
	KYCBlob   *string
}

func (u *UserStateUpdate) GetLoganFields() map[string]interface{} {
	fields := map[string]interface{}{
		"state_address":   u.Address,
		"state_timestamp": u.Timestamp,
	}

	if u.Type != nil {
		fields["state_type"] = u.Type
	}

	if u.State != nil {
		fields["state_state"] = u.State
	}

	if u.KYCBlob != nil {
		fields["state_kyc"] = u.KYCBlob
	}

	return fields
}

type UserStateQ interface {
	SetState(update UserStateUpdate) error
}
