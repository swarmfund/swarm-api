package config

import (
	"gitlab.com/distributed_lab/figure"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

const (
	txwatcherConfigKey = "tx_watcher"
)

var (
	txwatcher *TxWatcher
)

type TxWatcher struct {
	Disabled bool
}

func (c *ViperConfig) TxWatcher() TxWatcher {
	if txwatcher == nil {
		txwatcher = &TxWatcher{}
		config := c.GetStringMap(txwatcherConfigKey)
		if err := figure.Out(txwatcher).From(config).Please(); err != nil {
			panic(errors.Wrap(err, "failed to figure out txwatcher"))
		}
	}
	return *txwatcher
}
