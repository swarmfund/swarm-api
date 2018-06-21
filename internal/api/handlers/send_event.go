package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	. "github.com/go-ozzo/ozzo-validation"
	"github.com/pkg/errors"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
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

func NewSendEventRequest(r *http.Request) (SendEventRequest, error) {
	request := SendEventRequest{Address: types.Address(chi.URLParam(r, "address"))}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return request, errors.Wrap(err, "failed to unmarshal")
	}
	return request, request.Validate()
}

func (r SendEventRequest) Validate() error {
	return ValidateStruct(&r,
		Field(&r.Data, Required),
	)
}

func (r SendEventRequestData) Validate() error {
	return ValidateStruct(&r,
		Field(&r.Event, Required),
	)
}

type UserData struct {
	Name  string
	Email string
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
		ape.RenderErr(w, problems.InternalError())
		return
	}
	if kycData == nil {
		w.WriteHeader(204)
		return
	}

	name := kycData.FirstName + " " + kycData.LastName

	_, err = Salesforce(r).SendEvent(sphere, actionName, time.Now(), name, user.Email, 0, "")
	if err != nil {
		ape.RenderErr(w, problems.BadRequest(err)...)
	}

	w.WriteHeader(204)

}
