package config

import (
	"fmt"

	"reflect"

	"github.com/minio/minio-go"
	"github.com/pkg/errors"
	"github.com/spf13/cast"
	"gitlab.com/distributed_lab/figure"
	"gitlab.com/swarmfund/api/internal/types"
	"gitlab.com/swarmfund/api/log"
	"gitlab.com/swarmfund/api/storage"
)

const (
	storageConfigKey = "storage"
)

var MediaTypeHook = figure.Hooks{
	"map[string]storage.MediaTypes": func(value interface{}) (reflect.Value, error) {
		switch v := value.(type) {
		case map[string]interface{}:
			result := map[string]storage.MediaTypes{}
			for docName, val := range v {
				mediaTypes, err := cast.ToStringSliceE(val)
				if err != nil {
					return reflect.Value{}, errors.New("failed to cast")
				}

				result[docName] = storage.NewMediaTypes(mediaTypes)
			}

			return reflect.ValueOf(result), nil
		default:
			return reflect.Value{}, errors.New(fmt.Sprintf("unsupported conversion from %T", value))
		}
	},
}

type Storage struct {
	AccessKey        string                        `fig:"access_key"`
	SecretKey        string                        `fig:"secret_key"`
	Host             string                        `fig:"host"`
	ForceSSL         bool                          `fig:"force_ssl"`
	MinContentLength int64                         `fig:"min_content_length"`
	MaxContentLength int64                         `fig:"max_content_length"`
	MediaTypes       map[string]storage.MediaTypes `fig:"media_types"`
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

	allowedMediaTypes := map[types.DocumentType]storage.MediaTypes{}
	for docName, mediaType := range config.MediaTypes {
		var docType types.DocumentType

		if err := docType.UnmarshalJSON([]byte(fmt.Sprintf(`"%s"`, docName))); err != nil {
			panic(errors.Wrap(err, "failed to get document type"))
		}

		allowedMediaTypes[docType] = mediaType
	}

	minio, err := minio.New(config.Host, config.AccessKey, config.SecretKey, config.ForceSSL)
	if err != nil {
		panic(errors.Wrap(err, "failed to init client"))
	}

	connector := &storage.Connector{
		Minio:             minio,
		Log:               log.WithField("service", "storage"),
		MinContentLength:  config.MinContentLength,
		MaxContentLength:  config.MaxContentLength,
		AllowedMediaTypes: allowedMediaTypes,
	}

	c.storage = connector

	return c.storage
}
