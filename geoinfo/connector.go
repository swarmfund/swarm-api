package geoinfo

import (
	"encoding/json"
	"fmt"
	"net/http"

	"net/url"

	"gitlab.com/distributed_lab/geoinfo/resources"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

var (
	BaseURL         *url.URL
	ErrHTTPResponse = errors.New("request failed with status")
)

type Connector struct {
	accessKey string
	base      *url.URL
}

func (c *Connector) WithBase(url *url.URL) *Connector {
	return &Connector{
		accessKey: c.accessKey,
		base:      url,
	}
}

func NewConnector(accessKey string) *Connector {
	return &Connector{
		accessKey: accessKey,
		base:      BaseURL,
	}
}

func (c *Connector) LocationInfo(ip string) (*resources.LocationInfo, error) {
	endpoint := fmt.Sprintf("%s/%s?access_key=%s&output=json", c.base.String(), ip, c.accessKey)
	resp, err := http.Get(endpoint)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("failed to send request %s", endpoint))
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.Wrap(ErrHTTPResponse, resp.Status)
	}

	var info resources.LocationInfo
	if err := json.NewDecoder(resp.Body).Decode(&info); err != nil {
		return nil, errors.Wrap(err, "failed to decode location info")
	}

	return &info, nil
}

func init() {
	url, err := url.Parse("http://api.ipstack.com")
	if err != nil {
		panic(errors.Wrap(err, "failed to init base geoinfo url"))
	}

	BaseURL = url
}
