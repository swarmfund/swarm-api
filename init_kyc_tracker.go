package api

func initKYCTracker(app *App) {
	//tracker := kyctracker.New(
	//	log.WithField("service", "kyc-tracker"),
	//	app.APIQ(),
	//	app.notificator,
	//)
	//go tracker.Run()
}

func init() {
	appInit.Add("kyc-tracker", initKYCTracker, "api-db")
}
