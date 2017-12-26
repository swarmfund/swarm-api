package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

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
		Type          *types.UserType               `json:"type"`
		Attributes    PatchUserRequestAttributes    `json:"attributes"`
		Relationships PatchUserRequestRelationships `json:"relationships"`
	}
	PatchUserRequestAttributes struct {
		State *types.UserState `json:"state"`
	}
	PatchUserRequestRelationships struct {
		Transaction struct {
			Data struct {
				Attributes struct {
					Envelope string `json:"envelope"`
				}
			} `json:"data"`
		} `json:"transaction"`
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

func (r PatchUserRequestAttributes) Validate() error {
	return ValidateStruct(&r,
		Field(&r.State),
	)
}

func PatchUser(w http.ResponseWriter, r *http.Request) {
	request, err := NewPatchUserRequest(r)
	if err != nil {
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	// check is allowed
	signedByAdmin := false
	if err := Doorman(r, doorman.SignerOf(request.Address)); err != nil {
		// request not signed by user, let's check admin
		master := CoreInfo(r).GetMasterAccountID()
		Log(r).WithField("master", master).Debug("checking admin signature")
		if err := Doorman(r, doorman.SignerOf(master)); err != nil {
			// not by admin either
			movetoape.RenderDoormanErr(w, err)
			return
		}
		// seems like admin to me
		signedByAdmin = true
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

	Log(r).WithField("admin", signedByAdmin).Debug("attempting user update")

	if signedByAdmin {
		// admin can update state
		if request.Data.Attributes.State != nil && user.State != *request.Data.Attributes.State {
			// only when user is waiting for approval
			if user.State != types.UserStateWaitingForApproval {
				ape.RenderErr(w, problems.BadRequest(Errors{
					"/data/attributes/state": errors.New("state transition is not allowed"),
				})...)
				return
			}
			// transaction is required when approving
			if *request.Data.Attributes.State == types.UserStateApproved && request.Data.Relationships.Transaction.Data.Attributes.Envelope == "" {
				ape.RenderErr(w, problems.BadRequest(Errors{
					"/data/relationships/transaction/data/attributes/envelope": errors.New("required when approving"),
				})...)
				return
			}

			// TODO reject reason required when rejecting

			// all checks have passed, updating user state
			user.State = *request.Data.Attributes.State
		}
	} else {
		// user could update type
		if request.Data.Type != nil && user.UserType != *request.Data.Type {
			// only when it's currently not verified
			if user.UserType != types.UserTypeNotVerified {
				ape.RenderErr(w, &jsonapi.ErrorObject{
					Title:  http.StatusText(http.StatusForbidden),
					Status: fmt.Sprintf("%d", http.StatusForbidden),
					Detail: "Changing user type is not allowed",
				})
				return
			}

			user.UserType = *request.Data.Type
			user.State = user.CheckState()
		}

		// user could update state
		if request.Data.Attributes.State != nil && user.State != *request.Data.Attributes.State {
			// only to waiting for approval
			if *request.Data.Attributes.State != types.UserStateWaitingForApproval {
				ape.RenderErr(w, problems.BadRequest(Errors{
					"/data/attributes/state": errors.New("state transition is not allowed"),
				})...)
				return
			}
			// check if user is really able to change state
			if user.CheckState() != types.UserStateWaitingForApproval {
				ape.RenderErr(w, problems.BadRequest(Errors{
					"/data/attributes/state": errors.New("state transition is not allowed"),
				})...)
				return
			}

			// all checks have passed, updating user state
			user.State = *request.Data.Attributes.State
		}
	}

	err = UsersQ(r).Transaction(func(q api.UsersQI) error {
		if err := q.Update(user); err != nil {
			return errors.Wrap(err, "failed to update user")
		}

		if tx := request.Data.Relationships.Transaction.Data.Attributes.Envelope; tx != "" {
			if err := Horizon(r).SubmitTX(tx); err != nil {
				return errors.Wrap(err, "failed to submit transaction")
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
