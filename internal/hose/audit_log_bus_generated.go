// This file was automatically generated by genny.
// Any changes will be lost if this file is regenerated.
// see https://github.com/cheekybits/genny

package hose

import (
	"sync"

	"gitlab.com/distributed_lab/logan/v3"
)

type LogCallback func(LogEvent)
type LogDispatch func(LogEvent)

type LogBus struct {
	*sync.RWMutex
	log       *logan.Entry
	callbacks []LogCallback
}

func NewLogBus(log *logan.Entry) *LogBus {
	return &LogBus{
		&sync.RWMutex{},
		log,
		make([]LogCallback, 0),
	}
}

func (b *LogBus) Subscribe(cb LogCallback) *LogBus {
	b.Lock()
	defer b.Unlock()
	b.callbacks = append(b.callbacks, cb)
	return b
}

func (b *LogBus) Dispatch(event LogEvent) {
	b.Lock()
	defer b.Unlock()

	for _, cb := range b.callbacks {
		func() {
			defer func() {
				if v := recover(); v != nil { // a panic is detected.
					b.log.WithRecover(v).Error("LogBus crashed", logan.F{
						"event": event,
					})
				}
			}()
			cb(event)
		}()
	}
}
