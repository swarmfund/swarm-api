package resources

import "strings"

type LocationInfo struct {
	Ip            string  `json:"ip"`
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
