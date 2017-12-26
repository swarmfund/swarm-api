package middlewares

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/middleware"
	"gitlab.com/swarmfund/api/log"
)

var SlowRequestBound time.Duration

func init() {
	SlowRequestBound = time.Second
}

// TODO move to ape, once logan
func Logger(entry *log.Entry) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
			t1 := time.Now()

			defer func() {
				dur := time.Since(t1)
				lEntry := entry.WithFields(log.F{
					"path":     r.URL.Path,
					"duration": dur,
					"status":   ww.Status(),
				})

				lEntry.Info("request finished")
				if dur > SlowRequestBound {
					lEntry.Warning("too slow request")
				}
			}()

			entry.WithField("path", r.URL.Path).Info("request started")
			next.ServeHTTP(ww, r)
		})
	}
}
