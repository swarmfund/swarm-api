package middlewares

import (
	"net/http"
)

// TODO move to ape

func Recover(handle func(http.ResponseWriter, *http.Request, interface{})) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if rvr := recover(); rvr != nil {
					handle(w, r, rvr)
				}
			}()
			next.ServeHTTP(w, r)
		})
	}
}
