package config

import (
	"net/url"

	"github.com/pkg/errors"
	"gitlab.com/distributed_lab/figure"
	"gitlab.com/swarmfund/api/internal/discourse"
)

const (
	discourseConfigKey = "discourse"
)

func (c *ViperConfig) Discourse() *discourse.Connector {
	c.Lock()
	defer c.Unlock()

	if c.discourse == nil {
		var config struct {
			URL      url.URL
			Username string
			Key      string
		}

		err := figure.
			Out(&config).
			With(URLHook, figure.BaseHooks).
			From(c.GetStringMap(discourseConfigKey)).
			Please()
		if err != nil {
			panic(errors.Wrap(err, "failed to figure out discourse"))
		}

		c.discourse = discourse.NewConnector(&config.URL, config.Username, config.Key)
	}

	return c.discourse
}
