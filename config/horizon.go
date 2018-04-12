package config

import (
	"net/url"

	"github.com/pkg/errors"
	"gitlab.com/distributed_lab/figure"
	"gitlab.com/swarmfund/horizon-connector/v2"
	"gitlab.com/tokend/keypair"
)

const (
	horizonConfigKey = "horizon"
)

func (c *ViperConfig) Horizon() *horizon.Connector {
	c.Lock()
	defer c.Unlock()

	if c.horizon == nil {
		var config struct {
			URL    url.URL
			Signer keypair.Full
		}

		err := figure.
			Out(&config).
			With(URLHook, KeypairHook).
			From(c.GetStringMap(horizonConfigKey)).
			Please()
		if err != nil {
			panic(errors.Wrap(err, "failed to figure out horizon"))
		}

		c.horizon = horizon.NewConnector(&config.URL).WithSigner(config.Signer)

		accountQ := c.horizon.Accounts()
		info, err := c.horizon.Info()
		if err != nil {
			panic(errors.Wrap(err, "Failed to get horizon info"))
		}

		//if err := accountQ.IsSigner(info.MasterAccountID, config.Signer.Address(), xdr.SignerTypeNotVerifiedAccManager); err != nil {
		//	panic(errors.Wrap(err, "Check signer failed"))
		//}
	}

	return c.horizon
}
