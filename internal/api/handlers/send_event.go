package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	validation "github.com/go-ozzo/ozzo-validation"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/swarmfund/api/internal/kycv2"
	"gitlab.com/swarmfund/api/internal/types"
)

type (
	SendEventRequest struct {
		Address types.Address        `json:"-"`
		Data    SendEventRequestData `json:"data"`
	}
	SendEventRequestData struct {
		Event types.Event `json:"type"`
	}
)

func buildSalesforceError(errs []string) error {
	s := ""
	for i, err := range errs {
		if i > 0 {
			s += "; "
		}
		s += err
	}
	s += "."
	return errors.New(s)
}

func NewSendEventRequest(r *http.Request) (SendEventRequest, error) {
	request := SendEventRequest{Address: types.Address(chi.URLParam(r, "address"))}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return request, errors.Wrap(err, "failed to unmarshal")
	}
	return request, request.Validate()
}

func (r SendEventRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.Data, validation.Required),
	)
}

func (r SendEventRequestData) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.Event, validation.Required),
	)
}

func SendEvent(w http.ResponseWriter, r *http.Request) {
	request, err := NewSendEventRequest(r)
	if err != nil {
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	sphere := request.Data.Event.GetSalesforceSphere()
	actionName := request.Data.Event.GetSalesforceActionName()

	user, err := UsersQ(r).ByAddress(string(request.Address))
	if err != nil {
		Log(r).WithError(err).Error("failed to get user by address")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	if user == nil {
		ape.RenderErr(w, problems.NotFound())
		return
	}

	if user.KYCBlobValue == nil {
		w.WriteHeader(204)
		return
	}

	kycData, err := kyc.ParseKYCData(*user.KYCBlobValue)
	if err != nil {
		Log(r).WithError(err).Error("failed to parse kyc data")
		ape.RenderErr(w, problems.InternalError())
		return

	}

	var name string
	if kycData != nil {
		name = kycData.FirstName + " " + kycData.LastName
	}

	resp, err := Salesforce(r).SendEvent(sphere, actionName, time.Now(), name, user.Email, 0, "")
	if err != nil {
		Log(r).WithError(err).Error("failed to send event")
		ape.RenderErr(w, problems.InternalError())
		return
	}
	if resp.Errors != nil {
		ape.RenderErr(w, problems.BadRequest(buildSalesforceError(resp.Errors))...)
		return
	}

	w.WriteHeader(204)

}
