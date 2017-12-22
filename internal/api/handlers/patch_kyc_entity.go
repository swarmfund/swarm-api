package handlers

import (
	"net/http"

	"encoding/json"

	"github.com/go-chi/chi"
	. "github.com/go-ozzo/ozzo-validation"
	"github.com/pkg/errors"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"gitlab.com/swarmfund/api/internal/kyc"
)

type (
	PatchKYCEntityRequest struct {
		Address  string                    `json:"-"`
		EntityID string                    `json:"-"`
		Data     PatchKYCEntityRequestData `json:"data"`
	}
	PatchKYCEntityRequestData struct {
		kyc.Entity
		ID int64 `json:"id,string"`
	}
)

func NewPatchKYCEntityRequest(r *http.Request) (PatchKYCEntityRequest, error) {
	request := PatchKYCEntityRequest{
		Address:  chi.URLParam(r, "address"),
		EntityID: chi.URLParam(r, "entity"),
	}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return request, errors.Wrap(err, "failed to unmarshal")
	}
	return request, request.Validate()
}

func (r PatchKYCEntityRequestData) Validate() error {
	return ValidateStruct(&r,
		Field(&r.ID, Required),
		Field(&r.Type, Required),
	)
}

func (r PatchKYCEntityRequest) Validate() error {
	return ValidateStruct(&r,
		Field(&r.Address, Required),
		Field(&r.EntityID, Required),
	)
}

func PatchKYCEntity(w http.ResponseWriter, r *http.Request) {
	request, err := NewPatchKYCEntityRequest(r)
	if err != nil {
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	// TODO check allowed

	// load existing entity
	entity, err := UsersQ(r).KYC().Get(request.EntityID)
	if err != nil {
		Log(r).WithError(err).Error("failed to get entity")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	if entity == nil {
		ape.RenderErr(w, problems.NotFound())
		return
	}

	if entity.Entity.Type != request.Data.Type {
		ape.RenderErr(w, problems.BadRequest(Errors{
			"/data/type": errors.New("invalid"),
		})...)
		return
	}

	err = UsersQ(r).KYC().Update(request.EntityID, request.Data.Entity)
	if err != nil {
		Log(r).WithError(err).Error("failed to update entity")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	w.WriteHeader(204)
}
