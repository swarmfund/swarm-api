package api

import (
	"errors"

	"gitlab.com/swarmfund/api/db2/api"
	"gitlab.com/swarmfund/api/render/hal"
	"gitlab.com/swarmfund/api/render/problem"
	"gitlab.com/swarmfund/api/resource"
	"gitlab.com/swarmfund/go/xdr"
)

type GetTFAAction struct {
	Action

	walletID string
	username string
	wallet   *api.Wallet
	Resource resource.TFABackends
}

func (action *GetTFAAction) JSON() {
	action.Do(
		action.ValidateBodyType,
		action.loadParams,
		action.loadWallet,
		action.checkAllowed,
		action.performRequest,
		func() {
			hal.Render(action.W, action.Resource)
		},
	)
}

func (action *GetTFAAction) loadParams() {
	action.walletID = action.GetString("wallet_id")
	action.username = action.GetString("username")
}

func (action *GetTFAAction) loadWallet() {
	var err error
	var wallet = new(api.Wallet)

	if action.walletID != "" {
		wallet, err = action.APIQ().Wallet().ByWalletID(action.walletID)
	} else if action.username != "" {
		wallet, err = action.loadWalletViaRecovery()
	} else {
		action.SetInvalidField("wallet_id", errors.New("must not be empty"))
		return
	}

	if err != nil {
		action.Log.WithError(err).Error("failed to load wallet")
		action.Err = &problem.ServerError
		return
	}

	if wallet == nil {
		action.SetInvalidField("wallet_id", errors.New("does not exists"))
		return
	}

	action.wallet = wallet
}

func (action *GetTFAAction) loadWalletViaRecovery() (*api.Wallet, error) {
	recovery, err := action.APIQ().Recoveries().ByUsername(action.username)
	if err != nil {
		return nil, err
	}

	if recovery == nil {
		return action.APIQ().Wallet().ByEmail(action.username)
	}

	return action.APIQ().Wallet().ByID(recovery.WalletID)
}

func (action *GetTFAAction) checkAllowed() {
	action.checkSignerConstraints(
		SignerOf(action.wallet.AccountID),
	)
	if action.Err == &problem.NotAllowed {
		account, err := action.App.horizon.AccountSigned(action.App.AccountManagerKP(), action.wallet.AccountID)
		if err != nil {
			action.Log.WithError(err).Error("failed to get account")
			action.Err = &problem.ServerError
			return
		}
		if account == nil {
			return
		}
		if account.AccountType != xdr.AccountTypeNotVerified {
			return
		}
		action.Err = nil
		action.checkSignerConstraints(
			SignerType(action.App.CoreInfo.MasterAccountID, xdr.SignerTypeNotVerifiedAccManager),
		)
	}
}

func (action *GetTFAAction) performRequest() {
	backends, err := action.APIQ().TFA().Backends(action.wallet.WalletId)
	if err != nil {
		action.Log.WithError(err).Error("failed to load wallets")
		action.Err = &problem.ServerError
		return
	}

	action.Resource = resource.TFABackends{
		Backends: []resource.TFABackend{},
	}

	for _, backend := range backends {
		res := resource.TFABackend{}
		res.Populate(&backend)
		action.Resource.Backends = append(action.Resource.Backends, res)
	}
}
