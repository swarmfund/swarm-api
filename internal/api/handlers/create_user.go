package handlers

import (
	"net/http"

	"github.com/go-chi/chi"
	. "github.com/go-ozzo/ozzo-validation"
	"github.com/pkg/errors"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"gitlab.com/swarmfund/api/db2/api"
	"gitlab.com/swarmfund/api/internal/api/movetoape"
	"gitlab.com/swarmfund/api/internal/types"
	"gitlab.com/swarmfund/go/doorman"
	"gitlab.com/swarmfund/go/xdr"
	"gitlab.com/swarmfund/horizon-connector"
)

type CreateUserRequest struct {
	Address types.Address `json:"-"`
}

func NewCreateUserRequest(r *http.Request) (CreateUserRequest, error) {
	request := CreateUserRequest{
		Address: types.Address(chi.URLParam(r, "address")),
	}
	return request, request.Validate()
}

func (r *CreateUserRequest) Validate() error {
	return ValidateStruct(r,
		Field(&r.Address, Required),
	)
}

func CreateUser(w http.ResponseWriter, r *http.Request) {
	request, err := NewCreateUserRequest(r)
	if err != nil {
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	if err := Doorman(r, doorman.SignerOf(string(request.Address))); err != nil {
		movetoape.RenderDoormanErr(w, err)
		return
	}

	// wallet should exists and be verified when creating user
	wallet, err := WalletQ(r).ByAccountID(request.Address)
	if err != nil {
		Log(r).WithError(err).Error("failed to get wallet")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	if wallet == nil || !wallet.Verified {
		ape.RenderErr(w, movetoape.Forbidden("verification_required"))
		return
	}

	err = UsersQ(r).Transaction(func(q api.UsersQI) error {
		err := q.Create(&api.User{
			Address: request.Address,
			Email:   wallet.Username,
			// everybody is created equal
			UserType: types.UserTypeNotVerified,
			State:    types.UserStateNil,
		})
		if err != nil {
			return errors.Wrap(err, "failed to insert user")
		}

		err = Horizon(r).Transaction(&horizon.TransactionBuilder{Source: AccountManagerKP(r)}).
			Op(&horizon.CreateAccountOp{
				AccountID:   string(request.Address),
				AccountType: xdr.AccountTypeNotVerified,
			}).Sign(AccountManagerKP(r)).Submit()
		if err != nil {
			return errors.Wrap(err, "failed to submit tx")
		}

		return nil
	})
	if err != nil {
		cause := errors.Cause(err)
		if cause == api.ErrUsersConflict {
			ape.RenderErr(w, problems.Conflict())
			return
		}
		Log(r).WithError(err).Error("failed to create user")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
