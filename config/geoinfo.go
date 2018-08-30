package config

import (
	"net/url"

	"github.com/pkg/errors"
	"gitlab.com/distributed_lab/figure"
	"gitlab.com/swarmfund/api/geoinfo"
)

const geoInfoConfigKey = "geo_info"

func (c *ViperConfig) GeoInfo() *geoinfo.Connector {
	c.Lock()
	defer c.Unlock()

	if c.geoinfo == nil {
		var config struct {
			AccessKey string   `fig:"access_key,required"`
			URL       *url.URL `fig:"api_url,required"`
		}

		err := figure.
			Out(&config).
			With(figure.BaseHooks).
			From(c.GetStringMap(geoInfoConfigKey)).
			Please()
		if err != nil {
			panic(errors.Wrap(err, "failed to figure out geo_info"))
		}

		c.geoinfo = geoinfo.NewConnector(config.AccessKey, config.URL)
	}

	return c.geoinfo
}
