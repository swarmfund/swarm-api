package hose

import (
	"context"
	"time"

	"gitlab.com/swarmfund/api/internal/types"
	"gitlab.com/tokend/go/resources"
)

type User struct {
	Email   string
	Address types.Address
	IP      string
}

type LogEvent struct {
	Type types.LogEventType
	User
	time.Time
	context.Context
	resources.Signer
}
