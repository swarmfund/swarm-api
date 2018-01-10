package discourse

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/pkg/errors"
)

type Connector struct {
	endpoint *url.URL
	username string
	key      string
}

func NewConnector(endpoint *url.URL, username, key string) *Connector {
	return &Connector{
		endpoint: endpoint,
		username: username,
		key:      key,
	}
}

type AuthenticatedRequest struct {
	APIUsername string `json:"api_username"`
	APIKey      string `json:"api_key"`
}

type CreateUser struct {
	AuthenticatedRequest
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Username string `json:"username"`
}

func (r *CreateUser) Prepare(c *Connector) {
	name := strings.Split(r.Email, "@")[0]
	if r.Name == "" {
		r.Name = name
	}
	if r.Username == "" {
		r.Username = name
	}
	if r.Password == "" {
		r.Password = "sekrit"
	}
	if r.APIKey == "" {
		r.APIKey = c.key
	}
	if r.APIUsername == "" {
		r.APIUsername = c.username
	}
}

func (r *CreateUser) Validate() error {
	// TODO implement
	return nil
}

func (c *Connector) CreateUser(opts CreateUser) error {
	opts.Prepare(c)

	if err := opts.Validate(); err != nil {
		return errors.Wrap(err, "options are not valid")
	}

	bytes, err := json.Marshal(&opts)
	if err != nil {
		return errors.Wrap(err, "failed to marshal options")
	}

	var dict map[string]string
	if err := json.Unmarshal(bytes, &dict); err != nil {
		return errors.Wrap(err, "failed to unmarshal options")
	}

	form := url.Values{}
	for key, value := range dict {
		form.Add(key, value)
	}

	response, err := http.Post(
		// TODO proper url build
		fmt.Sprintf("%s/%s", c.endpoint.String(), "/users"),
		"application/x-www-form-urlencoded",
		strings.NewReader(form.Encode()),
	)
	if err != nil {
		return errors.Wrap(err, "request failed")
	}
	defer response.Body.Close()

	if response.StatusCode != 200 {
		body, err := ioutil.ReadAll(response.Body)
		if err != nil {
			return errors.Wrap(err, "failed to read body")
		}
		return errors.Wrap(errors.New("request failed"), string(body))
	}

	// TODO unmarshal json and check success code

	return nil
}
