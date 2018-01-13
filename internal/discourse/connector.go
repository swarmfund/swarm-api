package discourse

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"io"

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
	Active   bool   `json:"active"`
	Approved bool   `json:"approved"`
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
	requestBody, err := c.prepareBody(&opts)
	if err != nil {
		return errors.Wrap(err, "failed to prepare options")
	}

	response, err := http.Post(
		// TODO proper url build
		fmt.Sprintf("%s/%s", c.endpoint.String(), "/users"),
		"application/x-www-form-urlencoded",
		requestBody,
	)
	if err != nil {
		return errors.Wrap(err, "request failed")
	}
	defer response.Body.Close()

	return c.handleResponse(response)
}

type CreateCategory struct {
	AuthenticatedRequest
	Name      string `json:"name"`
	Color     string `json:"color"`
	TextColor string `json:"text_color"`
}

func (r *CreateCategory) Prepare(c *Connector) {
	// TODO random color
	color := "49d9e9"
	if r.Color == "" {
		r.Color = color
	}
	if r.TextColor == "" {
		r.TextColor = color
	}
	if r.APIKey == "" {
		r.APIKey = c.key
	}
	if r.APIUsername == "" {
		r.APIUsername = c.username
	}
}

func (r *CreateCategory) Validate() error {
	// TODO implement
	return nil
}
func (c *Connector) CreateCategory(opts CreateCategory) error {
	requestBody, err := c.prepareBody(&opts)
	if err != nil {
		return errors.Wrap(err, "failed to prepare options")
	}

	response, err := http.Post(
		// TODO proper url build
		fmt.Sprintf("%s/%s", c.endpoint.String(), "/categories.json"),
		"application/x-www-form-urlencoded",
		requestBody,
	)
	if err != nil {
		return errors.Wrap(err, "request failed")
	}
	defer response.Body.Close()

	return c.handleResponse(response)
}

type preparable interface {
	Prepare(*Connector)
}

type validatable interface {
	Validate() error
}

func (c *Connector) handleResponse(response *http.Response) error {
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

func (c *Connector) prepareBody(opts interface{}) (io.Reader, error) {
	if v, ok := opts.(preparable); ok {
		v.Prepare(c)
	}

	if v, ok := opts.(validatable); ok {
		if err := v.Validate(); err != nil {
			return nil, errors.Wrap(err, "options are not valid")
		}
	}

	bytes, err := json.Marshal(&opts)
	if err != nil {
		return nil, errors.Wrap(err, "failed to marshal options")
	}

	var dict map[string]string
	if err := json.Unmarshal(bytes, &dict); err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal options")
	}

	form := url.Values{}
	for key, value := range dict {
		form.Add(key, value)
	}

	return strings.NewReader(form.Encode()), nil
}
