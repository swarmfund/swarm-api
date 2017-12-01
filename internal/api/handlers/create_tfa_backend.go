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
	"gitlab.com/swarmfund/api/internal/api/movetoape"
	"gitlab.com/swarmfund/api/internal/secondfactor"
	"gitlab.com/swarmfund/api/internal/types"
	"gitlab.com/swarmfund/api/tfa"
	"gitlab.com/swarmfund/go/doorman"
)

type (
	CreateBackendRequestData struct {
		Type types.WalletFactorType `json:"type"`
	}
	CreateBackendRequest struct {
		WalletID string                   `json:"-"`
		Data     CreateBackendRequestData `json:"data"`
	}
)

func NewCreateBackendRequest(r *http.Request) (CreateBackendRequest, error) {
	request := CreateBackendRequest{
		WalletID: chi.URLParam(r, "wallet-id"),
	}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return request, errors.Wrap(err, "failed to unmarshal")
	}
	return request, request.Validate()
}

func (r CreateBackendRequestData) Validate() error {
	return ValidateStruct(&r,
		Field(&r.Type),
	)
}

func (r CreateBackendRequest) Validate() error {
	return ValidateStruct(&r,
		Field(&r.WalletID, Required),
		Field(&r.Data, Required),
	)
}

type (
	TFABackendData struct {
		Type       types.WalletFactorType `json:"type"`
		Attributes map[string]interface{} `json:"attributes"`
	}

	TFABackend struct {
		ID   string         `json:"id"`
		Data TFABackendData `json:"data"`
	}
)

func CreateTFABackend(w http.ResponseWriter, r *http.Request) {
	request, err := NewCreateBackendRequest(r)
	if err != nil {
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	wallet, err := WalletQ(r).ByWalletID(request.WalletID)
	if err != nil {
		Log(r).WithError(err).Error("failed to get wallet")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	if wallet == nil {
		ape.RenderErr(w, problems.NotFound())
		return
	}

	// check allowed
	if err := Doorman(r, doorman.SignerOf(wallet.AccountID)); err != nil {
		movetoape.RenderDoormanErr(w, err)
		return
	}

	var backend tfa.Backend

	switch request.Data.Type {
	case types.WalletFactorTOTP:
		backend, err = tfa.NewTOTPBackend("{{ .Project }}", wallet.Username)
		if err != nil {
			Log(r).WithError(err).Error("failed to generate backend")
			ape.RenderErr(w, problems.InternalError())
			return
		}
	case types.WalletFactorPassword:
		// factor should be created during signup
		ape.RenderErr(w, problems.Conflict())
		return
	default:
		Log(r).WithField("type", request.Data.Type).Error("unable to handle backend type")
		// TODO make 501 not implemented
		ape.RenderErr(w, problems.InternalError())
		return
	}

	// after all checks have passed, check 2fa
	if err := secondfactor.NewConsumer(TFAQ(r)).WithBackendType(types.WalletFactorPassword).Consume(r, wallet); err != nil {
		RenderFactorConsumeError(w, r, err)
		return
	}

	id, err := TFAQ(r).CreateBackend(wallet.WalletId, backend)
	if err != nil {
		cause := errors.Cause(err)

		if cause == api.ErrWalletBackendConflict {
			ape.RenderErr(w, problems.Conflict())
			return
		}

		Log(r).WithError(err).Error("failed to save backend")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	type (
		CreateTFABackendResponseData struct {
			ID         int64                  `json:"id"`
			Type       types.WalletFactorType `json:"type"`
			Attributes map[string]interface{} `json:"attributes"`
		}
		CreateTFABackendResponse struct {
			Data CreateTFABackendResponseData `json:"data"`
		}
	)

	response := CreateTFABackendResponse{
		Data: CreateTFABackendResponseData{
			ID:         *id,
			Type:       request.Data.Type,
			Attributes: backend.Attributes(),
		},
	}

	w.WriteHeader(201)
	json.NewEncoder(w).Encode(&response)
}
