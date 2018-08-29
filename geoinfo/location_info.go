package geoinfo

import (
	"fmt"
)

type LocationInfo struct {
	IP            string  `json:"ip"`
	Type          string  `json:"type"`
	ContinentCode string  `json:"continent_code"`
	ContinentName string  `json:"continent_name"`
	CountryCode   string  `json:"country_code"`
	CountryName   string  `json:"country_name"`
	RegionName    string  `json:"region_name"`
	City          string  `json:"city"`
	Zip           string  `json:"zip"`
	Latitude      float64 `json:"latitude"`
	Longitude     float64 `json:"longitude"`
}

func (l *LocationInfo) FullRegion() string {
	if (l.City == "") && (l.RegionName == "") && (l.ContinentName == "") {
		return "Unknown"
	}

	location := addNotEmptyString(l.City, l.RegionName)
	return addNotEmptyString(location, l.CountryName)
}

func addNotEmptyString(basic, toAdd string) string {
	if toAdd == "" {
		return basic
	}

	if basic == "" {
		return toAdd
	}

	return fmt.Sprintf("%s, %s", basic, toAdd)
}
