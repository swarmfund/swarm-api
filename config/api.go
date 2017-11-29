package config

import (
	"net/url"
	"reflect"

	"fmt"

	"github.com/pkg/errors"
	"github.com/spf13/cast"
	"gitlab.com/distributed_lab/figure"
	"gitlab.com/swarmfund/go/keypair"
)

const (
	apiConfigKey = "api"
)

var (
	apiConfig *API
	// TODO move to figure
	URLHook = figure.Hooks{
		"url.URL": func(value interface{}) (reflect.Value, error) {
			str, err := cast.ToStringE(value)
			if err != nil {
				return reflect.Value{}, errors.Wrap(err, "failed to parse string")
			}
			u, err := url.Parse(str)
			if err != nil {
				return reflect.Value{}, errors.Wrap(err, "failed to parse url")
			}
			return reflect.ValueOf(*u), nil
		},
	}
	KeypairHook = figure.Hooks{
		"keypair.KP": func(value interface{}) (reflect.Value, error) {
			switch v := value.(type) {
			case string:
				kp, err := keypair.Parse(v)
				if err != nil {
					return reflect.Value{}, errors.Wrap(err, "failed to parse kp")
				}
				return reflect.ValueOf(kp), nil
			case nil:
				return reflect.ValueOf(nil), nil
			default:
				return reflect.Value{}, fmt.Errorf("unsupported conversion from %T", value)
			}
		},
	}
)

type API struct {
	DatabaseURL        string
	HorizonURL         url.URL
	AccountManager     keypair.KP
	SkipSignatureCheck bool

	NoEmailVerify bool
	ClientDomain  string
}

func (c *ViperConfig) API() API {
	if apiConfig == nil {
		apiConfig = &API{}
		config := c.GetStringMap(apiConfigKey)
		if err := figure.Out(apiConfig).With(figure.BaseHooks, URLHook, KeypairHook).From(config).Please(); err != nil {
			panic(errors.Wrap(err, "failed to figure out api"))
		}
	}
	return *apiConfig
}
