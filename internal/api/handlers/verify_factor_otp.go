package handlers

import (
	"net/http"

	"encoding/json"

	"strconv"

	"github.com/go-chi/chi"
	. "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"github.com/pkg/errors"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
)

type (
	VerifyFactorOTPAttributes struct {
		Token string `json:"token"`
		OTP   string `json:"otp"`
	}
	VerifyFactorOTPData struct {
		Attributes VerifyFactorOTPAttributes `json:"attributes"`
	}
	VerifyFactorOTPRequest struct {
		WalletID string              `json:"-"`
		FactorID string              `json:"-"`
		Data     VerifyFactorOTPData `json:"data"`
	}
)

func NewVerifyFactorOTPRequest(r *http.Request) (VerifyFactorOTPRequest, error) {
	request := VerifyFactorOTPRequest{
		WalletID: chi.URLParam(r, "wallet-id"),
		FactorID: chi.URLParam(r, "backend"),
	}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return request, errors.Wrap(err, "failed to unmarshal")
	}
	return request, request.Validate()
}

func (r VerifyFactorOTPRequest) Validate() error {
	return ValidateStruct(&r,
		Field(&r.WalletID, Required),
		Field(&r.FactorID, Required, is.Int),
		Field(&r.Data, Required),
	)
}

func (r VerifyFactorOTPData) Validate() error {
	return ValidateStruct(&r,
		Field(&r.Attributes, Required),
	)
}

func (r VerifyFactorOTPAttributes) Validate() error {
	return ValidateStruct(&r,
		Field(&r.Token, Required),
		Field(&r.OTP, Required),
	)
}

func VerifyFactorOTP(w http.ResponseWriter, r *http.Request) {
	request, err := NewVerifyFactorOTPRequest(r)
	if err != nil {
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	// load backend
	fid, err := strconv.ParseInt(request.FactorID, 10, 64)
	if err != nil {
		Log(r).WithError(err).Error("failed to parse factor id")
		ape.RenderErr(w, problems.InternalError())
		return
	}
	record, err := TFAQ(r).Backend(fid)
	if err != nil {
		Log(r).WithError(err).Error("failed to get backend")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	if record == nil {
		ape.RenderErr(w, problems.NotFound())
		return
	}

	backend, err := record.Backend()
	if err != nil {
		Log(r).WithError(err).Error("failed to init backend")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	// load tfa
	tfa, err := TFAQ(r).Get(request.Data.Attributes.Token)
	if err != nil {
		Log(r).WithError(err).Error("failed to get tfa")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	if tfa == nil {
		ape.RenderErr(w, problems.BadRequest(Errors{
			"data/attributes/token": errors.New("not found"),
		})...)
		return
	}

	// verify otp
	ok, err := backend.Verify(request.Data.Attributes.OTP, tfa.Token)
	if err != nil {
		Log(r).WithError(err).Error("failed to verify otp")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	if !ok {
		ape.RenderErr(w, problems.BadRequest(Errors{
			"data/attributes/otp": errors.New("invalid"),
		})...)
		return
	}

	if err = TFAQ(r).Verify(backend.ID(), request.Data.Attributes.Token); err != nil {
		Log(r).WithError(err).Error("failed to mark token as verified")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	w.WriteHeader(204)
}
