package handlers

import (
	"net/http"

	"encoding/json"

	"fmt"

	"github.com/go-chi/chi"
	. "github.com/go-ozzo/ozzo-validation"
	"github.com/google/jsonapi"
	"github.com/pkg/errors"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"gitlab.com/swarmfund/api/db2/api"
	"gitlab.com/swarmfund/api/internal/api/movetoape"
	"gitlab.com/swarmfund/api/internal/types"
	"gitlab.com/swarmfund/go/doorman"
)

type (
	PatchUserRequest struct {
		Address string               `json:"-"`
		Data    PatchUserRequestData `json:"data"`
	}
	PatchUserRequestData struct {
		Type *types.UserType `json:"type"`
	}
)

func NewPatchUserRequest(r *http.Request) (PatchUserRequest, error) {
	request := PatchUserRequest{
		Address: chi.URLParam(r, "address"),
	}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return request, errors.Wrap(err, "failed to unmarshal")
	}
	return request, request.Validate()
}

func (r PatchUserRequest) Validate() error {
	return ValidateStruct(&r,
		Field(&r.Address, Required),
		Field(&r.Data, Required),
	)
}

func (r PatchUserRequestData) Validate() error {
	return ValidateStruct(&r,
		Field(&r.Type),
	)
}

func PatchUser(w http.ResponseWriter, r *http.Request) {
	request, err := NewPatchUserRequest(r)
	if err != nil {
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	// check allowed
	if err := Doorman(r, doorman.SignerOf(request.Address)); err != nil {
		movetoape.RenderDoormanErr(w, err)
		return
	}

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

	if request.Data.Type != nil {
		if user.UserType != types.UserTypeNotVerified && user.UserType != *request.Data.Type {
			ape.RenderErr(w, &jsonapi.ErrorObject{
				Title:  http.StatusText(http.StatusForbidden),
				Status: fmt.Sprintf("%d", http.StatusForbidden),
				Detail: "Changing user type is not allowed",
			})
			return
		}
		user.UserType = *request.Data.Type
	}

	err = UsersQ(r).Transaction(func(q api.UsersQI) error {
		if err := q.Update(user); err != nil {
			return errors.Wrap(err, "failed to update user")
		}

		state := user.CheckState()
		if state != user.State {
			if err := q.ChangeState(user.Address, state); err != nil {
				return errors.Wrap(err, "failed to change user state")
			}
		}
		return nil
	})
	if err != nil {
		Log(r).WithError(err).Error("update tx failed")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	w.WriteHeader(204)
}
