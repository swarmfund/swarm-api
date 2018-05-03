package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestViperConfig_Wallets(t *testing.T) {
	walletsConfigRaw := `
wallets:
  disable_confirm: true
  domains_blacklist:
    - mailinator.com
`
	config := ConfigHelper(t, walletsConfigRaw)

	expected := Wallets{
		DisableConfirm:   true,
		DomainsBlacklist: []string{"mailinator.com"},
	}

	got := config.Wallets()
	assert.Equal(t, expected, got)
}
