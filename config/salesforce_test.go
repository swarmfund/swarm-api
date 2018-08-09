package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestViperConfig_Salesforce(t *testing.T) {
	t.Run("disabled", func(t *testing.T) {
		raw := &mockRawGetter{}
		raw.On("GetStringMap", "salesforce").Return(map[string]interface{}{
			"disabled": true,
		}).Once()
		defer raw.AssertExpectations(t)

		config := newViperConfig(raw)

		assert.Nil(t, config.Salesforce())
	})
}
