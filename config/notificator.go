package config

import (
	"github.com/pkg/errors"
	"gitlab.com/distributed_lab/figure"
	"gitlab.com/swarmfund/api/notificator"
)

const (
	notificatorConfigKey = "notificator"
)

//Warning: Do not call the Notificator directly,
//use instead of it the notificator that is in api.API
func (c *ViperConfig) Notificator() *notificator.Connector {
	horizonConnect := c.Horizon()
	c.Lock()
	defer c.Unlock()

	if c.notificator == nil {
		config := notificator.Config{}
		err := figure.
			Out(&config).
			From(c.GetStringMap(notificatorConfigKey)).
			Please()
		if err != nil {
			panic(errors.Wrap(err, "failed to figure out notificator"))
		}

		notificator := notificator.NewConnector(config)
		if err := notificator.Init(horizonConnect); err != nil {
			panic(errors.Wrap(err, "failed to init notificator"))
		}
		c.notificator = notificator
	}

	return c.notificator
}
