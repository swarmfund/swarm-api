package config

import (
	"gitlab.com/distributed_lab/figure"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

const walletsConfigkey = "wallets"

type Wallets struct {
	DisableConfirm   bool     `fig:"disable_confirm"`
	DomainsBlacklist []string `fig:"domains_blacklist"`
}

func (c *ViperConfig) Wallets() Wallets {
	c.Lock()
	defer c.Unlock()

	if c.wallets != nil {
		return *c.wallets
	}

	wallets := &Wallets{}
	config := c.GetStringMap(walletsConfigkey)

	if err := figure.Out(wallets).From(config).Please(); err != nil {
		panic(errors.Wrap(err, "failed to figure out wallets"))
	}

	c.wallets = wallets

	return *c.wallets
}
