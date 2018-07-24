package api

import (
	"gitlab.com/swarmfund/api/txwatcher"
)

func initTxWatcher(app *App) {
	txWatcher := txwatcher.NewWatcher(
		app.Config().Log().WithField("service", "tx-watcher"),
		app.Config().Horizon(),
		app.txBus.Dispatch,
	)
	if !app.config.TxWatcher().Disabled {
		go txWatcher.Run()
	}

}

func init() {
	appInit.Add("tx-watcher", initTxWatcher, "bus-init")
}
