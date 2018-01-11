package api

import "gitlab.com/swarmfund/api/internal/hose"

func init() {
	appInit.Add("bus-init", func(app *App) {
		app.txBus = hose.NewTransactionBus()
		app.userBus = hose.NewUserBus()
	})
}
