package geoinfo_test

import (
	"testing"

	"net/http/httptest"

	"net/http"

	"net/url"

	"github.com/stretchr/testify/assert"
	"gitlab.com/swarmfund/api/geoinfo"
	"gitlab.com/swarmfund/api/geoinfo/resources"
)

func TestConnector(t *testing.T) {
	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		response := []byte(`
{
  "ip": "159.224.199.93",
  "type": "ipv4",
  "continent_code": "EU",
  "continent_name": "Europe",
  "country_code": "UA",
  "country_name": "Ukraine",
  "region_code": "63",
  "region_name": "Kharkivs'ka Oblast'",
  "city": "Kharkiv",
  "zip": "61024",
  "latitude": 49.9808,
  "longitude": 36.2527
}
`)
		w.Write(response)
	}))

	ts := httptest.NewServer(handler)
	defer ts.Close()

	url, err := url.Parse(ts.URL)
	if err != nil {
		t.Fatal(err)
	}

	connector := geoinfo.NewConnector("288fbbd921b99d9d60024c5a232d0546").WithBase(url)

	expected := resources.LocationInfo{
		Ip:            "159.224.199.93",
		Type:          "ipv4",
		ContinentCode: "EU",
		ContinentName: "Europe",
		CountryCode:   "UA",
		CountryName:   "Ukraine",
		RegionName:    "Kharkivs'ka Oblast'",
		City:          "Kharkiv",
		Zip:           "61024",
		Latitude:      49.9808,
		Longitude:     36.2527,
	}

	got, err := connector.LocationInfo("159.224.199.93")
	assert.NoError(t, err)
	assert.EqualValues(t, expected, *got)
}
