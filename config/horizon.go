package config

import (
	"net/url"

	"github.com/pkg/errors"
	"gitlab.com/distributed_lab/figure"
	"gitlab.com/swarmfund/horizon-connector/v2"
)

const (
	horizonConfigKey = "horizon"
)

func (c *ViperConfig) Horizon() *horizon.Connector {
	c.Lock()
	defer c.Unlock()

	if c.horizon == nil {
		var config struct {
			URL *url.URL
		}

		err := figure.
			Out(&config).
			From(c.GetStringMap(horizonConfigKey)).
			Please()
		if err != nil {
			panic(errors.Wrap(err, "failed to figure out horizon"))
		}

		c.horizon = horizon.NewConnector(config.URL)
	}

	return c.horizon
}
