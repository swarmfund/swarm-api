package config

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"gitlab.com/distributed_lab/figure"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

// AWS is not a part of public interface and used to initialize other service connectors
func (c *ViperConfig) AWS() *session.Session {
	c.Lock()
	defer c.Unlock()

	if c.aws != nil {
		return c.aws
	}

	raw := c.GetStringMap("aws")

	// first probe for credentials type we should use
	var probe struct {
		Credentials string `fig:"credentials,required"`
	}

	if err := figure.Out(&probe).From(raw).Please(); err != nil {
		panic(errors.Wrap(err, "failed to figure out credentials type"))
	}

	// now figure out config depending on credentials type
	switch probe.Credentials {
	case "static":
		var config struct {
			AccessKey string `fig:"access_key,required"`
			SecretKey string `fig:"secret_key,required"`
			Endpoint  string `fig:"endpoint,required"`
			Region    string `fig:"region,required"`
		}

		if err := figure.Out(&config).From(raw).Please(); err != nil {
			panic(errors.Wrap(err, "failed to figure out static config"))
		}

		cfg := &aws.Config{
			Credentials: credentials.NewStaticCredentials(config.AccessKey, config.SecretKey, ""),
			Region:      aws.String(config.Region),
		}

		session, err := session.NewSession(cfg)
		if err != nil {
			panic(errors.Wrap(err, "failed to AWS establish session"))
		}

		c.aws = session
	case "ec2":
		panic("not implemented")
	default:
		panic(fmt.Errorf("unknown credentials type: %s", probe.Credentials))
	}

	return c.aws
}
