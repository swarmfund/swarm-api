package config

import (
	"github.com/dukex/mixpanel"
	"github.com/pkg/errors"
	"gitlab.com/distributed_lab/figure"
)

func (c *ViperConfig) Mixpanel() mixpanel.Mixpanel {
	c.Lock()
	defer c.Unlock()

	if c.mixpanel == nil {
		var config struct {
			Token string
		}
		err := figure.
			Out(&config).
			From(c.Get("mixpanel")).
			Please()
		if err != nil {
			panic(errors.Wrap(err, "failed to figure out"))
		}
		c.mixpanel = mixpanel.New(config.Token, "")
	}

	return c.mixpanel
}
