package api

import (
	"gitlab.com/swarmfund/api/render/problem"
)

const (
	verificationStateNotFound     = "not_found"
	verificationStateInvalidToken = "invalid_token"
	verificationStateAccepted     = "accepted"
)

type VerifyWalletAction struct {
	Action
	State    string
	Username string
	Token    string
}

func (action *VerifyWalletAction) JSON() {
	action.Do(
		action.ValidateBodyType,
		action.loadParams,
		action.performRequest,
		func() {
			// TODO
			//requestUrl, err := url.Parse(action.App.config.Notificator.EmailConfirmation.ClientURL)
			//if err != nil {
			//	action.Log.WithError(err).Error("Failed to parse url")
			//	action.Err = &problem.ServerError
			//	return
			//}
			//
			//requestUrl.Path = fmt.Sprintf("/login/%s", action.State)
			//hal.Redirect(action.W, action.R, requestUrl.String())
		})
}

func (action *VerifyWalletAction) loadParams() {
	action.ValidateBodyType()
	action.Token = action.GetNonEmptyString("token")
	action.Username = action.GetNonEmptyString("username")
}

func (action *VerifyWalletAction) performRequest() {
	wallet, err := action.APIQ().Wallet().ByEmail(action.Username)
	if err != nil {
		action.Log.WithError(err).Error("Unable to get wallet")
		action.Err = &problem.ServerError
		return
	}

	if wallet == nil {
		action.State = verificationStateNotFound
		return
	}

	if wallet.VerificationToken != action.Token {
		action.State = verificationStateInvalidToken
		return
	}

	err = action.APIQ().Wallet().Verify(wallet.Id)

	if err != nil {
		action.Log.WithError(err).Error("Unable to verify wallet")
		action.Err = &problem.ServerError
		return
	}

	action.State = verificationStateAccepted
}
