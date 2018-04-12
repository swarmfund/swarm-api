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

type UserStateQ interface {
	SetState(update UserStateUpdate) error
}
