package config

import (
	"fmt"
	"reflect"

	"time"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"gitlab.com/distributed_lab/figure"
)

const (
	logConfigKey = "log"
)

var (
	logConfig    *Log
	logLevelHook = figure.Hooks{
		"logrus.Level": func(value interface{}) (reflect.Value, error) {
			switch v := value.(type) {
			case string:
				lvl, err := logrus.ParseLevel(v)
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

type Log struct {
	Level            logrus.Level   `fig:"level"`
	SlowQueryBound   *time.Duration `fig:"slow_query_bound"`
	SlowRequestBound *time.Duration `fig:"slow_request_bound"`
}

func (c *ViperConfig) Log() Log {
	if logConfig == nil {
		logConfig = &Log{}
		config := c.GetStringMap(logConfigKey)
		if err := figure.Out(logConfig).With(figure.BaseHooks, logLevelHook).From(config).Please(); err != nil {
			panic(errors.Wrap(err, "failed to figure out log"))
		}
		fmt.Printf("%+v", *logConfig)
	}
	return *logConfig
}
