package handlers

import (
	"net/http"

	"encoding/json"

	"github.com/go-chi/chi"
	. "github.com/go-ozzo/ozzo-validation"
	"github.com/pkg/errors"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"gitlab.com/swarmfund/api/db2/api"
	"gitlab.com/swarmfund/api/internal/kyc"
	"gitlab.com/swarmfund/api/internal/lorem"
)

type (
	CreateKYCEntityRequest struct {
		Address string     `json:"-"`
		Data    kyc.Entity `json:"data"`
	}
)

func NewCreateKYCEntityRequest(r *http.Request) (CreateKYCEntityRequest, error) {
	request := CreateKYCEntityRequest{
		Address: chi.URLParam(r, "address"),
	}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return request, errors.Wrap(err, "failed to unmarshal")
	}
	return request, request.Validate()
}

func (r CreateKYCEntityRequest) Validate() error {
	return ValidateStruct(&r,
		Field(&r.Address, Required),
		Field(&r.Data, Required),
	)
}

func CreateKYCEntity(w http.ResponseWriter, r *http.Request) {
	request, err := NewCreateKYCEntityRequest(r)
	if err != nil {
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	// TODO check allowed

	// load user
	user, err := UsersQ(r).ByAddress(request.Address)
	if err != nil {
		Log(r).WithError(err).Error("failed to get user")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	if user == nil {
		ape.RenderErr(w, problems.NotFound())
		return
	}

	err = UsersQ(r).Transaction(func(q api.UsersQI) error {
		if err = q.KYC().Create(user.ID, lorem.ULID(), request.Data); err != nil {
			return errors.Wrap(err, "failed to create entity")
		}

		return nil
	})
	if err != nil {
		cause := errors.Cause(err)
		if cause == api.ErrKYCEntitiesConstraintViolated {
			ape.RenderErr(w, problems.Conflict())
			return
		}
		Log(r).WithError(err).Error("db tx failed")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	w.WriteHeader(204)
}
