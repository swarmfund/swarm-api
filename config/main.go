package config

import (
	"sync"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/getsentry/raven-go"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/swarmfund/api/internal/data"
	"gitlab.com/swarmfund/api/internal/discourse"
	"gitlab.com/swarmfund/api/internal/mixpanel"
	"gitlab.com/swarmfund/api/internal/salesforce"
	"gitlab.com/swarmfund/api/notificator"
	"gitlab.com/tokend/horizon-connector"
)

type Config interface {
	Init() error
	API() API
	HTTP() HTTP
	Storage() data.Storage
	Log() *logan.Entry
	Wallets() Wallets

	Notificator() *notificator.Connector
	Sentry() *raven.Client
	Horizon() *horizon.Connector
	Discourse() *discourse.Connector
	Mixpanel() *mixpanel.Connector
	Salesforce() *salesforce.Connector
}

//go:generate mockery -case underscore -name rawGetter -testonly -inpkg
// rawGetter encapsulates raw config values provider
type rawGetter interface {
	GetStringMap(key string) map[string]interface{}
}

type ViperConfig struct {
	rawGetter
	*sync.RWMutex

	// runtime-initialized instances
	horizon     *horizon.Connector
	discourse   *discourse.Connector
	notificator *notificator.Connector
	sentry      *raven.Client
	logan       *logan.Entry
	mixpanel    *mixpanel.Connector
	salesforce  *salesforce.Connector
	wallets     *Wallets
	storage     data.Storage
	api         *API
	aws         *session.Session
}

func NewViperConfig(fn string) Config {
	// init underlying viper
	v := viper.GetViper()
	v.SetConfigFile(fn)

	return newViperConfig(v)
}

func newViperConfig(raw rawGetter) Config {
	config := &ViperConfig{
		RWMutex: &sync.RWMutex{},
	}
	config.rawGetter = raw
	return config
}

func (c *ViperConfig) Init() error {
	if err := viper.ReadInConfig(); err != nil {
		return errors.Wrap(err, "failed to read config file")
	}
	return nil
}
