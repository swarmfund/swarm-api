package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestViperConfig_API(t *testing.T) {
	raw := &mockRawGetter{}
	config := newViperConfig(raw)

	t.Run("correct config key", func(t *testing.T) {
		raw.On("GetStringMap", "api").Return(nil).Once()
		defer raw.AssertExpectations(t)
		assert.Panics(t, func() { config.API() })
	})
}
