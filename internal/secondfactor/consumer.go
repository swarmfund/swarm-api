package secondfactor

import (
	"net/http"

	"github.com/pkg/errors"
	"gitlab.com/swarmfund/api/db2/api"
	"gitlab.com/swarmfund/api/internal/types"
	"gitlab.com/swarmfund/api/tfa"
)

type Consumer struct {
	TFAQ api.TFAQI

	// forced* fields does not have priority,
	// error will be returned in case of conflict
	forcedBackendType *types.WalletFactorType
	forcedBackend     tfa.Backend
}

// NewConsumer returns new instance of Consumer with minimal required parameters,
// further configuration is possible using With* methods
func NewConsumer(tfaQ api.TFAQI) *Consumer {
	return &Consumer{
		TFAQ: tfaQ,
	}
}

// WithBackendType forces consumer to use only backends of provided type
func (c *Consumer) WithBackendType(tpe types.WalletFactorType) *Consumer {
	return &Consumer{
		TFAQ:              c.TFAQ,
		forcedBackendType: &tpe,
	}
}

// WithBackend forces consumer to use only provided backend entity
func (c *Consumer) WithBackend(backend tfa.Backend) *Consumer {
	return &Consumer{
		TFAQ:          c.TFAQ,
		forcedBackend: backend,
	}
}

func (c *Consumer) backends(wallet *api.Wallet) ([]tfa.Backend, error) {
	if c.forcedBackend != nil {
		return []tfa.Backend{c.forcedBackend}, nil
	}

	records, err := c.TFAQ.Backends(wallet.WalletId)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get backends")
	}
	backends := make([]tfa.Backend, 0, len(records))
	for _, record := range records {
		if record.Priority <= 0 {
			continue
		}
		if c.forcedBackendType != nil && record.BackendType != *c.forcedBackendType {
			continue
		}
		backend, err := record.Backend()
		if err != nil {
			return nil, errors.Wrap(err, "failed to init backend")
		}
		backends = append(backends, backend)
	}
	return backends, nil
}

func (c *Consumer) Consume(r *http.Request, wallet *api.Wallet) error {
	token := RequestHash(r)

	// try to consume tfa token
	ok, err := c.TFAQ.Consume(token)
	if err != nil {
		return errors.Wrap(err, "failed to consume tfa")
	}

	if ok {
		// token was already verified and now consumed
		return nil
	}

	// get active wallet tfa backends
	backends, err := c.backends(wallet)
	if err != nil {
		return errors.Wrap(err, "failed to get backends")
	}

	if len(backends) == 0 {
		// no backends are enabled, wallet is not tfa protected
		return nil
	}

	// check if there is active token already
	otp, err := c.TFAQ.Get(token)
	if err != nil {
		return errors.Wrap(err, "failed to get tfa")
	}

	if otp == nil {
		// no active tfa, let's go through backends and try create new one
		for _, backend := range backends {
			otp = &api.TFA{
				BackendID: backend.ID(),
				Token:     token,
			}
			break
		}

		if err = c.TFAQ.Create(otp); err != nil {
			return errors.Wrap(err, "failed to store otp")
		}
	}
	meta := map[string]interface{}{}
	for _, backend := range backends {
		if backend.ID() == otp.BackendID {
			if metable, ok := backend.(tfa.HasMeta); ok {
				meta = metable.Meta()
			}
			break
		}
	}
	meta["token"] = token
	meta["factor_id"] = otp.BackendID
	meta["wallet_id"] = wallet.WalletId

	return &FactorRequiredErr{
		token: token,
		meta:  meta,
	}
}
