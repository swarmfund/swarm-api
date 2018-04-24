package middlewares

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/middleware"
	"gitlab.com/distributed_lab/logan/v3"
)

// TODO move to ape
func Logger(entry *logan.Entry, durationThreshold time.Duration) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
			t1 := time.Now()

			defer func() {
				dur := time.Since(t1)
				lEntry := entry.WithFields(logan.F{
					"path":     r.URL.Path,
					"duration": dur,
					"status":   ww.Status(),
				})
				lEntry.Info("request finished")

				if dur > durationThreshold {
					lEntry.WithField("http_request", r).Warn("slow request")
				}
			}()

			entry.WithField("path", r.URL.Path).Info("request started")
			next.ServeHTTP(ww, r)
		})
	}
}
