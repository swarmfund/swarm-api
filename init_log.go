package api

import (
	"gitlab.com/swarmfund/api/db2"
	"gitlab.com/swarmfund/api/log"
)

// initLog initialized the logging subsystem, attaching app.log and
// app.logMetrics. It also configured the logger's level using Config.LogLevel.
func initLog(app *App) {
	log.DefaultLogger.Logger.Level = app.Config().Log().Level

	if app.Config().Log().SlowQueryBound != nil {
		db2.SlowQueryBound = *app.Config().Log().SlowQueryBound
	}
}

func init() {
	appInit.Add("log", initLog)
}
