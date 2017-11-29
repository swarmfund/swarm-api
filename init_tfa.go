package api

func initTFASMSBAckend(app *App) {
	//tfa.DefaultSMSConfig(tfa.SMSConfig{
	//	Dev:         app.config.TFA.Dev,
	//	Notificator: app.config.Notificator.Endpoint,
	//})
}

func init() {
	appInit.Add("tfa-sms-backend", initTFASMSBAckend)
}
