package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/distributed_lab/figure"
	"gitlab.com/swarmfund/api/internal/storage"
)

func TestViperConfig_Storage(t *testing.T) {
	t.Run("disabled", func(t *testing.T) {
		raw := &mockRawGetter{}
		raw.On("GetStringMap", "storage").Return(map[string]interface{}{
			"disabled": true,
		}).Once()
		defer raw.AssertExpectations(t)

		config := newViperConfig(raw)

		assert.Nil(t, config.Storage())
	})
}

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
