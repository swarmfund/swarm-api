package config

import (
	"github.com/pkg/errors"
	"gitlab.com/distributed_lab/figure"
	"gitlab.com/swarmfund/api/assets"
	"gitlab.com/swarmfund/api/notificator"
)

const (
	notificatorConfigKey = "notificator"
)

func (c *ViperConfig) Notificator() *notificator.Connector {
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
		config.EmailConfirmation = assets.Templates.Lookup("email_confirm")
		c.notificator = notificator.NewConnector(config)
	}

	return c.notificator
}
