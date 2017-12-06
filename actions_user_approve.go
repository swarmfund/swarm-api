package api

import (
	"encoding/json"

	"github.com/go-errors/errors"
	"gitlab.com/swarmfund/api/db2/api"
	"gitlab.com/swarmfund/api/render/hal"
	"gitlab.com/swarmfund/api/render/problem"
	"gitlab.com/swarmfund/go/xdr"
	"gitlab.com/swarmfund/horizon-connector"
)

type UserApproveRequest struct {
	// Approved default is false, so it's ok if field is missing
	Approved bool `json:"approved"`
	// RejectReason is conditionally required, should check down the processing pipe
	RejectReasons    json.RawMessage            `json:"reject_reasons"`
	Documents        api.DocumentsRejectReasons `json:"documents"`
	DocumentsVersion int64                      `json:"documents_version" valid:"required"`
}

type UserApproveAction struct {
	Action

	Request   UserApproveRequest
	AccountID string

	User *api.User
}

func (action *UserApproveAction) JSON() {
	action.Do(
		action.ValidateBodyType,
		action.loadParams,
		action.checkAllowed,
		action.loadUser,
		action.checkUserState,
		action.performRequest,
		func() {
			hal.Render(action.W, problem.Success)
		})
}

func (action *UserApproveAction) loadParams() {
	action.UnmarshalBody(&action.Request)
	action.AccountID = action.GetNonEmptyString("user")
}

func (action *UserApproveAction) checkAllowed() {
	action.checkSignerConstraints(
		SignerOf(action.App.CoreInfo.MasterAccountID),
	)
}

func (action *UserApproveAction) loadUser() {
	user, err := action.APIQ().Users().ByAddress(action.AccountID)
	if err != nil {
		action.Log.WithError(err).Error("failed to load user")
		action.Err = &problem.ServerError
		return
	}

	if user == nil {
		action.Err = &problem.NotFound
		return
	}

	action.User = user
}

func (action *UserApproveAction) checkUserState() {
	//if action.User.State != api.UserWaitingForApproval {
	//	action.Err = &problem.Forbidden
	//	return
	//}
}

func (action *UserApproveAction) performRequest() {
	//var err error
	// TODO implement
	// mark latest KYC docs, remove other. not writing anything until checks bellow passed
	//var docTypes map[string]bool
	//switch action.User.UserType {

	//case api.UserTypeIndividual:
	//	docTypes = api.IndividualDocTypes
	//case api.UserTypeBusiness:
	//	docTypes = api.NonIndividualDocTypes
	//default:
	//	action.Err = &problem.BadRequest
	//	return
	//}

	// TODO not sure it works
	docsToDelete := []string{}
	//for docType, _ := range docTypes {
	//	documents := action.User.Documents[docType]
	//	latestDoc := action.User.Documents.Latest(docType)
	//	if latestDoc == nil || documents == nil {
	//		// something went wrong
	//		// user don't have one of the required docs upload
	//		action.Err = &problem.BadRequest
	//		return
	//	}
	//
	//	if latestDoc.Meta == nil {
	//		latestDoc.Meta = map[string]interface{}{}
	//	}
	//	latestDoc.Meta["approved"] = action.Request.Approved
	//
	//	for _, document := range documents {
	//		if document.Key != latestDoc.Key {
	//			if document.Meta == nil {
	//				document.Meta = map[string]interface{}{}
	//			}
	//			document.Meta["approved"] = false
	//			docsToDelete = append(docsToDelete, document.Key)
	//		}
	//	}
	//}
	if action.Request.Approved {
		action.approveUser()
	} else {
		action.rejectUser()
	}
	if action.Err != nil {
		return
	}

	for _, document := range docsToDelete {
		ohaigo := document
		go func() {
			if err := action.App.Storage().Delete(string(action.User.Address), ohaigo); err != nil {
				action.Log.WithError(err).WithField("user", action.User.Address).WithField("doc", document).Error("failed to delete doc")
			}
		}()
	}
}

func (action *UserApproveAction) approveUser() {
	action.User.DocumentsVersion = action.Request.DocumentsVersion
	if err := action.APIQ().Users().Approve(action.User); err != nil {
		if err == api.ErrBadDocumentVersion {
			action.SetInvalidField("documents_version", api.ErrBadDocumentVersion)
			return
		}
		action.Log.WithError(err).Error("failed to approve")
		action.Err = &problem.ServerError
		return
	}

	transaction := action.App.horizon.Transaction(&horizon.TransactionBuilder{Source: action.App.MasterKP()})

	account, err := action.App.horizon.AccountSigned(action.App.AccountManagerKP(), string(action.User.Address))
	if err != nil {
		action.Log.WithError(err).Error("failed to get core account")
		action.Err = &problem.ServerError
		return
	}

	if account == nil {
		action.SetInvalidField("address", errors.New("core account does not exists"))
		return
	}

	ops := []horizon.Operation{}
	if xdr.BlockReasons(account.BlockReasons)&xdr.BlockReasonsKycUpdate != 0 {
		ops = append(ops, &horizon.ManageAccountOp{
			AccountID:     account.AccountID,
			AccountType:   xdr.AccountType(account.AccountType),
			RemoveReasons: xdr.BlockReasonsKycUpdate,
		})
	}

	if xdr.AccountType(account.AccountType) == xdr.AccountTypeNotVerified {
		ops = append(ops, &horizon.CreateAccountOp{
			AccountID:   account.AccountID,
			AccountType: xdr.AccountTypeGeneral,
		})
	}

	if len(ops) > 0 {
		for _, op := range ops {
			transaction = transaction.Op(op)
		}
		err = transaction.Sign(action.App.AccountManagerKP()).Submit()
		if err != nil {
			entry := action.Log.WithError(err)
			if serr, ok := err.(horizon.SubmitError); ok {
				entry = entry.
					WithField("tx code", serr.TransactionCode()).
					WithField("op codes", serr.OperationCodes())
			}
			entry.Error("failed to submit approve user tx")
			action.Err = &problem.ServerError
			return
		}
	}

	err = action.Notificator().NotifyApproval(action.User.Email)
	if err != nil {
		action.Log.WithError(err).Error("Emails sending failed")
	}
}

func (action *UserApproveAction) rejectUser() {
	if action.Err != nil {
		return
	}

	if action.Request.RejectReasons == nil {
		action.SetInvalidField("reject_reasons", errors.New("required"))
		return
	}

	rr, err := action.User.ValidateRejectReasons(action.Request.RejectReasons)
	if err != nil {
		action.SetInvalidField("reject_reasons", errors.New("invalid in some way"))
		return
	}

	var rrEntityType api.KYCEntityType
	switch action.User.UserType {
	//case api.UserTypeIndividual:
	//	rrEntityType = api.KYCEntityTypeIndividualRejectReasons
	//case api.UserTypeJoint:
	//	rrEntityType = api.KYCEntityTypeJointRejectReasons
	//case api.UserTypeBusiness:
	//	rrEntityType = api.KYCEntityTypeBusinessRejectReasons
	default:
		panic("unknown user type")
	}

	rrEntity := action.User.KYCEntities.GetSingle(rrEntityType)
	if rrEntity != nil {
		// update rr entity
		err := action.APIQ().Users().KYC().Update(rrEntity.ID, rr)
		if err != nil {
			action.Log.WithError(err).Error("failed to update entity")
			action.Err = &problem.ServerError
			return
		}
	} else {
		_, err = action.APIQ().Users().KYC().Create(api.KYCEntity{
			Type:   rrEntityType,
			UserID: action.User.ID,
			Data:   rr,
		})
		if err != nil {
			action.Log.WithError(err).Error("failed to save entity")
			action.Err = &problem.ServerError
			return
		}
	}

	// documents reject
	if action.Request.Documents != nil {
		data, err := json.Marshal(action.Request.Documents)
		if err != nil {
			action.Log.WithError(err).Error("failed to marshal documents reject reasons")
			action.Err = &problem.ServerError
			return
		}
		rrEntity = action.User.KYCEntities.GetSingle(api.KYCEntityTypeDocumentsRejectReasons)
		if rrEntity != nil {
			// update rr entity
			err := action.APIQ().Users().KYC().Update(rrEntity.ID, data)
			if err != nil {
				action.Log.WithError(err).Error("failed to update entity")
				action.Err = &problem.ServerError
				return
			}
		} else {
			_, err = action.APIQ().Users().KYC().Create(api.KYCEntity{
				Type:   api.KYCEntityTypeDocumentsRejectReasons,
				UserID: action.User.ID,
				Data:   data,
			})
			if err != nil {
				action.Log.WithError(err).Error("failed to save entity")
				action.Err = &problem.ServerError
				return
			}
		}
	}

	if err := action.APIQ().Users().Reject(action.User); err != nil {
		if err == api.ErrBadDocumentVersion {
			action.SetInvalidField("documents_version", api.ErrBadDocumentVersion)
			return
		}
		action.Log.WithError(err).Error("failed to reject user")
		action.Err = &problem.ServerError
		return
	}

	if err := action.Notificator().NotifyRejection(action.User.Email); err != nil {
		action.Log.WithError(err).Error("Emails sending failed")
	}
}
