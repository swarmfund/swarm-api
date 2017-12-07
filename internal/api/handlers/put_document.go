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
	"gitlab.com/swarmfund/api/internal/api/resources"
	"gitlab.com/swarmfund/api/internal/lorem"
	"gitlab.com/swarmfund/api/internal/types"
	"gitlab.com/swarmfund/api/storage"
)

type (
	PutDocumentRequest struct {
		AccountID types.Address          `json:"-"`
		Data      PutDocumentRequestData `json:"data"`
	}
	PutDocumentRequestData struct {
		Type       string                       `json:"type"`
		Attributes PutDocumentRequestAttributes `json:"attributes"`
	}
	PutDocumentRequestAttributes struct {
		ContentType string `json:"content_type"`
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

	// TODO check allowed

	if !storage.IsContentTypeAllowed(request.Data.Attributes.ContentType) {
		ape.RenderErr(w, problems.BadRequest(Errors{
			"/data/attributes/content_type": errors.New("not allowed"),
		})...)
		return
	}

	document := storage.Document{
		AccountID: request.AccountID,
		Type:      api.DocumentTypeAssetLogo,
		Version:   lorem.Token(),
		Extension: storage.ContentTypeExtension(request.Data.Attributes.ContentType),
	}

	form, err := Storage(r).UploadFormData(
		string(request.AccountID), document.Key(),
	)
	if err != nil {
		Log(r).WithError(err).Error("failed to build form data")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	json.NewEncoder(w).Encode(resources.NewUploadForm(form))
}
