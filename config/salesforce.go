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
		var probe struct {
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

		// check if service is enabled
		if err := figure.Out(&probe).From(v).Please(); err != nil {
			panic(errors.Wrap(err, "failed to figure out salesforce probe"))
		}

		if probe.Disabled {
			// FIXME returning nil here, will cause handlers to panic
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
