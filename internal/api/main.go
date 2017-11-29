package api

import (
	"net/http"
	"runtime/debug"

	"github.com/go-chi/chi"
	"github.com/google/jsonapi"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"gitlab.com/swarmfund/api/db2/api"
	"gitlab.com/swarmfund/api/errors"
	"gitlab.com/swarmfund/api/internal/api/handlers"
	"gitlab.com/swarmfund/api/internal/api/middlewares"
	"gitlab.com/swarmfund/api/internal/data"
	"gitlab.com/swarmfund/api/internal/secondfactor"
	"gitlab.com/swarmfund/api/log"
	"gitlab.com/swarmfund/go/doorman"
	"gitlab.com/swarmfund/go/keypair"
	"gitlab.com/swarmfund/horizon-connector"
)

func Router(
	entry *log.Entry, walletQ api.WalletQI, tokensQ data.EmailTokensQ,
	usersQ api.UsersQI, doorman doorman.Doorman, horizon *horizon.Connector,
	accountManager keypair.KP, tfaQ api.TFAQI,
) chi.Router {
	r := chi.NewRouter()

	r.Use(
		middlewares.Recover(func(w http.ResponseWriter, r *http.Request, rvr interface{}) {
			if entry != nil {
				entry.WithField("stack", string(debug.Stack())).
					WithError(errors.FromPanic(rvr)).Error("handler panicked")
			}
			ape.RenderErr(w, problems.InternalError())
		}),
		secondfactor.HashMiddleware(),
		middlewares.Logger(entry),
		middlewares.ContenType(jsonapi.MediaType),
		middlewares.Ctx(
			handlers.CtxWalletQ(walletQ),
			handlers.CtxLog(entry),
			handlers.CtxEmailTokensQ(tokensQ),
			handlers.CtxUsersQ(usersQ),
			handlers.CtxHorizon(horizon),
			handlers.CtxAccountManagerKP(accountManager),
			handlers.CtxTFAQ(tfaQ),
			handlers.CtxDoorman(doorman),
		),
	)
	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		ape.RenderErr(w, problems.NotFound())
	})

	// static stuff
	r.Get("/kdf", handlers.GetKDF)

	r.Route("/wallets", func(r chi.Router) {
		// signup
		r.Post("/", handlers.CreateWallet)

		// email verification
		r.Post("/{wallet-id}/verification", handlers.RequestVerification)
		r.Put("/{wallet-id}/verification", handlers.VerifyWallet)

		// login
		r.Get("/kdf", handlers.GetKDF)
		r.Get("/{wallet-id}", handlers.GetWallet)

		// change password
		r.Put("/{wallet-id}", handlers.ChangeWalletID)

		// 2fa
		r.Route("/{wallet-id}/factors", func(r chi.Router) {
			r.Post("/", handlers.CreateTFABackend)
			r.Get("/", handlers.GetWalletFactors)
			r.Delete("/{backend}", handlers.DeleteWalletFactor)
			r.Patch("/{backend}", handlers.UpdateWalletFactor)
			r.Put("/{backend}/verification", handlers.VerifyFactorOTP)
		})
	})

	r.Route("/users/{address}", func(r chi.Router) {
		//r.Use(middlewares.CheckAllowed("address", doorman.SignerOf))
		r.Put("/", handlers.CreateUser)
	})
	return r
}
