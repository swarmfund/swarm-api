package api

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/patrickmn/go-cache"
	"gitlab.com/swarmfund/api/config"
	"gitlab.com/swarmfund/api/db2"
	"gitlab.com/swarmfund/api/db2/api"
	api2 "gitlab.com/swarmfund/api/internal/api"
	"gitlab.com/swarmfund/api/internal/data"
	horizon2 "gitlab.com/swarmfund/api/internal/data/horizon"
	"gitlab.com/swarmfund/api/internal/data/postgres"
	"gitlab.com/swarmfund/api/internal/hose"
	"gitlab.com/swarmfund/api/internal/track"
	"gitlab.com/swarmfund/api/log"
	"gitlab.com/swarmfund/api/storage"
	"gitlab.com/tokend/go/doorman"
	"gitlab.com/tokend/horizon-connector"
	"gitlab.com/tokend/keypair"
	"golang.org/x/net/context"
	"golang.org/x/net/http2"
	"gopkg.in/tylerb/graceful.v1"
	"gitlab.com/tokend/go/xdrbuild"
)

// App represents the root of the state of a horizon instance.
type App struct {
	// DEPRECATED
	CoreInfo *horizon.Info

	config         config.Config
	web            *Web
	apiQ           api.QInterface
	ctx            context.Context
	cancel         func()
	ticks          *time.Ticker
	horizonVersion string
	memoryCache    *cache.Cache
	storage        *storage.Connector
	// DEPRECATED
	horizon *horizon.Connector
	txBus   *hose.TransactionBus
	userBus *hose.UserBus
}

// NewApp constructs an new App instance from the provided config.
func NewApp(config config.Config) (*App, error) {
	result := &App{
		config:  config,
		horizon: config.Horizon(),
	}
	result.ticks = time.NewTicker(10 * time.Second)
	result.init()
	result.UpdateStellarCoreInfo()
	return result, nil
}

func (a *App) Config() config.Config {
	return a.config
}

func (a *App) MasterSignerKP() keypair.Full {
	return a.Config().API().AccountManager
}

func (a *App) MasterKP() keypair.Address {
	return keypair.MustParseAddress(a.CoreInfo.MasterAccountID)
}

func (a *App) EmailTokensQ() data.EmailTokensQ {
	return api.NewEmailTokensQ(a.APIRepo(a.ctx))
}

func (a *App) Blobs() data.Blobs {
	return postgres.NewBlobs(a.APIRepo(a.ctx))
}

func (a *App) Tracker() *track.Tracker {
	return track.NewTracker(a.Config().Log(), postgres.NewTracking(a.APIRepo(a.ctx)))
}

// Serve starts the horizon web server, binding it to a socket, setting up
// the shutdown signals.
func (a *App) Serve() {
	a.web.router.Compile()

	source := a.MasterKP()
	signer := a.MasterSignerKP()
	builder:= xdrbuild.
		NewBuilder(a.CoreInfoConn().GetPassphrase(), a.CoreInfoConn().GetTXExpire()).
		Transaction(source).
		Sign(signer)

	r := api2.Router(
		a.Config().Log().WithField("service", "api"),
		a.APIQ().Wallet(),
		a.EmailTokensQ(),
		a.APIQ().Users(),
		doorman.New(
			a.Config().API().SkipSignatureCheck,
			horizon2.NewAccountQ(a.Config().Horizon()),
		),
		a.horizon,
		a.APIQ().TFA(),
		a.config.Storage(),
		a.CoreInfoConn(),
		a.Blobs(),
		a.Config().Sentry(),
		a.userBus.Dispatch,
		a.Config().Notificator(),
		a.APIRepo(a.ctx),
		a.config.Wallets(),
		a.Tracker(),
		builder,
	)

	r.Mount("/", a.web.router)
	http.Handle("/", r)

	addr := fmt.Sprintf("%s:%d", a.config.HTTP().Host, a.config.HTTP().Port)

	srv := &graceful.Server{
		Timeout: 10 * time.Second,

		Server: &http.Server{
			Addr:    addr,
			Handler: http.DefaultServeMux,
		},

		ShutdownInitiated: func() {
			log.Info("received signal, gracefully stopping")
			a.Close()
		},
	}

	http2.ConfigureServer(srv.Server, nil)

	log.Infof("Starting horizon on %s", addr)

	go a.run()

	if err := srv.ListenAndServe(); err != nil {
		log.Panic(err)
	}

	log.Info("stopped")
}

// Close cancels the app and forces the closure of db connections
func (a *App) Close() {
	a.cancel()
	a.ticks.Stop()
}

// APIQ returns a helper object for performing sql queries against the
// horizon api database.
func (a *App) APIQ() api.QInterface {
	return a.apiQ
}

// APIRepo returns a new repo that loads data from the api database. The
// returned repo is bound to `ctx`.
func (a *App) APIRepo(ctx context.Context) *db2.Repo {
	return &db2.Repo{DB: a.apiQ.GetRepo().DB, Ctx: ctx}
}

// CoreInfoConn create new instance of coreinfo.Connector.
func (a *App) CoreInfoConn() data.CoreInfoI {
	info, err := a.Config().Horizon().Info()
	if err != nil {
		panic(err)
	}
	return info
}

// UpdateStellarCoreInfo updates the value of coreVersion and networkPassphrase
// from the Stellar core API.
func (a *App) UpdateStellarCoreInfo() {
	info, err := a.horizon.Info()
	if err != nil {
		log.WithField("service", "app").WithError(err).Warn("could not load stellar-core info")
		return
	}
	a.CoreInfo = info
}

// Tick triggers horizon to update all of it's background processes such as
// transaction submission, metrics, ingestion and reaping.
func (a *App) Tick() {
	var wg sync.WaitGroup
	log.Debug("ticking app")
	// update ledger state and stellar-core info in parallel
	wg.Add(1)

	go func() {
		defer func() {
			wg.Done()
		}()
		a.UpdateStellarCoreInfo()
	}()

	wg.Wait()

	log.Debug("finished ticking app")
}

// Init initializes app, using the config to populate db connections and
// whatnot.
func (a *App) init() {
	appInit.Run(a)
}

// run is the function that runs in the background that triggers Tick each
// second
func (a *App) run() {
	for {
		select {
		case <-a.ticks.C:
			a.Tick()
		case <-a.ctx.Done():
			log.Info("finished background ticker")
			return
		}
	}
}

/*func (c *Connector) runUpdater() {
	entry := log.WithField("service", "corer")
	var info *Info
	var err error

	for {
		time.Sleep(1 * time.Hour)

		info, err = c.getCoreInfo()
		if err != nil {
			entry.WithError(err).Error("unable to update core info")
			continue
		}

		coreInfo = info
		entry.Debug("core info updated")
	}*/
