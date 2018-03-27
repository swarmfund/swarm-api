package config

import (
	"gitlab.com/distributed_lab/figure"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

const metricKey = "metrics"

type Metric struct {
	Host string `fig:"host"`
	Port int    `fig:"port"`
}

func (c *ViperConfig) Metric() Metric {
	c.Lock()
	defer c.Unlock()

	if c.metric != nil {
		return *c.metric
	}

	metric := &Metric{}
	config := c.GetStringMap(metricKey)

	if err := figure.Out(metric).From(config).Please(); err != nil {
		panic(errors.Wrap(err, "failed to figure out metrics"))
	}

	c.metric = metric

	return *c.metric
}
