package hose

import "gitlab.com/swarmfund/api/internal/types"

type UserEventType int32

const (
	UserEventTypeCreated UserEventType = iota + 1
)

type User struct {
	Email   string
	Address types.Address
}

type UserEvent struct {
	Type UserEventType
	User User
}
