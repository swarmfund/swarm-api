package txwatcher

import (
	"net/url"

	"gitlab.com/swarmfund/api/log"
	"gitlab.com/swarmfund/horizon-connector/v2"
)

type TxHandler func(horizon.TransactionEvent) error

type Watcher struct {
	horizonConn *horizon.Connector
	logger      *log.Entry
	handlers    []TxHandler
}

func NewWatcher(horizonURL *url.URL) *Watcher {
	return &Watcher{
		horizonConn: horizon.NewConnector(horizonURL),
		logger:      log.WithField("service", "tx_watcher"),
	}
}

func (w *Watcher) AddHandlers(ls ...TxHandler) {
	w.handlers = append(w.handlers, ls...)
}

func (w *Watcher) Run() {
	transactions := make(chan horizon.TransactionEvent)
	errs := w.horizonConn.Listener().Transactions(transactions)

	for {
		select {
		case tx := <-transactions:
			w.handleEvent(tx)
		case err := <-errs:
			w.logger.WithError(err).Warn("failed to get transaction")
		}
	}
}

func (w *Watcher) handleEvent(tx horizon.TransactionEvent) {
	for _, handler := range w.handlers {
		w.retryUntilSuccess(handler, tx)
	}
}

func (w *Watcher) retryUntilSuccess(handler TxHandler, tx horizon.TransactionEvent) {
	for {
		if err := handler(tx); err != nil {
			w.logger.WithError(err).Warn("handler failed")
			continue
		}
		return
	}
}
