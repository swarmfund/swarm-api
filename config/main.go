package config

import (
	"sync"

	raven "github.com/getsentry/raven-go"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/swarmfund/api/internal/discourse"
	"gitlab.com/swarmfund/api/internal/mixpanel"
	"gitlab.com/swarmfund/api/notificator"
	"gitlab.com/tokend/horizon-connector"
)

type Config interface {
	Init() error
	API() API
	HTTP() HTTP
	Storage() Storage
	Log() *logan.Entry
	Wallets() Wallets

	Notificator() *notificator.Connector
	Sentry() *raven.Client
	Horizon() *horizon.Connector
	Discourse() *discourse.Connector
	Mixpanel() *mixpanel.Connector

	Get(string) map[string]interface{}
}

type ViperConfig struct {
	*viper.Viper
	*sync.RWMutex

	// runtime-initialized instances
	horizon     *horizon.Connector
	discourse   *discourse.Connector
	notificator *notificator.Connector
	sentry      *raven.Client
	logan       *logan.Entry
	mixpanel    *mixpanel.Connector
	wallets     *Wallets
}

func NewViperConfig(fn string) Config {
	config := ViperConfig{
		Viper:   viper.GetViper(),
		RWMutex: &sync.RWMutex{},
	}
	config.SetConfigFile(fn)
	return &config
}

func (c *ViperConfig) Init() error {
	if err := viper.ReadInConfig(); err != nil {
		return errors.Wrap(err, "failed to read config file")
	}
	return nil
}

// Get will return value associated with config key, empty map if key is missing
func (c *ViperConfig) Get(key string) map[string]interface{} {
	m := c.Viper.GetStringMap(key)
	if m == nil {
		m = map[string]interface{}{}
	}
	return m
}
