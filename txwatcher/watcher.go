package txwatcher

import (
	"time"

	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/swarmfund/api/internal/hose"
	"gitlab.com/tokend/horizon-connector"
)

type Watcher struct {
	horizon  *horizon.Connector
	log      *logan.Entry
	dispatch hose.TransactionDispatch
}

func NewWatcher(log *logan.Entry, connector *horizon.Connector, dispatch hose.TransactionDispatch) *Watcher {
	return &Watcher{
		dispatch: dispatch,
		horizon:  connector,
		log:      log,
	}
}

func (w *Watcher) Run() {
	// ticker to slow down requests leaving quota for other API requests
	// FIXME find a better way to prioritise requests from API
	ticker := time.NewTicker(3 * time.Second)
	defer func() {
		if rvr := recover(); rvr != nil {
			w.log.WithRecover(rvr).Error("watcher panicked")
		}
		ticker.Stop()
	}()
	events := make(chan horizon.TransactionEvent)
	errs := w.horizon.Listener().Transactions(events)
	for {
		select {
		case event := <-events:
			if event.Transaction != nil {
				w.log.WithFields(logan.F{
					"tx":     event.Transaction.PagingToken,
					"ledger": event.Transaction.CreatedAt,
				}).Debug("received tx")
				w.dispatch(event)
			}
			<-ticker.C
		case err := <-errs:
			w.log.WithError(err).Warn("failed to get transaction")
		}
	}
}
