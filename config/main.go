package config

import (
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"gitlab.com/distributed_lab/logan/v3"
)

type Config interface {
	Init() error
	API() API
	HTTP() HTTP
	Storage() Storage
	Log() *logan.Entry
	Notificator() Notificator
	Get(string) map[string]interface{}
}

type ViperConfig struct {
	*viper.Viper
}

func NewViperConfig(fn string) Config {
	config := ViperConfig{
		Viper: viper.GetViper(),
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
