package config

import (
	"fmt"

	"reflect"

	"github.com/pkg/errors"
	"github.com/spf13/cast"
	"gitlab.com/distributed_lab/figure"
	"gitlab.com/swarmfund/api/internal/data"
	"gitlab.com/swarmfund/api/internal/data/s3storage"
	"gitlab.com/swarmfund/api/internal/storage"
)

var MediaTypeHook = figure.Hooks{
	"storage.MediaTypes": func(value interface{}) (reflect.Value, error) {
		switch v := value.(type) {
		case map[string]interface{}:
			temp := map[string][]string{}
			for docName, val := range v {
				mediaTypes, err := cast.ToStringSliceE(val)
				if err != nil {
					return reflect.Value{}, errors.New("failed to cast")
				}

				temp[docName] = mediaTypes
			}

			result, err := storage.NewMediaTypes(temp)
			if err != nil {
				return reflect.Value{}, errors.Wrap(err, "Failed to create new media types")
			}

			return reflect.ValueOf(result), nil
		default:
			return reflect.Value{}, errors.New(fmt.Sprintf("unsupported conversion from %T", value))
		}
	},
}

func (c *ViperConfig) Storage() data.Storage {
	raw := c.GetStringMap("storage")

	// check if storage is enabled
	var disabled struct {
		Disabled bool `fig:"disabled"`
	}
	if err := figure.Out(&disabled).From(raw).Please(); err != nil {
		panic(errors.Wrap(err, "failed to figure out storage disabled"))
	}
	if disabled.Disabled {
		// FIXME nil will cause consumers to panic
		return nil
	}

	// before acquiring lock we must get backend instance to avoid deadlock
	var probe struct {
		Backend string `fig:"backend,required"`
	}

	if err := figure.Out(&probe).From(raw).Please(); err != nil {
		panic(errors.Wrap(err, "failed to figure out storage backend"))
	}

	switch probe.Backend {
	case "aws":
		aws := c.AWS()
		// now it should be safe to acquire lock
		c.Lock()
		defer c.Unlock()

		var config struct {
			Bucket string `fig:"bucket,required"`
			//MediaTypes       storage.MediaTypes `fig:"media_types"`
		}

		if err := figure.Out(&config).From(raw).Please(); err != nil {
			panic(errors.Wrap(err, "failed to figure out storage"))
		}
		storage, err := s3storage.NewStorage(aws, config.Bucket)
		if err != nil {
			panic(errors.Wrap(err, "failed to init storage"))
		}
		c.storage = storage
	default:
		panic(fmt.Errorf("unknown backend type: %s", probe.Backend))
	}

	return c.storage
}
