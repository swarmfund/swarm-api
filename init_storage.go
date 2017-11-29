package api

func initStorage(app *App) {
	//if app.config.Storage.DisableStorage {
	//	return
	//}
	//connector, err := storage.New(app.config.Storage)
	//if err != nil {
	//	log.WithField("service", "storage").WithError(err).Fatal("failed to init storage")
	//}
	//app.storage = connector
}

func init() {
	appInit.Add("storage", initStorage, "app-context", "log")
}
