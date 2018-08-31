package api

import (
	"gitlab.com/swarmfund/api/db2/api"
	"gitlab.com/swarmfund/api/internal/hose"
)

func init() {
	appInit.Add("audit-logs", func(app *App) {
		log := app.Config().Log().WithField("service", "audit-logs")

		app.logBus.Subscribe(func(event hose.LogEvent) {
			log.WithField("audit log event", event).Debug("processing")

			location, err := app.Config().GeoInfo().LocationInfo(event.User.IP)
			if err != nil {
				log.WithField("audit log event", event).Error(err, "failed to get location")
			}

			err = app.APIQ().AuditLog().Create(&api.AuditLogAction{
				ActionType: int32(event.Type),
				Details: api.LogInfo{
					GeoInfo: location,
				},
				PerformedAt: event.Time,
				UserAddress: string(event.User.Address),
			})
			if err != nil {
				log.WithField("audit log event", event).
					Error(err, "failed to add entry to db")
				return
			}

		})
	})
}
