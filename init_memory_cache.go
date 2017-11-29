package api

import (
	"time"

	"github.com/patrickmn/go-cache"
)

func initMemoryCache(app *App) {
	app.memoryCache = cache.New(24*time.Hour, 1*time.Hour)
}

func init() {
	appInit.Add("memory_cache", initMemoryCache, "app-context", "log")
}
