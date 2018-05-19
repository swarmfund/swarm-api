package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/distributed_lab/figure"
	"gitlab.com/swarmfund/api/storage"
)

func TestViperConfig_StorageMediaTypesHook(t *testing.T) {
	var config struct {
		MediaTypes storage.MediaTypes `fig:"media_types"`
	}

	configData := map[string]interface{}{
		"media_types": map[string]interface{}{
			"general": []string{"image/jpeg", "image/tiff", "image/png", "image/gif"},
			"alpha":   []string{"application/pdf"},
		},
	}

	err := figure.Out(&config).With(figure.BaseHooks, MediaTypeHook).From(configData).Please()
	assert.NoError(t, err)

	expectedMediaTypes, err := storage.NewMediaTypes(
		map[string][]string{
			"general": {"image/jpeg", "image/tiff", "image/png", "image/gif"},
			"alpha":   {"application/pdf"},
		})
	assert.EqualValues(t, expectedMediaTypes, config.MediaTypes)
}
