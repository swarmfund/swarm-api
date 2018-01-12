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

	r.Get("/user_id", &GetUserIdAction{})

	// limit review
	r.Get("/users/:id/poi", &UserProofOfIncomeIndexAction{})
	r.Post("/users/:id/poi/:version", &UserProofOfIncomeApproveAction{})

	r.Get("/data/enums", &GetEnumsAction{})

	// transaction submission
	r.Post("/transactions", &PendingTransactionCreateAction{})

	// pending transactions
	r.Get("/accounts/:id/transactions", &PendingTransactionIndexAction{})
	r.Patch("/transactions/:tx_hash", &PendingTransactionRejectAction{})
	r.Delete("/transactions/:tx_hash", &PendingTransactionDeleteAction{})

	r.Get("/notifications/:id", &GetNotificationsAction{})
	r.Patch("/notifications/:id", &PatchNotificationsAction{})

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
