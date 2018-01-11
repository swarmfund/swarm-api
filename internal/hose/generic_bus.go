package hose

import (
	"sync"
)

//go:generate genny -in=$GOFILE -out=tx_bus_generated.go gen "Generic=Transaction"
//go:generate genny -in=$GOFILE -out=user_bus_generated.go gen "Generic=User"

type GenericCallback func(GenericEvent)
type GenericDispatch func(GenericEvent)

type GenericBus struct {
	*sync.RWMutex
	callbacks []GenericCallback
}

func NewGenericBus() *GenericBus {
	return &GenericBus{
		&sync.RWMutex{},
		make([]GenericCallback, 0),
	}
}

func (b *GenericBus) Subscribe(cb GenericCallback) *GenericBus {
	b.Lock()
	defer b.Unlock()
	b.callbacks = append(b.callbacks, cb)
	return b
}

func (b *GenericBus) Dispatch(event GenericEvent) {
	b.Lock()
	defer b.Unlock()

	for _, cb := range b.callbacks {
		func() {
			defer func() {
				// todo recover
			}()
			cb(event)
		}()
	}
}
