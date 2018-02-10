package config

import (
	"fmt"
	"reflect"

	"time"

	"github.com/evalphobia/logrus_sentry"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cast"
	"gitlab.com/distributed_lab/figure"
	"gitlab.com/distributed_lab/logan/v3"
)

const (
	logConfigKey = "log"
)

var (
	logLevelHook = figure.Hooks{
		"map[string]string": func(value interface{}) (reflect.Value, error) {
			result, err := cast.ToStringMapStringE(value)
			if err != nil {
				return reflect.Value{}, errors.Wrap(err, "failed to parse map[string]string")
			}
			return reflect.ValueOf(result), nil
		},
		"logan.Level": func(value interface{}) (reflect.Value, error) {
			switch v := value.(type) {
			case string:
				lvl, err := logan.ParseLevel(v)
				if err != nil {
					return reflect.Value{}, errors.Wrap(err, "failed to parse log level")
				}
				return reflect.ValueOf(lvl), nil
			case nil:
				return reflect.ValueOf(nil), nil
			default:
				return reflect.Value{}, fmt.Errorf("unsupported conversion from %T", value)
			}
		},
	}
)

func (c *ViperConfig) Log() *logan.Entry {
	// get sentry before acquiring lock
	sentry := c.Sentry()

	// acquire lock for log init
	c.Lock()
	defer c.Unlock()

	if c.logan != nil {
		return c.logan
	}

	var config struct {
		Level            logan.Level
		QueryThreshold   time.Duration
		RequestThreshold time.Duration
	}

	err := figure.
		Out(&config).
		With(figure.BaseHooks, logLevelHook).
		From(c.Get(logConfigKey)).
		Please()
	if err != nil {
		panic(errors.Wrap(err, "failed to figure out log"))
	}

	entry := logan.New().Level(config.Level)

	if sentry != nil {
		// sentry error hook
		levels := []logrus.Level{
			logrus.ErrorLevel,
			logrus.FatalLevel,
			logrus.PanicLevel,
		}

		hook, err := logrus_sentry.NewWithClientSentryHook(sentry, levels)
		if err != nil {
			panic(errors.Wrap(err, "failed to init sentry hook"))
		}
		hook.Timeout = 1 * time.Second
		entry.AddLogrusHook(hook)
	}

	c.logan = entry
	return c.logan
}
