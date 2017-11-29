package config

import (
	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

type Config interface {
	Init() error
	API() API
	HTTP() HTTP
	Storage() Storage
	Log() Log
	Notificator() Notificator
}

type ViperConfig struct {
	*viper.Viper
	//apiDatabaseURL string
	//port           int
	//
	//LogLevel     logrus.Level
	//ClientRouter *url.URL
	//ClientDomain string
	//// SkipCheck disables signature validation, for testing and development purpose
	//SkipCheck     bool
	//Notificator   Notificator
	//TFA           TFA
	//Core          Core
	//Storage       Storage
	//HorizonURL    string
	//Deposit       Deposit
	//NoEmailVerify bool
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

func (c *ViperConfig) Storage() Storage {
	return Storage{}
}
