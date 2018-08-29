package geoinfo

import (
	"testing"

	"github.com/magiconair/properties/assert"
)

func TestLocationInfo_FullRegion(t *testing.T) {
	t.Run("Get full region name", func(t *testing.T) {
		locationInfo := LocationInfo{
			CountryName: "Ukraine",
			RegionName:  "Kharkivs'ka Oblast'",
			City:        "Kharkiv",
		}

		expected := "Kharkiv, Kharkivs'ka Oblast', Ukraine"

		got := locationInfo.FullRegion()
		assert.Equal(t, expected, got)
	})

	t.Run("If fields is empty", func(t *testing.T) {
		locationInfo := LocationInfo{}

		got := locationInfo.FullRegion()
		assert.Equal(t, "Unknown", got)
	})
}
