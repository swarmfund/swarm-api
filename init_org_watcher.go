package api

func initOrganizationWatcher(app *App) {
	//u, err := url.Parse(app.config.HorizonURL)
	//if err != nil {
	//	log.WithField("service", "org-watcher").WithError(err).Error("failed to init org watcher")
	//}
	//
	//watcher := orgwatcher.New(
	//	app.AccountManagerKP(),
	//	log.WithField("service", "org-watcher"),
	//	app.APIQ(),
	//	*u,
	//)
	//go watcher.Run()
}

func init() {
	appInit.Add("org-watcher", initOrganizationWatcher, "api-db")
}
