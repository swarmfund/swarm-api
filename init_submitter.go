package api

import "gitlab.com/swarmfund/api/pentxsub"

func initPendingSubmitter(app *App) {
	app.pendingSubmitter = pentxsub.New(app.APIQ().PenTXSub(), app.horizon, app.MasterSignerKP())
}

func init() {
	appInit.Add("pentxsub", initPendingSubmitter, "app-context", "api-db")
}
