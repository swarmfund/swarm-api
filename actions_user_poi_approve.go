package api

import (
	"errors"

	"gitlab.com/swarmfund/api/db2/api"
	"gitlab.com/swarmfund/api/render/hal"
	"gitlab.com/swarmfund/api/render/problem"
	horizon "gitlab.com/swarmfund/horizon-connector"
)

type UserProofOfIncomeApproveRequest struct {
	DocumentsVersion int64  `json:"documents_version" valid:"required"`
	Approved         bool   `json:"approved"`
	RejectReason     string `json:"reject_reason"`
	TX               string `json:"tx"`
}

type UserProofOfIncomeApproveAction struct {
	Action

	Request   UserProofOfIncomeApproveRequest
	AccountID string
	Version   string

	User     *api.User
	Document *api.Document
}

func (action *UserProofOfIncomeApproveAction) JSON() {
	action.Do(
		action.ValidateBodyType,
		action.loadParams,
		action.checkAllowed,
		action.loadUser,
		action.performRequest,
		func() {
			hal.Render(action.W, problem.Success)
		},
	)
}

func (action *UserProofOfIncomeApproveAction) loadParams() {
	action.UnmarshalBody(&action.Request)
	if action.Err != nil {
		return
	}
	if action.Request.Approved && action.Request.TX == "" {
		action.SetInvalidField("tx", errors.New("should not be empty"))
	}
	if !action.Request.Approved && action.Request.RejectReason == "" {
		action.SetInvalidField("reject_reason", errors.New("should not be empty"))
	}

	action.AccountID = action.GetNonEmptyString("id")
	action.Version = action.GetNonEmptyString("version")
}

func (action *UserProofOfIncomeApproveAction) checkAllowed() {
	action.checkSignerConstraints(
		SignerOf(action.App.CoreInfo.MasterAccountID),
	)
}

func (action *UserProofOfIncomeApproveAction) loadUser() {
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

func (action *UserProofOfIncomeApproveAction) performRequest() {
	document := action.User.Documents.Get(func(d *api.Document) bool {
		return d.Version == action.Version && d.Type.IsProofOfIncome()
	})
	if document == nil {
		action.SetInvalidField("document_key", errors.New("document not found"))
		return
	}

	document.Meta["reviewed"] = true
	document.Meta["approved"] = action.Request.Approved
	if !action.Request.Approved {
		document.Meta["reject_reason"] = action.Request.RejectReason
	} else {
		delete(document.Meta, "reject_reason")
	}

	err := action.APIQ().Users().Documents(action.User.DocumentsVersion).Set(action.User.ID, document)
	if err != nil {
		if err == api.ErrBadDocumentVersion {
			action.SetInvalidField("documents_version", err)
			return
		}
		action.Log.WithError(err).
			WithField("account", action.User.Address).
			WithField("document", document.Version).
			Error("failed to update document")
		action.Err = &problem.ServerError
		return
	}

	if action.Request.Approved {
		err := action.App.horizon.SubmitTX(action.Request.TX)
		if err != nil {
			entry := action.Log.WithError(err)

			if serr, ok := err.(horizon.SubmitError); ok {
				entry = entry.
					WithField("tx code", serr.TransactionCode()).
					WithField("op codes", serr.OperationCodes())

				for _, code := range serr.OperationCodes() {
					if code == "op_bad_auth" {
						action.Err = &problem.Forbidden
						return
					}
				}
			}

			entry.Error("failed to submit poi tx")
			action.Err = &problem.ServerError
			return
		}
	}
}
