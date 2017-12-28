package middlewares

import (
	"net/http"

	"time"

	"github.com/go-chi/chi/middleware"
	"gitlab.com/distributed_lab/logan/v3"
)

// TODO move to ape, once logan
func Logger(entry *logan.Entry) func(http.Handler) http.Handler {

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
			t1 := time.Now()
			defer func() {
				entry.WithFields(logan.F{
					"path":     r.URL.Path,
					"duration": time.Since(t1),
					"status":   ww.Status(),
				}).Info("request finished")
			}()
			entry.WithField("path", r.URL.Path).Info("request started")
			next.ServeHTTP(ww, r)
		})
	}
}
