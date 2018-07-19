package api

import (
	"gitlab.com/swarmfund/api/db2"
	"gitlab.com/swarmfund/api/db2/api"
)

func initAPIDB(app *App) {
	repo, err := db2.Open(app.config.API().DatabaseURL)

	if err != nil {
		panic(err)
	}
	repo.DB.SetMaxIdleConns(4)
	repo.DB.SetMaxOpenConns(12)

	app.apiQ = &api.Q{Repo: repo}
}

func init() {
	appInit.Add("api-db", initAPIDB, "app-context")
}
