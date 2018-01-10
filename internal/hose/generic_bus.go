package hose

import (
	"sync"
)

//go:generate genny -in=$GOFILE -out=changes_bus_generated.go gen "T=Changes"

type TCallback func(TEvent)

type TBus struct {
	*sync.RWMutex
	callbacks []TCallback
}

func NewTBus() *TBus {
	return &TBus{
		&sync.RWMutex{},
		make([]TCallback, 0),
	}
}

func (b *TBus) Subscribe(cb TCallback) *TBus {
	b.Lock()
	defer b.Unlock()
	b.callbacks = append(b.callbacks, cb)
	return b
}

func (b *TBus) Dispatch(event TEvent) {
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
