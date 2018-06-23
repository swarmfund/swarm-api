package config

import (
	"github.com/pkg/errors"
	"gitlab.com/distributed_lab/figure"
	"gitlab.com/swarmfund/api/internal/mixpanel"
)

func (c *ViperConfig) Mixpanel() *mixpanel.Connector {
	c.Lock()
	defer c.Unlock()

	if c.mixpanel == nil {
		var config struct {
			Token string
		}
		err := figure.
			Out(&config).
			From(c.GetStringMap("mixpanel")).
			Please()
		if err != nil {
			panic(errors.Wrap(err, "failed to figure out"))
		}
		c.mixpanel = mixpanel.New(config.Token)
	}

	return c.mixpanel
}
