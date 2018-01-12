package txwatcher

import (
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/swarmfund/api/internal/hose"
	"gitlab.com/swarmfund/horizon-connector/v2"
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
	defer func() {
		// TODO recover
	}()
	events := make(chan horizon.TransactionEvent)
	errs := w.horizon.Listener().Transactions(events)

	for {
		select {
		case event := <-events:
			if event.Transaction != nil {
				w.log.WithField("tx", event.Transaction.PagingToken).Debug("received tx")
				w.dispatch(event)
			}
		case err := <-errs:
			w.log.WithError(err).Warn("failed to get transaction")
		}
	}
}