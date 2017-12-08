package config

import (
	"time"

	"github.com/pkg/errors"
	"gitlab.com/distributed_lab/figure"
)

const (
	storageConfigKey = "storage"
)

type Storage struct {
	Disabled             bool
	AccessKey            string
	SecretKey            string
	Host                 string
	ForceSSL             bool
	FormDataExpire       time.Duration
	MinContentLength     int64
	MaxContentLength     int64
	ObjectCreateARN      string
	ListenerBrokerURL    string
	ListenerExchange     string
	ListenerExchangeType string
	ListenerBindingKey   string
}

func (c *ViperConfig) Storage() (config Storage) {
	err := figure.Out(&config).From(c.GetStringMap(storageConfigKey)).Please()
	if err != nil {
		panic(errors.Wrap(err, "failed to figure out storage"))
	}
	return config
}
