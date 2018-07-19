package config

import (
	"fmt"

	"reflect"

	"github.com/minio/minio-go"
	"github.com/pkg/errors"
	"github.com/spf13/cast"
	"gitlab.com/distributed_lab/figure"
	"gitlab.com/swarmfund/api/storage"
)

const (
	storageConfigKey = "storage"
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

type Storage struct {
	AccessKey        string             `fig:"access_key"`
	SecretKey        string             `fig:"secret_key"`
	Host             string             `fig:"host"`
	ForceSSL         bool               `fig:"force_ssl"`
	MinContentLength int64              `fig:"min_content_length"`
	MaxContentLength int64              `fig:"max_content_length"`
	MediaTypes       storage.MediaTypes `fig:"media_types"`
}

func (c *ViperConfig) Storage() *storage.Connector {
	c.Lock()
	defer c.Unlock()

	if c.storage != nil {
		return c.storage
	}

	config := &Storage{}

	err := figure.Out(config).With(figure.BaseHooks, MediaTypeHook).From(c.GetStringMap(storageConfigKey)).Please()
	if err != nil {
		panic(errors.Wrap(err, "failed to figure out storage"))
	}

	minio, err := minio.New(config.Host, config.AccessKey, config.SecretKey, config.ForceSSL)
	if err != nil {
		panic(errors.Wrap(err, "failed to init client"))
	}

	connector := &storage.Connector{
		Minio:             minio,
		Log:               c.Log().WithField("service", "storage"),
		MinContentLength:  config.MinContentLength,
		MaxContentLength:  config.MaxContentLength,
		AllowedMediaTypes: config.MediaTypes,
	}

	c.storage = connector

	return c.storage
}
