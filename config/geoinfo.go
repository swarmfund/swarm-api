package config

import (
	"net/url"

	"github.com/pkg/errors"
	"gitlab.com/distributed_lab/figure"
	"gitlab.com/distributed_lab/geoinfo"
)

const geoInfoConfigKey = "geo_info"

type GeoConnector interface {
	LocationInfo(ip string) (*geoinfo.LocationInfo, error)
}

type DisabledGeoConnector struct{}

func (c *DisabledGeoConnector) LocationInfo(ip string) (*geoinfo.LocationInfo, error) {
	return &geoinfo.LocationInfo{}, nil
}

func (c *ViperConfig) GeoInfo() GeoConnector {
	c.Lock()
	defer c.Unlock()

	if c.geoinfo == nil {
		var config struct {
			AccessKey string   `fig:"access_key,required"`
			URL       *url.URL `fig:"api_url,required"`
			Disabled  bool     `fig:"disabled"`
		}

		err := figure.
			Out(&config).
			With(figure.BaseHooks).
			From(c.GetStringMap(geoInfoConfigKey)).
			Please()
		if err != nil {
			panic(errors.Wrap(err, "failed to figure out geo_info"))
		}

		c.geoinfo = &DisabledGeoConnector{}

		if !config.Disabled {
			c.geoinfo = geoinfo.NewConnector(config.AccessKey, config.URL)
		}
	}

	return c.geoinfo
}
