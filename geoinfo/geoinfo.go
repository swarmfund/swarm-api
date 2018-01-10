package geoinfo

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type LocationInfo struct {
	Ip          string  `json:"ip"`
	CountryCode string  `json:"country_code"`
	CountryName string  `json:"country_name"`
	RegionCode  string  `json:"region_code"`
	RegionName  string  `json:"region_name"`
	City        string  `json:"city"`
	ZipCode     string  `json:"zip_code"`
	TimeZone    string  `json:"time_zone"`
	Latitudes   float64 `json:"latitudes"`
	Longitude   float64 `json:"longitude"`
	MetroCode   float64 `json:"metro_code"`
}

const serviceURL = "http://freegeoip.net/json/"

func GetLocationInfo(ip string) (*LocationInfo, error) {
	resp, err := http.Get(serviceURL + ip)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("request failed with status: %d", resp.StatusCode)
	}
	defer resp.Body.Close()

	var info LocationInfo
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&info)

	return &info, err
}

func (l *LocationInfo) FullRegion() string {
	sep := ", "
	locationParts := []string{l.City, l.RegionName, l.CountryName}
	location := strings.Join(locationParts, sep)

	// set location = "Unknown" if all locationParts is empty strings
	if len(location) == len(sep)*(len(locationParts)-1) {
		location = "Unknown"
	}

	// remove unnecessary whitespaces and separators when some locationParts empty
	return strings.Replace(location, " ,", "", -1)
}
