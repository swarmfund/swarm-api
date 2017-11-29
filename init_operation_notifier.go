package api

func initOperationNotifier(app *App) {
	//if app.config.Notificator.OperationsNotification.Disable {
	//	return
	//}
	//
	//notifier := operation_notifier.New(
	//	app.config.Notificator.OperationsNotification.LastCursor,
	//	app.config.HorizonURL,
	//	app.AccountManagerKP(),
	//	app.notificator,
	//	app.apiQ,
	//)
	//
	//notifier.NewListener()
	//go notifier.Run()
}

func init() {
	appInit.Add("operation-notifier", initOperationNotifier, "api-db")
}
