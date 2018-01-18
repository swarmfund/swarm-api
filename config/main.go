package config

import (
	"sync"

	raven "github.com/getsentry/raven-go"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/swarmfund/api/internal/discourse"
	"gitlab.com/swarmfund/api/notificator"
	"gitlab.com/swarmfund/horizon-connector/v2"
)

type Config interface {
	Init() error
	API() API
	HTTP() HTTP
	Storage() Storage
	Log() *logan.Entry

	Notificator() notificator.ConnectorI
	Sentry() *raven.Client
	Horizon() *horizon.Connector
	Discourse() *discourse.Connector
	Get(string) map[string]interface{}
}

type ViperConfig struct {
	*viper.Viper
	*sync.RWMutex

	// runtime-initialized instances
	horizon     *horizon.Connector
	discourse   *discourse.Connector
	notificator notificator.ConnectorI
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

func (c *ViperConfig) Get(key string) map[string]interface{} {
	return c.Viper.GetStringMap(key)
}
