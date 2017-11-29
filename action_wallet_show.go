package api

import (
	"gitlab.com/swarmfund/api/db2/api"
	"gitlab.com/swarmfund/api/render/hal"
	"gitlab.com/swarmfund/api/render/problem"
	"gitlab.com/swarmfund/api/resource"
	"gitlab.com/swarmfund/go/xdr"
)

type ShowWalletAction struct {
	Action
	DeviceInfo  *api.DeviceInfo
	Username    string
	WalletId    string
	fingerprint string
	wallet      *api.Wallet
}

func (action *ShowWalletAction) JSON() {
	action.Do(
		action.ValidateBodyType,
		action.loadParams,
		action.loadWallet,
		action.checkBlockState,
		action.checkTFA,
		action.checkDevice,
		func() {
			var response resource.Wallet
			response.Populate(action.wallet)
			hal.Render(action.W, response)
		})
}

func (action *ShowWalletAction) loadParams() {
	action.Username = action.GetNonEmptyString("username")
	action.WalletId = action.GetNonEmptyString("walletId")
}

func (action *ShowWalletAction) loadWallet() {
	wallet, err := action.APIQ().Wallet().ByEmail(action.Username)
	if err != nil {
		action.Log.WithError(err).Error("failed to get wallet")
		action.Err = &problem.ServerError
		return
	}

	if wallet == nil || action.WalletId != wallet.WalletId {
		action.Err = &problem.NotFound
		return
	}

	if !wallet.Verified {
		// email should be verified before login
		action.Err = &problem.Forbidden
		return
	}

	action.wallet = wallet
}

func (action *ShowWalletAction) checkBlockState() {
	account, err := action.App.horizon.AccountSigned(action.App.AccountManagerKP(), action.wallet.AccountID)
	if err != nil {
		action.Log.WithError(err).Error("failed to load account")
		action.Err = &problem.ServerError
		return
	}

	if account == nil {
		return
	}

	if account.BlockReasons&xdr.BlockReasonsRecoveryRequest != 0 {
		action.Err = &problem.ForbiddenReasonRecovery
		return
	}
}

func (action *ShowWalletAction) checkTFA() {
	action.consumeTFA(action.wallet, api.TFAActionLogin)
}

func (action *ShowWalletAction) checkDevice() {
	var err error

	action.DeviceInfo, err = action.GetSenderDeviceInfo(action.Username, action.App.Config().API().ClientDomain)
	if err != nil {
		action.Log.WithError(err).Error("Unable to get sender device info")
		action.Err = &problem.ServerError
		return
	}

	action.fingerprint, err = action.DeviceInfo.Fingerprint()
	if err != nil {
		action.Log.WithError(err).Error("Unable to get fingerprint of the device")
		action.Err = &problem.ServerError
		return
	}

	device, err := action.APIQ().AuthorizedDevice().ByFingerprint(action.fingerprint)
	if err != nil {
		action.Log.WithError(err).Error("Unable to get authorized device by fingerprint")
		action.Err = &problem.ServerError
		return
	}
	if device == nil {
		action.resolveNewDevice()
		return
	}

	err = action.APIQ().AuthorizedDevice().UpdateLastLoginTime(device)
	if err != nil {
		action.Log.WithError(err).Error("Failed to update last login time")
		action.Err = &problem.ServerError
		return
	}
}

func (action *ShowWalletAction) resolveNewDevice() {
	device := api.AuthorizedDevice{
		WalletID:    action.wallet.Id,
		Details:     *action.DeviceInfo,
		Fingerprint: action.fingerprint,
	}
	err := action.APIQ().AuthorizedDevice().Create(&device)
	if err != nil {
		action.Log.WithError(err).Error("Failed to add authorized device")
		action.Err = &problem.ServerError
		return
	}

	err = action.Notificator().SendNewDeviceLogin(action.wallet.Username, device)
	if err != nil {
		action.Log.WithError(err).Error("New device login emails sending failed")
	}
}
