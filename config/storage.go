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
	Disabled             bool          `fig:"disable"`
	AccessKey            string        `fig:"access_key"`
	SecretKey            string        `fig:"secret_key"`
	Host                 string        `fig:"host"`
	ForceSSL             bool          `fig:"force_ssl"`
	FormDataExpire       time.Duration `fig:"form_data_expire"`
	MinContentLength     int64         `fig:"min_content_length"`
	MaxContentLength     int64         `fig:"max_content_length"`
	ObjectCreateARN      string        `fig:"object_create_arn"`
	ListenerBrokerURL    string        `fig:"listener_broker_url"`
	ListenerExchange     string        `fig:"listener_exchange"`
	ListenerExchangeType string        `fig:"listener_exchange_type"`
	ListenerBindingKey   string        `fig:"listener_binding_key"`
	MediaTypes           []string      `fig:"media_types"`
}

func (c *ViperConfig) Storage() Storage {
	c.Lock()
	defer c.Unlock()

	if c.storage != nil {
		return *c.storage
	}

	config := &Storage{}

	err := figure.Out(config).From(c.GetStringMap(storageConfigKey)).Please()
	if err != nil {
		panic(errors.Wrap(err, "failed to figure out storage"))
	}

	c.storage = config

	return *c.storage
}
