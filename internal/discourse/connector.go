package discourse

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/pkg/errors"
	"gitlab.com/swarmfund/api/internal/lorem"
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
	// TODO fix username limitation
	//name := strings.Split(r.Email, "@")[0]
	name := lorem.Token()
	if r.Name == "" {
		r.Name = name
	}
	if r.Username == "" {
		r.Username = name
	}
	if r.Password == "" {
		r.Password = "supersekritp@ssw0rt"
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

type apiResponse struct {
	/*
		{"success":false,"message":"Username must be no more than 20 characters\nPassword is too short (minimum is 10 characters)","errors":{"username":["must be no more than 20 characters"],"password":["is too short (minimum is 10 characters)"]},"values":{"name":"yr0a3ke29d78skm2030kep14i","username":"yr0a3ke29d78skm2030kep14i","email":"yr0a3ke29d78skm2030kep14i@test.com"},"is_developer":false}
	*/
	Success bool   `json:"success"`
	Message string `json:"message"`
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

	var resp apiResponse
	if err := json.NewDecoder(response.Body).Decode(&resp); err != nil {
		return errors.Wrap(err, "failed to unmarshal body")
	}

	if !resp.Success {
		return errors.Wrap(errors.New("request failed"), resp.Message)
	}

	return nil
}
