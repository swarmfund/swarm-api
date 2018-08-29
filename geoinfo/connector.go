package geoinfo

import (
	"encoding/json"
	"fmt"
	"net/http"

	"net/url"

	"gitlab.com/distributed_lab/logan/v3/errors"
)

type Connector struct {
	accessKey string
	base      *url.URL
}

func NewConnector(accessKey string, url *url.URL) *Connector {
	return &Connector{
		accessKey: accessKey,
		base:      url,
	}
}

func (c *Connector) LocationInfo(ip string) (*LocationInfo, error) {
	endpoint := fmt.Sprintf("%s/%s?access_key=%s&output=json", c.base.String(), ip, c.accessKey)
	resp, err := http.Get(endpoint)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("failed to send request %s", endpoint))
	}

	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, errors.Wrap(errors.New("request failed with status"), resp.Status)
	}

	var info LocationInfo
	if err := json.NewDecoder(resp.Body).Decode(&info); err != nil {
		return nil, errors.Wrap(err, "failed to decode location info")
	}

	return &info, nil
}
