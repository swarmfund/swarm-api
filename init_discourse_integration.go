package api

import (
	"fmt"

	"gitlab.com/swarmfund/api/internal/hose"
)

func init() {
	appInit.Add("discourse-integration", func(app *App) {
		//connector := app.Config().Discourse()

		// create user listener
		app.userBus.Subscribe(func(event hose.UserEvent) {
			fmt.Println("EVENT", event)
		})
	})
}
