package config

import (
	"github.com/pkg/errors"
	"gitlab.com/distributed_lab/figure"
)

// WalletCleaner is config for
type WalletCleaner struct {
	Enabled    bool
	Expiration string //time duration after than unconfirmed wallet will be deleted
}

const walletCleanerConfKey = "wallet_cleaner"

var walletCleaner *WalletCleaner

func (c *ViperConfig) WalletCleaner() WalletCleaner {
	if walletCleaner == nil {
		walletCleaner = &WalletCleaner{}
		config := c.GetStringMap(walletCleanerConfKey)
		if err := figure.Out(walletCleaner).From(config).Please(); err != nil {
			panic(errors.Wrap(err, "failed to figure out wallet_cleaner"))
		}
	}
	return *walletCleaner
}
