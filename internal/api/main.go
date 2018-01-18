package api

import (
	"net/http"

	"time"

	"github.com/getsentry/raven-go"
	"github.com/go-chi/chi"
	"github.com/google/jsonapi"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/swarmfund/api/db2/api"
	"gitlab.com/swarmfund/api/internal/api/handlers"
	"gitlab.com/swarmfund/api/internal/api/middlewares"
	"gitlab.com/swarmfund/api/internal/data"
	"gitlab.com/swarmfund/api/internal/discourse/sso"
	"gitlab.com/swarmfund/api/internal/hose"
	"gitlab.com/swarmfund/api/internal/secondfactor"
	"gitlab.com/swarmfund/api/notificator"
	"gitlab.com/swarmfund/api/storage"
	"gitlab.com/swarmfund/go/doorman"
	"gitlab.com/swarmfund/horizon-connector"
	"gitlab.com/tokend/keypair"
)

func Router(
	entry *logan.Entry, walletQ api.WalletQI, tokensQ data.EmailTokensQ,
	usersQ api.UsersQI, doorman doorman.Doorman, horizon *horizon.Connector,
	tfaQ api.TFAQI, storage *storage.Connector, master keypair.Address, signer keypair.Full,
	coreInfo data.CoreInfoI, blobQ data.Blobs, sentry *raven.Client,
	userDispatch hose.UserDispatch, authorizedDeviceQ api.AuthorizedDeviceQI,
	notificatorConn notificator.ConnectorI,
) chi.Router {
	r := chi.NewRouter()

	r.Use(
		ape.RecoverMiddleware(entry, sentry),
		secondfactor.HashMiddleware(),
		middlewares.Logger(entry, 200*time.Second),
		middlewares.ContenType(jsonapi.MediaType),
		middlewares.Ctx(
			handlers.CtxWalletQ(walletQ),
			handlers.CtxLog(entry),
			handlers.CtxEmailTokensQ(tokensQ),
			handlers.CtxUsersQ(usersQ),
			handlers.CtxHorizon(horizon),
			handlers.CtxTransaction(master, signer),
			handlers.CtxTFAQ(tfaQ),
			handlers.CtxDoorman(doorman),
			handlers.CtxStorage(storage),
			handlers.CtxCoreInfo(coreInfo),
			handlers.CtxUserBusDispatch(userDispatch),
			handlers.CtxAuthorizedDeviceQ(authorizedDeviceQ),
			handlers.CtxNotificator(notificatorConn),
		),
	)

	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		ape.RenderErr(w, problems.NotFound())
	})

	// static stuff
	r.Get("/kdf", handlers.GetKDF)

	r.Route("/wallets", func(r chi.Router) {
		// admin endpoints
		r.Get("/", handlers.WalletsIndex)

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

	r.Route("/users", func(r chi.Router) {
		r.Get("/", handlers.UsersIndex)
		r.Get("/{address}", handlers.GetUser)
		r.Put("/{address}", handlers.CreateUser)
		r.Patch("/{address}", handlers.PatchUser)

		// documents
		r.Route("/{address}/documents", func(r chi.Router) {
			r.Post("/", handlers.PutDocument)
		})

		// kyc
		r.Route("/{address}/entities", func(r chi.Router) {
			r.Post("/", handlers.CreateKYCEntity)
			r.Get("/", handlers.KYCEntitiesIndex)
			r.Put("/{entity}", handlers.PatchKYCEntity)
		})

		// blobs
		r.Route("/{address}/blobs", func(r chi.Router) {
			r.Use(middlewares.Ctx(
				handlers.CtxBlobQ(blobQ),
			))
			r.Post("/", handlers.CreateBlob)
			r.Get("/", handlers.BlobIndex)
			r.Get("/{blob}", handlers.GetBlob)
		})
	})

	r.Route("/integrations", func(r chi.Router) {
		// discourse ping-pong
		r.Get("/discourse-sso", sso.SSOReceiver)
		r.Post("/discourse-sso", sso.SSORedirect)
	})

	return r
}
