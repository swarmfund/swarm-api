package config

import (
	"fmt"
	"reflect"

	"github.com/pkg/errors"
	"gitlab.com/distributed_lab/figure"
	"gitlab.com/distributed_lab/logan/v3"
)

const (
	logConfigKey = "log"
)

var (
	entry        *logan.Entry
	logLevelHook = figure.Hooks{
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
	if entry != nil {
		return entry
	}

	var config struct {
		Level logan.Level
	}
	err := figure.
		Out(&config).
		With(logLevelHook).
		From(c.GetStringMap(logConfigKey)).
		Please()
	if err != nil {
		panic(errors.Wrap(err, "failed to figure out log"))
	}

	entry = logan.New().Level(config.Level)

	return entry
}
