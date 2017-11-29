package api

import (
	"gitlab.com/swarmfund/api/render/hal"
	"gitlab.com/swarmfund/api/render/problem"
)

type VerifyTFAAction struct {
	Action

	Token string
	Code  string
}

func (action *VerifyTFAAction) JSON() {
	action.Do(
		action.loadParams,
		action.performRequest,
		func() {
			hal.Render(action.W, problem.Success)
		},
	)
}

func (action *VerifyTFAAction) loadParams() {
	action.Token = action.GetNonEmptyString("token")
	action.Code = action.GetNonEmptyString("code")
}

func (action *VerifyTFAAction) performRequest() {
	otp, err := action.APIQ().TFA().Get(action.Token)
	if err != nil {
		action.Log.WithError(err).Error("failed to get otp")
		action.Err = &problem.ServerError
		return
	}

	if otp == nil {
		action.Err = &problem.BadRequest
		return
	}

	record, err := action.APIQ().TFA().Backend(otp.BackendID)
	if err != nil {
		action.Log.WithError(err).Error("failed to get otp")
		action.Err = &problem.ServerError
		return
	}

	if record == nil {
		action.Log.WithField("otp", otp.ID).Error("can't find backend")
		action.Err = &problem.ServerError
		return
	}

	//backend, err := record.Backend()
	//if err != nil {
	//	action.Log.WithError(err).WithField("backend", record.ID).Error("failed to init backend")
	//	action.Err = &problem.ServerError
	//}

	//ok, err := backend.Verify(action.Code, otp.OTPData)
	//if err != nil {
	//	action.Log.WithError(err).Error("failed to verify otp")
	//	action.Err = &problem.ServerError
	//	return
	//}
	//
	//if !ok {
	//	action.Err = &problem.BadRequest
	//	return
	//}

	err = action.APIQ().TFA().Verify(action.Token)
	if err != nil {
		action.Log.WithError(err).Error("failed to mark token as verified")
		action.Err = &problem.ServerError
		return
	}
}
