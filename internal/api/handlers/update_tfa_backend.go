package handlers

import (
	"net/http"

	"encoding/json"

	"strconv"

	"fmt"

	"github.com/go-chi/chi"
	. "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"github.com/google/jsonapi"
	"github.com/pkg/errors"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"gitlab.com/swarmfund/api/internal/api/movetoape"
	"gitlab.com/swarmfund/api/internal/secondfactor"
	"gitlab.com/tokend/go/doorman"
)

type (
	UpdateTFABackendRequestData struct {
		Attributes UpdateTFABackendRequestAttributes `json:"attributes"`
	}
	UpdateTFABackendRequestAttributes struct {
		Priority int `json:"priority"`
	}
	UpdateTFABackendRequest struct {
		WalletID  string                      `json:"-"`
		BackendID string                      `json:"-"`
		Data      UpdateTFABackendRequestData `json:"data"`
	}
)

func (r UpdateTFABackendRequestData) Validate() error {
	return ValidateStruct(&r,
		Field(&r.Attributes),
	)
}

func (r UpdateTFABackendRequestAttributes) Validate() error {
	return ValidateStruct(&r,
		Field(&r.Priority, Required, Min(0)),
	)
}

func NewUpdateTFABackendRequest(r *http.Request) (UpdateTFABackendRequest, error) {
	request := UpdateTFABackendRequest{
		WalletID:  chi.URLParam(r, "wallet-id"),
		BackendID: chi.URLParam(r, "backend"),
	}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return request, errors.Wrap(err, "failed to unmarshal body")
	}
	return request, request.Validate()
}

func (r UpdateTFABackendRequest) Validate() error {
	return ValidateStruct(&r,
		Field(&r.WalletID, Required),
		Field(&r.BackendID, Required, is.Int),
		Field(&r.Data),
	)
}

func UpdateWalletFactor(w http.ResponseWriter, r *http.Request) {
	request, err := NewUpdateTFABackendRequest(r)
	if err != nil {
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	// load wallet
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
	if err := Doorman(r, doorman.SignerOf(string(wallet.AccountID))); err != nil {
		movetoape.RenderDoormanErr(w, err)
		return
	}

	// load backend
	bid, err := strconv.ParseInt(request.BackendID, 10, 64)
	if err != nil {
		Log(r).WithError(err).Error("unexpected id parse fail")
		ape.RenderErr(w, problems.InternalError())
		return
	}
	record, err := TFAQ(r).Backend(bid)
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

	// check tfa
	if err := secondfactor.NewConsumer(TFAQ(r)).WithBackend(backend).Consume(r, wallet); err != nil {
		RenderFactorConsumeError(w, r, err)
		return
	}

	// update record
	if err := TFAQ(r).SetBackendPriority(record.ID, request.Data.Attributes.Priority); err != nil {
		Log(r).WithError(err).Error("failed to update factor")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	w.WriteHeader(204)
}

func RenderFactorConsumeError(w http.ResponseWriter, r *http.Request, err error) {
	cause := errors.Cause(err)
	switch terr := cause.(type) {
	case *secondfactor.FactorRequiredErr:
		ape.RenderErr(w, &jsonapi.ErrorObject{
			Code:   "tfa_required",
			Title:  http.StatusText(http.StatusForbidden),
			Status: fmt.Sprintf("%d", http.StatusForbidden),
			Detail: "Additional factor required",
			Meta:   terr.Meta(),
		})
	default:
		Log(r).WithError(err).Error("failed to consume second factor")
		ape.RenderErr(w, problems.InternalError())
	}
}
