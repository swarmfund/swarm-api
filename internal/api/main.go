package api

import (
	"net/http"

	"time"

	"net/http/httputil"
	"net/url"

	raven "github.com/getsentry/raven-go"
	"github.com/go-chi/chi"
	"github.com/google/jsonapi"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/swarmfund/api/blacklist"
	"gitlab.com/swarmfund/api/config"
	"gitlab.com/swarmfund/api/db2"
	"gitlab.com/swarmfund/api/db2/api"
	"gitlab.com/swarmfund/api/internal/api/handlers"
	"gitlab.com/swarmfund/api/internal/api/middlewares"
	"gitlab.com/swarmfund/api/internal/data"
	"gitlab.com/swarmfund/api/internal/discourse/sso"
	"gitlab.com/swarmfund/api/internal/favorites"
	"gitlab.com/swarmfund/api/internal/hose"
	"gitlab.com/swarmfund/api/internal/salesforce"
	"gitlab.com/swarmfund/api/internal/secondfactor"
	"gitlab.com/swarmfund/api/notificator"
	"gitlab.com/tokend/go/doorman"
	"gitlab.com/tokend/horizon-connector"
)

func Router(
	entry *logan.Entry, walletQ api.WalletQI, tokensQ data.EmailTokensQ,
	usersQ api.UsersQI, auditLogsQ api.AuditLogQI, doorman doorman.Doorman, horizon *horizon.Connector,
	tfaQ api.TFAQI, storage data.Storage, coreInfo data.Info, blobQ data.Blobs, sentry *raven.Client,
	logDispatch hose.LogDispatch, notificator *notificator.Connector, repo *db2.Repo,
	wallets config.Wallets, salesforce *salesforce.Connector, txbuilder data.Infobuilder,
) chi.Router {
	r := chi.NewRouter()

	r.Use(
		ape.RecoverMiddleware(entry, sentry),
		secondfactor.HashMiddleware(),
		middlewares.Logger(entry, 300*time.Millisecond),
		middlewares.ContenType(jsonapi.MediaType),
		middlewares.Ctx(
			handlers.CtxWalletQ(walletQ),
			handlers.CtxLog(entry),
			handlers.CtxEmailTokensQ(tokensQ),
			handlers.CtxUsersQ(usersQ),
			handlers.CtxAuditLogs(auditLogsQ),
			handlers.CtxHorizon(horizon),
			handlers.CtxTransaction(txbuilder),
			handlers.CtxTFAQ(tfaQ),
			handlers.CtxDoorman(doorman),
			handlers.CtxStorage(storage),
			handlers.CtxCoreInfo(coreInfo),
			handlers.CtxLogBusDispatch(logDispatch),
			handlers.CtxNotificator(notificator),
			handlers.CtxWallets(wallets),
			handlers.CtxBlobQ(blobQ),
			handlers.CtxDomainApprover(blacklist.NewApprover(wallets.DomainsBlacklist...)),
			handlers.CtxSalesforce(salesforce),
		),
	)

	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		ape.RenderErr(w, problems.NotFound())
	})

	//participants
	r.Get("/user_id", handlers.GetUserId)
	r.Get("/data/enums", handlers.GetEnums)
	r.Post("/details", handlers.GetUsersDetails)
	r.Post("/participants", handlers.GetParticipants)

	// static stuff
	r.Get("/kdf", handlers.GetKDF)

	r.Route("/wallets", func(r chi.Router) {
		// admin endpoints
		r.Get("/", handlers.WalletsIndex)
		r.Delete("/{wallet-id}", handlers.DeleteWallets)

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

		r.Post("/{address}/events", handlers.SendEvent)

		// documents
		r.Route("/{address}/documents", func(r chi.Router) {
			r.Post("/", handlers.PutDocument)
			r.Get("/{document}", handlers.GetDocument)
		})

		// kyc
		r.Route("/{address}/entities", func(r chi.Router) {
			r.Post("/", handlers.CreateKYCEntity)
			r.Get("/", handlers.KYCEntitiesIndex)
			r.Put("/{entity}", handlers.PatchKYCEntity)
		})

		// DEPRECATED use /blobs family instead
		r.Route("/{address}/blobs", func(r chi.Router) {
			r.Post("/", handlers.CreateBlob)
			// DEPRECATED
			// at this point everything should have reference in core
			// and blobs index is not needed anymore
			r.Get("/", handlers.BlobIndex)
			r.Get("/{blob}", handlers.GetBlob)
		})

		// favorites
		r.Route("/{address}/favorites", favorites.Router(repo))

		// favorites aka settings
		r.Route("/{address}/settings", favorites.Router(repo))

		//get users statistics
		r.Get("/stats", handlers.UserStats)
	})

	r.Route("/v2/users", func(r chi.Router) {
		url, _ := url.Parse("http://localhost:7006")
		proxy := httputil.ReverseProxy{
			Director: func(request *http.Request) {
				request.URL.Scheme = url.Scheme
				request.URL.Host = url.Host
				request.URL.Path = url.Path
			},
		}
		r.Get("/", proxy.ServeHTTP)
	})

	// blobs
	r.Route("/blobs", func(r chi.Router) {
		r.Post("/", handlers.CreateBlob)
		r.Get("/{blob}", handlers.GetBlob)
		r.Delete("/{blob}", handlers.DeleteBlob)
	})

	r.Route("/documents", func(r chi.Router) {
		r.Post("/", handlers.PutDocument)
		r.Get("/{document}", handlers.GetDocument)
	})

	r.Route("/integrations", func(r chi.Router) {
		// discourse ping-pong
		r.Get("/discourse-sso", sso.SSOReceiver)
		r.Post("/discourse-sso", sso.SSORedirect)
	})

	// guest-by-email favorites
	r.Route("/favorites", favorites.Router(repo))
	return r
}
