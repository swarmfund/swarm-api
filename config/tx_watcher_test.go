package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestViperConfig_TXWatcher(t *testing.T) {
	t.Run("disabled", func(t *testing.T) {
		raw := &mockRawGetter{}
		raw.On("GetStringMap", "tx_watcher").Return(map[string]interface{}{
			"disabled": true,
		}).Once()
		defer raw.AssertExpectations(t)

		config := newViperConfig(raw)

		assert.True(t, config.TXWatcher().Disabled)
	})
}
