package handlers

import (
	"context"
	"net/http"

	"gitlab.com/swarmfund/api/internal/favorites/internal/data"
)

type ctxKey int

const (
	favoritesQCtx ctxKey = iota
)

func CtxFavoritesQ(q data.Favorites) func(context.Context) context.Context {
	return func(ctx context.Context) context.Context {
		return context.WithValue(ctx, favoritesQCtx, q)
	}
}

func FavoritesQ(r *http.Request) data.Favorites {
	return r.Context().Value(favoritesQCtx).(data.Favorites)
}
