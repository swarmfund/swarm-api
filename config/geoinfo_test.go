package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestViperConfig_GeoInfo(t *testing.T) {
	raw := &mockRawGetter{}
	raw.On("GetStringMap", "geo_info").Return(map[string]interface{}{
		"api_url":    "http://api.ipstack.com",
		"access_key": "fb170f98697192973f33434cb35157b4",
	}).Once()
	defer raw.AssertExpectations(t)

	config := newViperConfig(raw)

	assert.NotNil(t, config.GeoInfo())
}

func TestViperConfig_GeoInfoPanics(t *testing.T) {
	raw := &mockRawGetter{}
	raw.On("GetStringMap", "geo_info").Return(map[string]interface{}{
		"api_url": "http://api.ipstack.com",
	}).Once()
	defer raw.AssertExpectations(t)

	config := newViperConfig(raw)

	assert.Panics(t, func() { config.GeoInfo() })

	raw = &mockRawGetter{}
	raw.On("GetStringMap", "geo_info").Return(map[string]interface{}{
		"access_key": "fb170f98697192973f33434cb35157b4",
	}).Once()

	config = newViperConfig(raw)

	assert.Panics(t, func() { config.GeoInfo() })
}
