package config

import (
	"github.com/getsentry/raven-go"
	"github.com/pkg/errors"
	"gitlab.com/distributed_lab/figure"
)

const (
	sentryConfigKey = "sentry"
)

func (c *ViperConfig) Sentry() *raven.Client {
	c.Lock()
	defer c.Unlock()

	if c.sentry != nil {
		return c.sentry
	}

	var config struct {
		Disabled bool
		DSN      string
		Tags     map[string]string
	}

	err := figure.
		Out(&config).
		With(figure.BaseHooks, logLevelHook).
		From(c.Get(sentryConfigKey)).
		Please()
	if err != nil {
		panic(errors.Wrap(err, "failed to figure out sentry"))
	}

	if config.Disabled {
		return nil
	}

	client, err := raven.New(config.DSN)
	if err != nil {
		panic(errors.Wrap(err, "failed to init sentry client"))
	}

	if config.Tags != nil {
		client.Tags = config.Tags
	}

	c.sentry = client
	return c.sentry
}
