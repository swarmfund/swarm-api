package api

import (
	"database/sql"

	"github.com/rcrowley/go-metrics"
	"github.com/zenazn/goji/web"
	"github.com/zenazn/goji/web/middleware"
	"gitlab.com/swarmfund/api/internal/secondfactor"
	"gitlab.com/swarmfund/api/render/problem"
)

// Web contains the http server related fields for horizon: the router,
// rate limiter, etc.
type Web struct {
	router *web.Mux

	requestTimer metrics.Timer
	failureMeter metrics.Meter
	successMeter metrics.Meter
}

// initWeb installed a new Web instance onto the provided app object.
func initWeb(app *App) {
	app.web = &Web{
		router:       web.New(),
		requestTimer: metrics.NewTimer(),
		failureMeter: metrics.NewMeter(),
		successMeter: metrics.NewMeter(),
	}

	// register problems
	problem.RegisterError(sql.ErrNoRows, problem.NotFound)
}

// initWebMiddleware installs the middleware stack used for horizon onto the
// provided app.
func initWebMiddleware(app *App) {
	r := app.web.router
	r.Use(secondfactor.HashMiddleware())
	r.Use(stripTrailingSlashMiddleware())
	r.Use(middleware.EnvInit)
	r.Use(app.Middleware)
	r.Use(middleware.RequestID)
	r.Use(contextMiddleware(app.ctx))
	r.Use(LoggerMiddleware)
	r.Use(requestMetricsMiddleware)
	r.Use(RecoverMiddleware)
	r.Use(middleware.AutomaticOptions)
}

// initWebActions installs the routing configuration of horizon onto the
// provided app.  All route registration should be implemented here.
func initWebActions(app *App) {
	// ok

	r := app.web.router

	// participants
	r.Post("/details", &DetailsAction{})
	r.Post("/participants", &ParticipantsAction{})

	// user actions
	//r.Post("/users", &CreateUserAction{})
	r.Get("/users", &UserIndexAction{})
	r.Get("/users/:id", &UserShowAction{})

	// kyc
	r.Patch("/users/:id", &PatchUserAction{})
	r.Post("/users/:user/entities", &CreateKYCEntityAction{})
	r.Patch("/users/:user/entities/:entity", &PatchKYCEntityAction{})
	r.Delete("/users/:user/entities/:entity", &DeleteKYCEntityAction{})
	r.Post("/users/:user/approve", &UserApproveAction{})

	//r.Post("/users/unverified/delete", &DeleteWalletAction{})
	//r.Post("/users/unverified/resend_token", &ResendTokenAction{})

	// documents
	r.Get("/users/:id/documents", &GetUserDocsAction{})
	r.Get("/users/:id/documents/:version", &GetUserFileAction{})
	r.Post("/users/:id/documents", &PutDocumentAction{})
	r.Get("/user_id", &GetUserIdAction{})

	// wallet
	r.Get("/wallets/unverified", &GetUnverifiedWalletsAction{})
	//r.Get("/wallets/verify", &VerifyWalletAction{})

	//r.Post("/wallets/create", &CreateWalletAction{})
	//r.Post("/wallets/update", &UpdateWalletAction{})
	//r.Post("/wallets/get_tfa_secret", &GetTfaKeychainAction{})

	//r.Post("/wallets/show", &ShowWalletAction{})
	//r.Post("/wallets/show_login_params", &ShowLoginParamsAction{})
	r.Patch("/wallets/:id", &PatchWalletAction{})
	r.Get("/wallets/:id/organization", &GetWalletOrganizationAction{})

	// wallet recovery
	//   user endpoints
	r.Post("/wallets/recovery", &CreateRecoveryRequestAction{})
	r.Get("/wallets/recovery", &ApproveRecoveryRequestAction{})
	//   admin endpoints
	r.Get("/recoveries", &GetRecoveryRequestsAction{})
	r.Get("/recoveries/:id", &GetRecoveryRequestAction{})
	r.Post("/users/:id/recovery", &ResolveUserRecoveryRequestAction{})

	// 2fa
	//r.Get("/tfa", &GetTFAAction{})
	//r.Post("/tfa", &EnableTFABackendAction{})
	//r.Patch("/tfa/:tfa", &UpdateTFABackendAction{})
	//r.Get("/tfa/verify", &VerifyTFAAction{})

	//   admin endpoint
	//r.Post("/tfa/delete", &DeleteTFABackendsAction{})

	// limit review
	r.Get("/users/:id/poi", &UserProofOfIncomeIndexAction{})
	r.Post("/users/:id/poi/:version", &UserProofOfIncomeApproveAction{})

	r.Get("/data/enums", &GetEnumsAction{})

	// deposit
	r.Get("/deposit/:method", &DepositInfoAction{})
	r.Post("/deposit", &DepositAction{})

	// transaction submission
	r.Post("/transactions", &PendingTransactionCreateAction{})

	// pending transactions
	r.Get("/accounts/:id/transactions", &PendingTransactionIndexAction{})
	r.Patch("/transactions/:tx_hash", &PendingTransactionRejectAction{})
	r.Delete("/transactions/:tx_hash", &PendingTransactionDeleteAction{})

	r.Get("/notifications/:id", &GetNotificationsAction{})
	r.Patch("/notifications/:id", &PatchNotificationsAction{})

	//r.Get("/kdf_params", &KdfParamsAction{})

	r.NotFound(&NotFoundAction{})
}

func init() {
	appInit.Add(
		"web.init",
		initWeb,
		"app-context", "stellarCoreInfo", "api-db", "memory_cache",
	)

	appInit.Add(
		"web.middleware",
		initWebMiddleware,

		"web.init",
	)
	appInit.Add(
		"web.actions",
		initWebActions,

		"web.init",
	)
}
