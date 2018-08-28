package config

import (
	"gitlab.com/distributed_lab/figure"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

type TXWatcher struct {
	Disabled bool
}

func (c *ViperConfig) TXWatcher() TXWatcher {
	c.Lock()
	defer c.Unlock()

	var config TXWatcher

	if err := figure.Out(&config).From(c.GetStringMap("tx_watcher")).Please(); err != nil {
		panic(errors.Wrap(err, "failed to figure out tx_watcher"))
	}

	return config
}
