package api

import (
	"gitlab.com/swarmfund/api/txwatcher"
)

func initTxWatcher(app *App) {
	app.txWatcher = txwatcher.NewWatcher(
		app.Config().Log().WithField("service", "tx-watcher"),
		app.Config().Horizon(),
		app.txBus.Dispatch,
	)
	go app.txWatcher.Run()
}

func init() {
	appInit.Add("tx-watcher", initTxWatcher, "bus-init")
}
