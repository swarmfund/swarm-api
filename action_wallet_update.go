package api

import (
	"github.com/pkg/errors"
	"gitlab.com/swarmfund/api/db2/api"
	"gitlab.com/swarmfund/api/render/hal"
	"gitlab.com/swarmfund/api/render/problem"
	"gitlab.com/swarmfund/horizon-connector"
)

type UpdateWalletRequest struct {
	Username         string `json:"username" valid:"required"`
	WalletId         string `json:"walletId" valid:"required"`
	OldWalletId      string `json:"oldWalletId" valid:"required"`
	Salt             string `json:"salt" valid:"required"`
	KeychainData     string `json:"keychainData" valid:"required"`
	CurrentAccountID string `json:"accountId" valid:"required"`
	TfaSalt          string `json:"tfaSalt" valid:"required"`
	TfaKeychainData  string `json:"tfaKeychainData" valid:"required"`
	TfaPublicKey     string `json:"tfaPublicKey" valid:"required"`
	Transaction      string `json:"tx" valid:"required"`
}

type UpdateWalletAction struct {
	Action
	Request UpdateWalletRequest
	Wallet  *api.Wallet
}

// JSON is a method for actions.JSON
func (action *UpdateWalletAction) JSON() {
	action.Do(
		action.ValidateBodyType,
		action.loadParams,
		action.loadWallet,
		action.checkAllowed,
		action.submitTX,
	)
}

func (action *UpdateWalletAction) loadParams() {
	action.UnmarshalBody(&action.Request)
}

func (action *UpdateWalletAction) loadWallet() {
	var err error
	action.Wallet, err = action.APIQ().Wallet().ByEmail(action.Request.Username)
	if err != nil {
		action.Log.WithError(err).Error("Failed to get wallet from db")
		action.Err = &problem.ServerError
		return
	}

	if action.Wallet == nil {
		action.Err = &problem.NotFound
		return
	}

	if action.Request.OldWalletId != action.Wallet.WalletId {
		action.SetInvalidField("walletId", errors.New(" mismatched"))
		return
	}
}

func (action *UpdateWalletAction) checkAllowed() {
	action.checkSignerConstraints(
		SignerOf(action.Wallet.CurrentAccountID),
	)
}

func (action *UpdateWalletAction) submitTX() {
	err := action.App.horizon.SubmitTX(action.Request.Transaction)
	if err != nil {
		entry := action.Log.WithError(err)
		if serr, ok := errors.Cause(err).(horizon.SubmitError); ok {
			switch serr.ResponseCode() {
			case 403:
				action.W.WriteHeader(403)
				action.W.Write(serr.ResponseBody())
				return
			default:
				entry = entry.
					WithField("tx code", serr.TransactionCode()).
					WithField("op codes", serr.OperationCodes())
			}
		}
		entry.Error("failed to submit wallet update tx")
		action.Err = &problem.ServerError
		return
	}

	// update wallet
	action.Wallet.WalletId = action.Request.WalletId

	action.Wallet.Salt = action.Request.Salt
	action.Wallet.KeychainData = action.Request.KeychainData
	action.Wallet.CurrentAccountID = action.Request.CurrentAccountID

	//action.Wallet.TfaSalt = action.Request.TfaSalt
	//action.Wallet.TfaKeychainData = action.Request.TfaKeychainData
	//action.Wallet.TfaPublicKey = action.Request.TfaPublicKey

	err = action.APIQ().Wallet().Update(action.Wallet)
	if err != nil {
		action.Log.WithError(err).Error("Failed to add wallet into database")
		action.Err = &problem.ServerError
		return
	}

	hal.Render(action.W, &problem.Success)
}
