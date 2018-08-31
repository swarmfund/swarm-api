package hose

import (
	"sync"

	"gitlab.com/distributed_lab/logan/v3"
)

//go:generate genny -in=$GOFILE -out=tx_bus_generated.go gen "Generic=Transaction"
//go:generate genny -in=$GOFILE -out=audit_log_bus_generated.go gen "Generic=Log"

type GenericCallback func(GenericEvent)
type GenericDispatch func(GenericEvent)

type GenericBus struct {
	*sync.RWMutex
	log       *logan.Entry
	callbacks []GenericCallback
}

func NewGenericBus(log *logan.Entry) *GenericBus {
	return &GenericBus{
		&sync.RWMutex{},
		log,
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
				if v := recover(); v != nil { // a panic is detected.
					b.log.WithRecover(v).Error("GenericBus crashed", logan.F{
						"event": event,
					})
				}
			}()
			cb(event)
		}()
	}
}
