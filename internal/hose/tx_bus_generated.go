// This file was automatically generated by genny.
// Any changes will be lost if this file is regenerated.
// see https://github.com/cheekybits/genny

package hose

import (
	"sync"

	"gitlab.com/distributed_lab/logan/v3"
)

type TransactionCallback func(TransactionEvent)
type TransactionDispatch func(TransactionEvent)

type TransactionBus struct {
	*sync.RWMutex
	log       *logan.Entry
	callbacks []TransactionCallback
}

func NewTransactionBus(log *logan.Entry) *TransactionBus {
	return &TransactionBus{
		&sync.RWMutex{},
		log,
		make([]TransactionCallback, 0),
	}
}

func (b *TransactionBus) Subscribe(cb TransactionCallback) *TransactionBus {
	b.Lock()
	defer b.Unlock()
	b.callbacks = append(b.callbacks, cb)
	return b
}

func (b *TransactionBus) Dispatch(event TransactionEvent) {
	b.Lock()
	defer b.Unlock()

	for _, cb := range b.callbacks {
		func() {
			defer func() {
				if v := recover(); v != nil { // a panic is detected.
					b.log.WithRecover(v).Error("TransactionBus crashed", logan.F{
						"event": event,
					})
				}
			}()
			cb(event)
		}()
	}
}
