package api

const (
	ServiceDocumentConsumer = "document-consumer"
)

func initConsumer(app *App) {
	//if app.config.Storage.DisableStorage {
	//	return
	//}
	//entry := log.WithField("service", ServiceDocumentConsumer)
	//consumer, err := storage.NewConsumer(app.config.Storage, app.APIQ(), app.Storage())
	//if err != nil {
	//	entry.WithError(err).Fatal("failed to init document consumer")
	//}
	//err = consumer.Prepare()
	//if err != nil {
	//	entry.WithError(err).Fatal("failed to prepare document consumer")
	//}
	//go consumer.Consume()
}

func init() {
	appInit.Add(ServiceDocumentConsumer, initConsumer, "app-context", "log")
}
