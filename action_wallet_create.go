package api

import (
	"github.com/go-errors/errors"
	"gitlab.com/swarmfund/api/db2/api"
	"gitlab.com/swarmfund/api/render/hal"
	"gitlab.com/swarmfund/api/render/problem"
	"gitlab.com/swarmfund/api/resource"
	"gitlab.com/swarmfund/api/utils"
	"gitlab.com/swarmfund/go/xdr"
	horizon "gitlab.com/swarmfund/horizon-connector"
)

type CreateWalletAction struct {
	Action

	Wallet       api.Wallet
	RecoveryCode string

	RecoveryRequest *api.RecoveryRequest
	InitialWallet   *api.Wallet
	Account         *horizon.Account
	WalletExists    bool

	Resource resource.Wallet
}

// JSON is a method for actions.JSON
func (action *CreateWalletAction) JSON() {
	action.Do(
		action.ValidateBodyType,
		action.loadParams,
		action.checkWalletExists,
	)

	if action.WalletExists {
		// recovery flow
		action.Do(
			action.loadRecoveryRequest,
			action.checkCode,
			action.loadAccount,
			action.loadInitialWallet,
			action.checkTFA,
			action.prepareRecoveryWallet,
			action.createWallet,
			action.postCreateWallet,
		)
	}

	if !action.WalletExists && action.Wallet.AccountID != action.Wallet.CurrentAccountID {
		// pending wallet flow
		action.Do(
			action.generateToken,
			action.createWallet,
			action.updateSigner,
			action.updateUser,
			action.sendNotification,
		)
	}

	if !action.WalletExists && action.Wallet.AccountID == action.Wallet.CurrentAccountID {
		// signup flow
		action.Do(
			action.generateToken,
			action.createWallet,
			action.sendNotification,
		)
	}

	if action.Err == nil {
		action.Resource.Populate(&action.Wallet)
		hal.Render(action.W, action.Resource)
	}
}

func (action *CreateWalletAction) loadParams() {
	action.Wallet.Username = action.GetRestrictedString("username", 3, 255)
	action.Wallet.AccountID = action.GetNonEmptyString("accountId")

	action.Wallet.CurrentAccountID = action.GetString("currentAccountId")
	if action.Wallet.CurrentAccountID == "" {
		action.Wallet.CurrentAccountID = action.Wallet.AccountID
	}

	action.Wallet.WalletId = action.GetNonEmptyString("walletId")
	action.Wallet.Salt = action.GetByteArray("salt", 16)

	//action.Wallet.Kdf = action.GetNonEmptyString("kdfParams")

	action.Wallet.KeychainData = action.GetNonEmptyString("keychainData")

	//action.Wallet.TfaSalt = action.GetByteArray("tfaSalt", 16)
	//action.Wallet.TfaPublicKey = action.GetNonEmptyString("tfaPublicKey")
	//action.Wallet.TfaKeychainData = action.GetNonEmptyString("tfaKeychainData")

	if action.App.Config().API().NoEmailVerify {
		action.Wallet.Verified = true
	}

	action.RecoveryCode = action.GetString("recovery_code")
}

func (action *CreateWalletAction) checkWalletExists() {
	wallet, err := action.APIQ().Wallet().ByEmail(action.Wallet.Username)
	if err != nil {
		action.Log.WithError(err).Error("failed to get recovery request")
		action.Err = &problem.ServerError
		return
	}

	action.WalletExists = wallet != nil
}

func (action *CreateWalletAction) loadRecoveryRequest() {
	recoveryRequest, err := action.APIQ().Recoveries().ByUsername(action.Wallet.Username)
	if err != nil {
		action.Log.WithError(err).Error("failed to get recovery request")
		action.Err = &problem.ServerError
		return
	}

	if recoveryRequest == nil || !recoveryRequest.CodeShown() {
		// there are no valid recovery request so gtfo
		action.Err = &problem.Conflict
		return
	}

	action.RecoveryRequest = recoveryRequest
}

func (action *CreateWalletAction) checkCode() {
	// checking second-factor code to filter out defected admins trying to steal your wallet
	if action.RecoveryCode == "" {
		action.SetInvalidField("recovery_code", errors.New("required"))
		return
	}

	if action.RecoveryCode != action.RecoveryRequest.Code {
		action.SetInvalidField("recovery_code", errors.New("invalid"))
		return
	}
}

func (action *CreateWalletAction) loadAccount() {
	account, err := action.App.horizon.AccountSigned(action.App.AccountManagerKP(), action.RecoveryRequest.AccountID)
	if err != nil {
		action.Log.WithError(err).Error("failed to get core account")
		action.Err = &problem.ServerError
		return
	}

	if account == nil {
		action.Log.WithField("address", action.RecoveryRequest.AccountID).Error("account expected to exist")
		action.Err = &problem.ServerError
		return
	}

	action.Account = account
}

func (action *CreateWalletAction) loadInitialWallet() {
	wallet, err := action.APIQ().Wallet().ByID(action.RecoveryRequest.WalletID)
	if err != nil {
		action.Log.WithError(err).Error("failed to get initial wallet")
		action.Err = &problem.ServerError
		return
	}

	if wallet == nil {
		action.Log.WithField("id", action.RecoveryRequest.WalletID).Error("wallet expected to exist")
		action.Err = &problem.ServerError
		return
	}

	action.InitialWallet = wallet
}

func (action *CreateWalletAction) checkTFA() {
	if xdr.AccountType(action.Account.AccountType) == xdr.AccountTypeNotVerified {
		// allowing user without approved KYC to recover account with TFA
		action.consumeTFA(action.InitialWallet, api.TFAActionRecovery)
		if action.Err != nil {
			return
		}
	}
}

func (action *CreateWalletAction) prepareRecoveryWallet() {
	action.Wallet.AccountID = action.RecoveryRequest.AccountID
	action.Wallet.Verified = action.InitialWallet.Verified
}

func (action *CreateWalletAction) postCreateWallet() {
	if xdr.AccountType(action.Account.AccountType) == xdr.AccountTypeNotVerified {
		err := action.recoveryApprove(action.RecoveryRequest, action.InitialWallet, &action.Wallet, "")
		if err != nil {
			action.Log.WithError(err).Error("failed to approve recovery request")
			action.Err = &problem.ServerError
			return
		}
	}
}

func (action *CreateWalletAction) generateToken() {
	action.Wallet.VerificationToken = utils.GenerateToken()
}

func (action *CreateWalletAction) createWallet() {
	err := action.APIQ().Wallet().Create(&action.Wallet)
	if err != nil {
		action.Log.WithError(err).Error("failed to save wallet")
		action.Err = &problem.ServerError
		return
	}
}

func (action *CreateWalletAction) updateSigner() {
	err := action.App.horizon.Transaction(&horizon.TransactionBuilder{Source: action.App.MasterKP()}).
		Op(&horizon.RecoverOp{
			AccountID: action.Wallet.AccountID,
			OldSigner: action.Wallet.AccountID,
			NewSigner: action.Wallet.CurrentAccountID,
		}).Sign(action.App.AccountManagerKP()).Submit()

	if err != nil {
		action.Log.WithError(err).Error("failed to update signers")
		action.Err = &problem.ServerError
		return
	}

	err = action.APIQ().Users().ChangeState(action.Wallet.AccountID, api.UserRejected)
	if err != nil {
		action.Log.WithError(err).Error("failed to update user state")
		action.Err = &problem.ServerError
		return
	}
}

func (action *CreateWalletAction) updateUser() {
	err := action.App.APIQ().Users().SetEmail(action.Wallet.AccountID, action.Wallet.Username)
	if err != nil {
		action.Log.WithError(err).Error("failed to update user")
		action.Err = &problem.ServerError
		return
	}
}

func (action *CreateWalletAction) sendNotification() {
	//if !action.Wallet.Verified {
	//	err := action.Notificator().SendVerificationLink(action.Wallet.Username, action.Wallet.VerificationToken)
	//	if err != nil {
	//		action.Log.WithError(err).Error("Failed to sed email verification link")
	//		action.Err = &problem.ServerError
	//		return
	//	}
	//}
}
