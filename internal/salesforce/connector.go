package salesforce

import (
	"net/url"
	"time"

	"github.com/pkg/errors"
)

const salesforceTimeLayout = "2006-01-02T15:04:05.999-0700"

// EmptyConnector is used for signalizing about special conditions
var EmptyConnector = &Connector{}

// Connector provides salesforce-interface to be used in PSIM services
type Connector struct {
	client *Client
}

// NewConnector construct a connector from arguments and gets accessToken
func NewConnector(authURL *url.URL, secret string, id string, username string, password string) (*Connector, error) {
	client := NewClient(authURL, secret, id)
	authResponse, err := client.PostAuthRequest(username, password)
	if err != nil {
		return nil, errors.Wrap(err, "failed to authenticate while constructing salesforce connector")
	}

	return &Connector{
		client: &Client{
			httpClient:  client.httpClient,
			authURL:     authURL,
			apiURL:      authResponse.InstanceURL,
			secret:      client.secret,
			accessToken: authResponse.AccessToken,
			id:          client.id,
		},
	}, nil
}

// SendEvent sends an event from arguments to salesforce
func (c *Connector) SendEvent(sphere string, actionName string, time time.Time, actorName string, actorEmail string, investmentAmount int64, investmentCountry string) (*EventResponse, error) {
	return c.client.PostEvent(sphere, actionName, time.Format(salesforceTimeLayout), actorName, actorEmail, investmentAmount, investmentCountry)
}
