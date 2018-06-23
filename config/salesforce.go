package config

import (
	"net/url"

	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/swarmfund/api/internal/salesforce"
)

// TODO use figure out

// Salesforce returns a ready-to-use salesforce connector
func (c *ViperConfig) Salesforce() *salesforce.Connector {
	c.Lock()
	defer c.Unlock()

	if c.salesforce == nil {
		v := c.GetStringMap("salesforce")
		if v == nil {
			panic("salesforce config entry is missing")
		}

		apiRawURL := v["api_url"].(string)
		apiURL, err := url.Parse(apiRawURL)
		if err != nil {
			panic(errors.Wrap(err, "failed to parse salesforce api url", logan.F{
				"api_url": apiRawURL,
			}))
		}
		secret := v["client_secret"].(string)
		id := v["client_id"].(string)
		username := v["username"].(string)
		password := v["password"].(string)

		salesforce, err := salesforce.NewConnector(apiURL, secret, id, username, password)
		if err != nil {
			panic(errors.Wrap(err, "failed to create connector"))
		}
		c.salesforce = salesforce
	}

	return c.salesforce
}
