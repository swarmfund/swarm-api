package horizon

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/pkg/errors"
)

type Signers struct {
	Signers []Signer `json:"signers"`
}

type SignerType struct {
	Name  string `json:"name"`
	Value int32  `json:"value"`
}

// Signer represents one of an account's signers.
type Signer struct {
	PublicKey      string       `json:"public_key"`
	Weight         int32        `json:"weight"`
	SignerTypeI    int32        `json:"signer_type_i"`
	SignerTypes    []SignerType `json:"signer_types"`
	SignerIdentity int32        `json:"signer_identity"`
	SignerName     string       `json:"signer_name"`
}

var (
	endpointAccountSigners = func(address string) string {
		return fmt.Sprintf("/accounts/%s/signers", address)
	}
)

type Client struct {
	base   url.URL
	client *http.Client
}

func New(endpoint url.URL) *Client {
	return &Client{
		base:   endpoint,
		client: http.DefaultClient,
	}
}

func (c *Client) NewRequest(method, endpoint string, body io.Reader) (*http.Request, error) {
	u, err := url.Parse(endpoint)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse url")
	}
	request, err := http.NewRequest("GET", c.base.ResolveReference(u).String(), body)
	request.Header.Set("Content-Type", "application/vnd.api+json")
	return request, nil
}

func (c *Client) Account(address string) *AccountBuilder {
	return &AccountBuilder{
		c, address,
	}
}

type AccountBuilder struct {
	client  *Client
	address string
}

func (b *AccountBuilder) Signers() ([]Signer, error) {
	request, err := b.client.NewRequest("GET", endpointAccountSigners(b.address), nil)
	if err != nil {
		return nil, errors.Wrap(err, "failed to build request")
	}
	response, err := b.client.client.Do(request)
	if err != nil {
		return nil, errors.Wrap(err, "request failed")
	}
	defer response.Body.Close()

	if response.StatusCode != 200 {
		return nil, fmt.Errorf("request failed with %s", response.Status)
	}

	var result Signers
	if err := json.NewDecoder(response.Body).Decode(&result); err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal body")
	}

	return result.Signers, nil
}
