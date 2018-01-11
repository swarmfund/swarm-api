package api

import (
	"gitlab.com/swarmfund/api/internal/discourse"
	"gitlab.com/swarmfund/api/internal/hose"
)

func init() {
	appInit.Add("discourse-integration", func(app *App) {
		log := app.Config().Log().WithField("service", "user-create-listener")

		connector := app.Config().Discourse()

		// create user listener
		app.userBus.Subscribe(func(event hose.UserEvent) {
			if event.Type != hose.UserEventTypeCreated {
				return
			}
			err := connector.CreateUser(discourse.CreateUser{
				Email: event.User.Email,
			})
			entry := log.WithField("user", event.User.Address)
			if err != nil {
				entry.WithError(err).
					Error("failed to create discourse user")
				return
			}
			log.Debug("discourse user created")
		})
	})
}
