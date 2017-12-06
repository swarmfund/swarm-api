package api

import (
	"encoding/json"
	"errors"

	"gitlab.com/swarmfund/api/db2/api"
	"gitlab.com/swarmfund/api/render/hal"
	"gitlab.com/swarmfund/api/render/problem"
)

type EnableTFABackendRequest struct {
	WalletID string `json:"wallet_id" valid:"required"`
	//Type     api.TFABackend  `json:"type" valid:"required"`
	Details json.RawMessage `json:"details"`
}

type EnableTFABackendAction struct {
	Action

	Request  EnableTFABackendRequest
	Resource map[string]interface{} // actual response depends on backend
	wallet   *api.Wallet
}

func (action *EnableTFABackendAction) JSON() {
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

func (action *EnableTFABackendAction) loadParams() {
	action.UnmarshalBody(&action.Request)
}

func (action *EnableTFABackendAction) loadWallet() {
	wallet, err := action.APIQ().Wallet().ByWalletID(action.Request.WalletID)
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

func (action *EnableTFABackendAction) checkAllowed() {
	action.checkSignerConstraints(
		SignerOf(action.wallet.AccountID),
	)
}

func (action *EnableTFABackendAction) performRequest() {
	//action.Resource = map[string]interface{}{}
	//var backend api.Backend
	//switch action.Request.Type {
	//case api.TFABackendGoogleTOTP:
	//	key, err := totp.Generate(totp.GenerateOpts{
	//		Issuer:      "Swarm Fund",
	//		AccountName: action.wallet.Username,
	//	})
	//	if err != nil {
	//		action.Log.WithError(err).Error("failed to generate totp key")
	//		action.Err = &problem.ServerError
	//		return
	//	}
	//	details := tfa.GoogleTOPTDetails{
	//		Secret: key.Secret(),
	//	}
	//	action.Resource["details"] = tfa.GoogleTOPTDetails{
	//		Secret:     key.String(),
	//		SecretSeed: key.Secret(),
	//	}
	//	bytes, err := json.Marshal(&details)
	//	if err != nil {
	//		action.Log.WithError(err).Error("failed to marshal details")
	//		action.Err = &problem.ServerError
	//		return
	//	}
	//	backend = api.Backend{
	//		BackendType: action.Request.Type,
	//		Details:     bytes,
	//	}
	//default:
	//	action.SetInvalidField("type", errors.New("invalid backend type"))
	//}
	//
	//id, err := action.APIQ().TFA().EnableBackend(action.wallet, &backend)
	//if err != nil {
	//	action.Log.WithError(err).Error("failed to add backend")
	//	action.Err = &problem.ServerError
	//	return
	//}
	//action.Resource["id"] = id
}
