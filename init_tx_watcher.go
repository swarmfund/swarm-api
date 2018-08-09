package api

import (
	"gitlab.com/swarmfund/api/txwatcher"
)

func initTxWatcher(app *App) {
	if app.config.TXWatcher().Disabled {
		return
	}
	go txwatcher.NewWatcher(
		app.Config().Log().WithField("service", "tx-watcher"),
		app.Config().Horizon(),
		app.txBus.Dispatch,
	).Run()
}

func init() {
	appInit.Add("tx-watcher", initTxWatcher)
}
