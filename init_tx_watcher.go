package api

import (
	"gitlab.com/swarmfund/api/txwatcher"
)

func initTxWatcher(app *App) {
	u := app.config.API().HorizonURL
	app.txWatcher = txwatcher.NewWatcher(&u)
	go app.txWatcher.Run()
}

func init() {
	appInit.Add("tx-watcher", initTxWatcher, "api-db")
}
