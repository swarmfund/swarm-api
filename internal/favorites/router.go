package favorites

import (
	"github.com/go-chi/chi"
	"gitlab.com/swarmfund/api/db2"
	"gitlab.com/swarmfund/api/internal/api/middlewares"
	"gitlab.com/swarmfund/api/internal/favorites/internal/data/postgres"
	"gitlab.com/swarmfund/api/internal/favorites/internal/handlers"
)

func Router(repo *db2.Repo) func(chi.Router) {
	return func(r chi.Router) {
		r.Use(middlewares.Ctx(handlers.CtxFavoritesQ(postgres.NewFavorites(repo))))
		r.Post("/", handlers.CreateFavorite)
		r.Delete("/{favorite}", handlers.DeleteFavorite)
		r.Get("/", handlers.FavoriteIndex)
		//TODO MOVE ME to api.Handlers
		r.Post("/notifications/{key}", handlers.SendNotifications)
	}
}
