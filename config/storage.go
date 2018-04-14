package config

import (
	"github.com/minio/minio-go"
	"github.com/pkg/errors"
	"gitlab.com/distributed_lab/figure"
	"gitlab.com/swarmfund/api/log"
	"gitlab.com/swarmfund/api/storage"
)

const (
	storageConfigKey = "storage"
)

type Storage struct {
	AccessKey        string   `fig:"access_key"`
	SecretKey        string   `fig:"secret_key"`
	Host             string   `fig:"host"`
	ForceSSL         bool     `fig:"force_ssl"`
	MinContentLength int64    `fig:"min_content_length"`
	MaxContentLength int64    `fig:"max_content_length"`
	MediaTypes       []string `fig:"media_types"`
}

func (c *ViperConfig) Storage() *storage.Connector {
	c.Lock()
	defer c.Unlock()

	if c.storage != nil {
		return c.storage
	}

	config := &Storage{}

	err := figure.Out(config).From(c.GetStringMap(storageConfigKey)).Please()
	if err != nil {
		panic(errors.Wrap(err, "failed to figure out storage"))
	}

	storage.SetMediaTypes(config.MediaTypes)

	minio, err := minio.New(config.Host, config.AccessKey, config.SecretKey, config.ForceSSL)
	if err != nil {
		panic(errors.Wrap(err, "failed to init client"))
	}

	connector := &storage.Connector{
		Minio:            minio,
		Log:              log.WithField("service", "storage"),
		MinContentLength: config.MinContentLength,
		MaxContentLength: config.MaxContentLength,
	}

	c.storage = connector

	return c.storage
}
