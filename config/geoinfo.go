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

type disabledGeoConnector struct{}

func (c *disabledGeoConnector) LocationInfo(ip string) (*geoinfo.LocationInfo, error) {
	return &geoinfo.LocationInfo{}, nil
}

func (c *ViperConfig) GeoInfo() GeoConnector {
	c.Lock()
	defer c.Unlock()

	if c.geoinfo == nil {
		// check if geo_info is disabled
		var disabled struct {
			Disabled bool `fig:"disabled"`
		}

		configData := c.GetStringMap(geoInfoConfigKey)

		err := figure.
			Out(&disabled).
			From(configData).
			Please()
		if err != nil {
			panic(errors.Wrap(err, "failed to figure out geo_info disabled"))
		}

		if disabled.Disabled {
			c.geoinfo = &disabledGeoConnector{}
			return c.geoinfo
		}

		var config struct {
			AccessKey string   `fig:"access_key,required"`
			URL       *url.URL `fig:"api_url,required"`
		}

		err = figure.
			Out(&config).
			With(figure.BaseHooks).
			From(configData).
			Please()
		if err != nil {
			panic(errors.Wrap(err, "failed to figure out geo_info"))
		}

		c.geoinfo = geoinfo.NewConnector(config.AccessKey, config.URL)
	}

	return c.geoinfo
}
