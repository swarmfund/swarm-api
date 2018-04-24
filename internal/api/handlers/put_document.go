package handlers

import (
	"net/http"

	"encoding/json"

	"github.com/go-chi/chi"
	. "github.com/go-ozzo/ozzo-validation"
	"github.com/pkg/errors"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"gitlab.com/swarmfund/api/internal/api/movetoape"
	"gitlab.com/swarmfund/api/internal/api/resources"
	storage2 "gitlab.com/swarmfund/api/internal/storage"
	"gitlab.com/swarmfund/api/internal/types"
	"gitlab.com/swarmfund/api/storage"
	"gitlab.com/tokend/go/doorman"
)

type (
	PutDocumentRequest struct {
		AccountID types.Address          `json:"-"`
		Data      PutDocumentRequestData `json:"data"`
	}
	PutDocumentRequestData struct {
		Type       types.DocumentType           `json:"type"`
		Attributes PutDocumentRequestAttributes `json:"attributes"`
	}
	PutDocumentRequestAttributes struct {
		ContentType string `json:"content_type"`
	}
	PutDocumentResponse struct {
		Data resources.UploadForm `json:"data"`
	}
)

func NewPutDocumentRequest(r *http.Request) (PutDocumentRequest, error) {
	request := PutDocumentRequest{
		AccountID: types.Address(chi.URLParam(r, "address")),
	}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return request, errors.Wrap(err, "failed to unmarshal")
	}
	return request, request.Validate()
}

func (r PutDocumentRequest) Validate() error {
	return ValidateStruct(&r,
		Field(&r.AccountID, Required),
		Field(&r.Data, Required),
	)
}

func (r PutDocumentRequestData) Validate() error {
	return ValidateStruct(&r,
		Field(&r.Type, Required),
		Field(&r.Attributes, Required),
	)
}

func (r PutDocumentRequestAttributes) Validate() error {
	return ValidateStruct(&r,
		Field(&r.ContentType, Required, By(storage.IsAllowedContentType)),
	)
}

func PutDocument(w http.ResponseWriter, r *http.Request) {
	request, err := NewPutDocumentRequest(r)
	if err != nil {
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	user, err := UsersQ(r).ByAddress(string(request.AccountID))
	if err != nil {
		Log(r).WithError(err).Error("failed to get user")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	var ownerID int64
	if user != nil {
		// user exists means docs has an owner and we check against it's signers
		if err := Doorman(r, doorman.SignerOf(string(user.Address))); err != nil {
			movetoape.RenderDoormanErr(w, err)
			return
		}
		ownerID = user.ID
	} else {
		if err := Doorman(r, doorman.SignerOf(CoreInfo(r).GetMasterAccountID())); err != nil {
			movetoape.RenderDoormanErr(w, err)
			return
		}
	}

	if !storage.IsContentTypeAllowed(request.Data.Attributes.ContentType) {
		ape.RenderErr(w, problems.BadRequest(Errors{
			"/data/attributes/content_type": errors.New("not allowed"),
		})...)
		return
	}

	key := storage2.NewKey(ownerID, request.Data.Type)

	form, err := Storage(r).UploadFormData(
		storage2.EncodeKey(key),
	)
	if err != nil {
		Log(r).WithError(err).Error("failed to build form data")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	response := PutDocumentResponse{
		Data: resources.NewUploadForm(form),
	}
	json.NewEncoder(w).Encode(response)
}
