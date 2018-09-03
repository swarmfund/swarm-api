package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestViperConfig_GeoInfo(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		raw := &mockRawGetter{}
		raw.On("GetStringMap", "geo_info").Return(map[string]interface{}{
			"api_url":    "http://api.ipstack.com",
			"access_key": "fb170f98697192973f33434cb35157b4",
		}).Once()
		defer raw.AssertExpectations(t)

		config := newViperConfig(raw)

		assert.NotNil(t, config.GeoInfo())
	})

	t.Run("panic because access key not specified", func(t *testing.T) {
		raw := &mockRawGetter{}
		raw.On("GetStringMap", "geo_info").Return(map[string]interface{}{
			"api_url": "http://api.ipstack.com",
		}).Once()
		defer raw.AssertExpectations(t)

		config := newViperConfig(raw)

		assert.Panics(t, func() { config.GeoInfo() })
	})

	t.Run("panic because api url not specified", func(t *testing.T) {
		raw := &mockRawGetter{}
		raw.On("GetStringMap", "geo_info").Return(map[string]interface{}{
			"access_key": "fb170f98697192973f33434cb35157b4",
		}).Once()
		defer raw.AssertExpectations(t)

		config := newViperConfig(raw)

		assert.Panics(t, func() { config.GeoInfo() })
	})

	t.Run("disabled geoConnector", func(t *testing.T) {
		raw := &mockRawGetter{}
		raw.On("GetStringMap", "geo_info").Return(map[string]interface{}{
			"access_key": "fb170f98697192973f33434cb35157b4",
			"api_url":    "http://api.ipstack.com",
			"disabled":   "true",
		}).Once()
		defer raw.AssertExpectations(t)

		config := newViperConfig(raw)
		geoConnector := config.GeoInfo()
		assert.NotNil(t, geoConnector)

		location, err := geoConnector.LocationInfo("127.0.0.1")
		assert.NotNil(t, location)
		assert.NoError(t, err)
		assert.Equal(t, location.FullRegion(), "Unknown")
	})
}
