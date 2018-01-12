package api

import (
	errors "errors"
	"net/http"

	"gitlab.com/swarmfund/api/db2/api"
	"gitlab.com/swarmfund/api/internal/lorem"
	"gitlab.com/swarmfund/api/internal/types"
	"gitlab.com/swarmfund/api/render/hal"
	"gitlab.com/swarmfund/api/render/problem"
	"gitlab.com/swarmfund/api/storage"
)

type PutDocumentRequest struct {
	ContentType  string           `json:"content_type" valid:"required"`
	DocumentType api.DocumentType `json:"document_type" valid:"required"`
	EntityID     int64            `json:"entity_id" valid:"optional"`
	WalletID     string           `json:"wallet_id" valid:"optional"`
}

type PutDocumentAction struct {
	Action

	AccountID string
	Request   PutDocumentRequest

	User     *api.User
	Document storage.Document

	Response map[string]string
}

func (action *PutDocumentAction) JSON() {
	action.Do(
		action.checkAvailable,
		action.ValidateBodyType,
		action.loadParams,
		action.checkContentType,
		action.checkAllowed,
		action.loadUser,
		action.performRequest,
		action.signForm,
		func() {
			hal.Render(action.W, action.Response)
		},
	)
}

func (action *PutDocumentAction) checkAvailable() {
	if action.App.Config().Storage().Disabled {
		action.Log.Warn("storage service disabled")
		action.Err = &problem.P{
			Status: http.StatusServiceUnavailable,
		}
		return
	}
}

func (action *PutDocumentAction) loadParams() {
	action.UnmarshalBody(&action.Request)
	action.AccountID = action.GetNonEmptyString("id")
}

func (action *PutDocumentAction) checkAllowed() {
	switch {
	default:
		action.checkSignerConstraints(
			SignerOf(action.AccountID),
		)
	}
}

func (action *PutDocumentAction) loadUser() {
	user, err := action.APIQ().Users().ByAddress(action.AccountID)
	if err != nil {
		action.Log.WithError(err).Error("failed to get user")
		action.Err = &problem.ServerError
		return
	}

	if user == nil {
		action.Err = &problem.NotFound
		return
	}

	action.User = user
}

func (action *PutDocumentAction) checkContentType() {
	if !storage.IsContentTypeAllowed(action.Request.ContentType) {
		action.SetInvalidField("content_type", errors.New("Content type not allowed"))
		return
	}
}

func (action *PutDocumentAction) performRequest() {
	action.Document = storage.Document{
		AccountID: types.Address(action.AccountID),
		Type:      action.Request.DocumentType,
		EntityID:  action.Request.EntityID,
		Version:   lorem.Token(),
		Extension: storage.ContentTypeExtension(action.Request.ContentType),
	}

	switch {
	case action.Document.Type.IsKYC():
	case action.Document.Type.IsProofOfIncome():
	default:
		action.SetInvalidField("document_type", errors.New("invalid"))
	}
}

func (action *PutDocumentAction) signForm() {
	//form, err := action.App.Storage().UploadFormData(
	//	string(action.Document.AccountID), action.Document.Key(),
	//)
	//if err != nil {
	//	action.Log.WithError(err).Error("failed to build form data")
	//	action.Err = &problem.ServerError
	//	return
	//}
	//
	//action.Response = form
}
