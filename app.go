package api

import (
	"fmt"
	"net/http"
	"time"

	"gitlab.com/swarmfund/api/config"
	"gitlab.com/swarmfund/api/db2/api"
	api2 "gitlab.com/swarmfund/api/internal/api"
	"gitlab.com/swarmfund/api/internal/data"
	horizon2 "gitlab.com/swarmfund/api/internal/data/horizon"
	"gitlab.com/swarmfund/api/internal/data/postgres"
	"gitlab.com/swarmfund/api/internal/hose"
	"gitlab.com/tokend/go/doorman"
	"gitlab.com/tokend/go/support/log"
	"gitlab.com/tokend/go/xdrbuild"
	"gitlab.com/tokend/keypair"
	"golang.org/x/net/context"
	"golang.org/x/net/http2"
	"gopkg.in/tylerb/graceful.v1"
)

// App represents the root of the state of a horizon instance.
type App struct {
	config  config.Config
	ctx     context.Context
	cancel  func()
	infoer  data.Info
	storage data.Storage
	txBus   *hose.TransactionBus
	logBus  *hose.LogBus
}

// NewApp constructs an new App instance from the provided config.
func NewApp(config config.Config) (*App, error) {
	ctx, cancel := context.WithCancel(context.Background())
	result := &App{
		config: config,
		ctx:    ctx,
		cancel: cancel,
		txBus:  hose.NewTransactionBus(config.Log().WithField("service", "tx-bus")),
		logBus: hose.NewLogBus(config.Log().WithField("service", "audit-log-bus")),
		infoer: NewLazyInfo(ctx, config.Log(), config.Horizon()),
	}
	result.init()
	return result, nil
}

func (a *App) Config() config.Config {
	return a.config
}

func (a *App) EmailTokensQ() data.EmailTokensQ {
	return api.NewEmailTokensQ(a.Config().DB())
}

func (a *App) Blobs() data.Blobs {
	return postgres.NewBlobs(a.Config().DB())
}

// Serve starts the horizon web server, binding it to a socket, setting up
// the shutdown signals.
func (a *App) Serve() {
	builder := func(info data.Info) *xdrbuild.Transaction {
		inf, err := info.Info()
		if err != nil {
			//TODO handle error
			panic(err)
		}

		source := keypair.MustParseAddress(inf.GetMasterAccountID())
		signer := a.Config().API().AccountManager
		return xdrbuild.
			NewBuilder(inf.GetPassphrase(), inf.GetTXExpire()).
			Transaction(source).
			Sign(signer)
	}

	r := api2.Router(
		a.Config().Log().WithField("service", "api"),
		a.APIQ().Wallet(),
		a.EmailTokensQ(),
		a.APIQ().Users(),
		a.APIQ().AuditLog(),
		doorman.New(
			a.Config().API().SkipSignatureCheck,
			horizon2.NewAccountQ(a.Config().Horizon()),
		),
		a.Config().Horizon(),
		a.APIQ().TFA(),
		a.config.Storage(),
		a.Info(),
		a.Blobs(),
		a.Config().Sentry(),
		a.logBus.Dispatch,
		a.Config().Notificator(),
		a.Config().DB(),
		a.config.Wallets(),
		a.Config().Salesforce(),
		builder,
	)

	http.Handle("/", r)

	addr := fmt.Sprintf("%s:%d", a.config.HTTP().Host, a.config.HTTP().Port)

	srv := &graceful.Server{
		Timeout: 10 * time.Second,
		Server: &http.Server{
			Addr:    addr,
			Handler: http.DefaultServeMux,
		},
		ShutdownInitiated: func() {
			a.Close()
		},
	}

	http2.ConfigureServer(srv.Server, nil)

	if err := srv.ListenAndServe(); err != nil {
		log.Panic(err)
	}
}

// Close cancels the app and forces the closure of db connections
func (a *App) Close() {
	a.cancel()
}

// APIQ returns a helper object for performing sql queries against the
// horizon api database.
func (a *App) APIQ() api.QInterface {
	return &api.Q{a.Config().DB()}
}

// CoreInfoConn create new instance of coreinfo.Connector.
func (a *App) Info() data.Info {
	return a.infoer
}

// Init initializes app, using the config to populate db connections and
// whatnot.
func (a *App) init() {
	appInit.Run(a)
}
