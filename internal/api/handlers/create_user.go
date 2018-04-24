package handlers

import (
	"net/http"

	"github.com/go-chi/chi"
	. "github.com/go-ozzo/ozzo-validation"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/swarmfund/api/db2/api"
	"gitlab.com/swarmfund/api/internal/api/movetoape"
	"gitlab.com/swarmfund/api/internal/hose"
	"gitlab.com/swarmfund/api/internal/types"
	"gitlab.com/tokend/go/doorman"
	"gitlab.com/tokend/go/xdr"
	"gitlab.com/tokend/go/xdrbuild"
)

var (
	ErrWalletNotVerified = errors.New("wallet not verified")
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

func performUserCreate(r *http.Request, wallet *api.Wallet) error {
	// wallet should exists and be verified when creating user
	if wallet == nil || !wallet.Verified {
		return ErrWalletNotVerified
	}

	err := UsersQ(r).Transaction(func(q api.UsersQI) error {
		user := &api.User{
			Address: wallet.AccountID,
			Email:   wallet.Username,
		}
		state := api.UserStateUpdate{
			Address: wallet.AccountID,
			Type:    &types.DefaultUserType,
			State:   &types.DefaultUserState,
		}

		if err := q.Create(user); err != nil {
			return errors.Wrap(err, "failed to insert user")
		}

		if err := q.SetState(state); err != nil {
			return errors.Wrap(err, "failed to insert user state")
		}

		envelope, err := Transaction(r).
			Op(xdrbuild.CreateAccountOp{
				Address:     string(wallet.AccountID),
				AccountType: xdr.AccountTypeNotVerified,
				Recovery:    string(wallet.RecoveryAddress),
				Referrer:    wallet.Referrer,
			}).Marshal()
		if err != nil {
			return errors.Wrap(err, "failed to build tx envelope")
		}
		if result := Horizon(r).Submitter().Submit(r.Context(), envelope); result.Err != nil {
			// TODO assert fail reasons
			return errors.Wrap(result.Err, "failed to submit tx", result.GetLoganFields())
		}

		// dispatch user create event
		UserBusDispatch(r, hose.UserEvent{
			Type: hose.UserEventTypeCreated,
			User: hose.User{
				Email:   wallet.Username,
				Address: wallet.AccountID,
			},
		})
		return nil
	})
	if err != nil {
		cause := errors.Cause(err)
		if cause == api.ErrUsersConflict {
			return cause
		}
		return err
	}

	return nil
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

	wallet, err := WalletQ(r).ByAccountID(request.Address)
	if err != nil {
		Log(r).WithError(err).Error("failed to get wallet")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	if err := performUserCreate(r, wallet); err != nil {
		switch errors.Cause(err) {
		case ErrWalletNotVerified:
			ape.RenderErr(w, movetoape.Forbidden("verification_required"))
		case api.ErrUsersConflict:
			ape.RenderErr(w, problems.Conflict())
		default:
			Log(r).WithError(err).Error("failed to create user")
			ape.RenderErr(w, problems.InternalError())
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
