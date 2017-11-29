package middlewares

import (
	"net/http"

	"github.com/go-chi/chi"
	"gitlab.com/swarmfund/api/internal/api/movetoape"
	"gitlab.com/swarmfund/go/doorman"
	"gitlab.com/swarmfund/go/doorman/types"
)

func CheckAllowed(key string, checks ...func(string) types.SignerConstraint) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			value := chi.URLParam(r, key)
			constraints := make([]types.SignerConstraint, 0, len(checks))
			for _, check := range checks {
				constraints = append(constraints, check(value))
			}
			if err := doorman.Check(r, constraints...); err != nil {
				movetoape.RenderDoormanErr(w, err)
				return
			}
			next.ServeHTTP(w, r)
		}
		return http.HandlerFunc(fn)
	}
}
