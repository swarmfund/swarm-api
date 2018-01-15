package api

import (
	"net/http"

	"github.com/pkg/errors"
	"gitlab.com/swarmfund/api/render/problem"
	"gitlab.com/swarmfund/go/keypair"
	"gitlab.com/swarmfund/go/signcontrol"
	"gitlab.com/swarmfund/go/xdr"
	"gitlab.com/swarmfund/horizon-connector"
)

var (
	ErrNotAllowed = errors.New("not allowed")
)

type AccountGetter func(address string) (*horizon.Account, error)
type SignerConstraint func(app *App, r *http.Request, getter AccountGetter, signer string) error
type AccountSignerConstraint func(address string, app *App, r *http.Request, getter AccountGetter, signer string) error

// SignedBy is only useful for wallet-bound operations, where only entity involved
// is wallet itself.
func SignedBy(address string) SignerConstraint {
	return func(app *App, r *http.Request, getter AccountGetter, signer string) error {
		if r.Header.Get(signcontrol.SignatureHeader) == "" {
			return ErrNotAllowed
		}
		if signer == address {
			return nil
		}
		return ErrNotAllowed
	}
}

func SignerOf(address string, extra ...AccountSignerConstraint) SignerConstraint {
	return func(app *App, r *http.Request, getter AccountGetter, signer string) error {
		if r.Header.Get(signcontrol.SignatureHeader) == "" {
			return ErrNotAllowed
		}
		if signer == address && len(extra) == 0 {
			return nil
		}

		account, err := getter(address)
		if err != nil {
			return err
		}
		if account == nil {
			return ErrNotAllowed
		}
		// TODO make it readable
		for _, accountSigner := range account.Signers {
			if accountSigner.AccountID == signer && accountSigner.Weight > 0 {
				for _, constraint := range extra {
					if err := constraint(address, app, r, getter, signer); err != nil {
						return err
					}
				}
				return nil
			}
		}
		return ErrNotAllowed
	}
}

func SignerType(address string, signerType xdr.SignerType) SignerConstraint {
	return func(_ *App, _ *http.Request, getter AccountGetter, signer string) error {
		account, err := getter(address)
		if err != nil {
			return err
		}
		if account == nil {
			return ErrNotAllowed
		}
		for _, accountSigner := range account.Signers {
			if accountSigner.AccountID == signer &&
				accountSigner.Weight > 0 &&
				accountSigner.SignerType&int32(signerType) > 0 {
				return nil
			}
		}
		return ErrNotAllowed
	}
}

func (action *Action) checkSignerConstraints(constraints ...SignerConstraint) {
	//if action.App.config.SkipCheck {
	//	return
	//}

	accounts := map[string]*horizon.Account{}
	signer := action.R.Header.Get(signcontrol.PublicKeyHeader)
	getter := func(address string) (*horizon.Account, error) {
		var err error
		account, ok := accounts[address]
		if !ok {
			// FIXME move connector to proper keypair
			account, err = action.App.horizon.AccountSigned(keypair.MustParse(action.App.MasterSignerKP().Seed()), address)
			if err != nil {
				return nil, err
			}
			accounts[address] = account
		}
		return account, nil
	}

	for _, constraint := range constraints {
		switch err := constraint(action.App, action.R, getter, signer); err {
		case nil:
			return
		case ErrNotAllowed:
			continue
		default:
			action.Log.WithError(err).Error("failed to check constraints")
			action.Err = &problem.ServerError
			return
		}
	}

	action.Err = &problem.NotAllowed
}
