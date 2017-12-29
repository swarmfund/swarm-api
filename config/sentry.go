package config

import (
	"time"

	"github.com/evalphobia/logrus_sentry"
	"github.com/getsentry/raven-go"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"gitlab.com/distributed_lab/figure"
	"gitlab.com/distributed_lab/logan/v3"
)

const (
	sentryConfigKey = "sentry"
)

var (
	sentry *raven.Client
)

func (c *ViperConfig) Sentry() *raven.Client {
	if sentry != nil {
		return sentry
	}

	var config struct {
		DSN   string
		Level logan.Level
	}

	err := figure.Out(&config).With(figure.BaseHooks, logLevelHook).From(c.Get(sentryConfigKey)).Please()
	if err != nil {
		panic(errors.Wrap(err, "failed to figure out sentry"))
	}

	client, err := raven.New(config.DSN)
	if err != nil {
		panic(errors.Wrap(err, "failed to init sentry client"))
	}

	var levels []logrus.Level
	for i := config.Level + 1; i > 0; i-- {
		levels = append(levels, logrus.Level(i-1))
	}

	hook, err := logrus_sentry.NewWithClientSentryHook(client, levels)
	if err != nil {
		panic(errors.Wrap(err, "failed to init sentry hook"))
	}
	hook.Timeout = 1 * time.Second

	c.Log().AddLogrusHook(hook)

	sentry = client
	return sentry

}
