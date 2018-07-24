package config

import (
	"net/url"

	"gitlab.com/distributed_lab/figure"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/swarmfund/api/internal/salesforce"

)

// Salesforce returns a ready-to-use salesforce connector
// or nil if it's disabled by config
func (c *ViperConfig) Salesforce() *salesforce.Connector {
	c.Lock()
	defer c.Unlock()

	if c.salesforce == nil {
		var toggle struct {
			Disabled bool `fig:"disabled"`
		}
		var config struct {
			APIUrl       *url.URL `fig:"api_url,required"`
			ClientSecret string   `fig:"client_secret,required"`
			ClientID     string   `fig:"client_id,required"`
			Username     string   `fig:"username,required"`
			Password     string   `fig:"password,required"`
		}

		v := c.GetStringMap("salesforce")

		// getting toggle flag, to check if initialization is needed
		if err := figure.Out(&toggle).From(v).Please(); err != nil {
			panic(errors.Wrap(err, "failed to figure out salesforce toggle"))
		}
		if toggle.Disabled {
			// connector is disabled
			// FIXME find a better way w/o re-reading config
			return nil
		}

		if err := figure.Out(&config).From(v).Please(); err != nil {
			panic(errors.Wrap(err, "failed to figure out salesforce"))
		}

		salesforce, err := salesforce.NewConnector(
			config.APIUrl, config.ClientSecret, config.ClientID,
			config.Username, config.Password)
		if err != nil {
			panic(errors.Wrap(err, "failed to create connector"))
		}
		c.salesforce = salesforce
	}

	return c.salesforce
}
