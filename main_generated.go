package api

import (
	"github.com/zenazn/goji/web"
	"net/http"
)

// ServeHTTPC is a method for web.Handler
func (action ApproveRecoveryRequestAction) ServeHTTPC(c web.C, w http.ResponseWriter, r *http.Request) {
	ap := &action.Action
	ap.Prepare(c, w, r)
	action.Log = action.Log.WithField("action", "ApproveRecoveryRequestAction")
	ap.Execute(&action)
}

// ServeHTTPC is a method for web.Handler
func (action CreateRecoveryRequestAction) ServeHTTPC(c web.C, w http.ResponseWriter, r *http.Request) {
	ap := &action.Action
	ap.Prepare(c, w, r)
	action.Log = action.Log.WithField("action", "CreateRecoveryRequestAction")
	ap.Execute(&action)
}

// ServeHTTPC is a method for web.Handler
func (action DeleteContactShareAction) ServeHTTPC(c web.C, w http.ResponseWriter, r *http.Request) {
	ap := &action.Action
	ap.Prepare(c, w, r)
	action.Log = action.Log.WithField("action", "DeleteContactShareAction")
	ap.Execute(&action)
}

// ServeHTTPC is a method for web.Handler
func (action DetailsAction) ServeHTTPC(c web.C, w http.ResponseWriter, r *http.Request) {
	ap := &action.Action
	ap.Prepare(c, w, r)
	action.Log = action.Log.WithField("action", "DetailsAction")
	ap.Execute(&action)
}

// ServeHTTPC is a method for web.Handler
func (action GetEnumsAction) ServeHTTPC(c web.C, w http.ResponseWriter, r *http.Request) {
	ap := &action.Action
	ap.Prepare(c, w, r)
	action.Log = action.Log.WithField("action", "GetEnumsAction")
	ap.Execute(&action)
}

// ServeHTTPC is a method for web.Handler
func (action GetNotificationsAction) ServeHTTPC(c web.C, w http.ResponseWriter, r *http.Request) {
	ap := &action.Action
	ap.Prepare(c, w, r)
	action.Log = action.Log.WithField("action", "GetNotificationsAction")
	ap.Execute(&action)
}

// ServeHTTPC is a method for web.Handler
func (action GetRecoveryRequestAction) ServeHTTPC(c web.C, w http.ResponseWriter, r *http.Request) {
	ap := &action.Action
	ap.Prepare(c, w, r)
	action.Log = action.Log.WithField("action", "GetRecoveryRequestAction")
	ap.Execute(&action)
}

// ServeHTTPC is a method for web.Handler
func (action GetRecoveryRequestsAction) ServeHTTPC(c web.C, w http.ResponseWriter, r *http.Request) {
	ap := &action.Action
	ap.Prepare(c, w, r)
	action.Log = action.Log.WithField("action", "GetRecoveryRequestsAction")
	ap.Execute(&action)
}

// ServeHTTPC is a method for web.Handler
func (action GetUserDocsAction) ServeHTTPC(c web.C, w http.ResponseWriter, r *http.Request) {
	ap := &action.Action
	ap.Prepare(c, w, r)
	action.Log = action.Log.WithField("action", "GetUserDocsAction")
	ap.Execute(&action)
}

// ServeHTTPC is a method for web.Handler
func (action GetUserFileAction) ServeHTTPC(c web.C, w http.ResponseWriter, r *http.Request) {
	ap := &action.Action
	ap.Prepare(c, w, r)
	action.Log = action.Log.WithField("action", "GetUserFileAction")
	ap.Execute(&action)
}

// ServeHTTPC is a method for web.Handler
func (action GetUserIdAction) ServeHTTPC(c web.C, w http.ResponseWriter, r *http.Request) {
	ap := &action.Action
	ap.Prepare(c, w, r)
	action.Log = action.Log.WithField("action", "GetUserIdAction")
	ap.Execute(&action)
}

// ServeHTTPC is a method for web.Handler
func (action NotFoundAction) ServeHTTPC(c web.C, w http.ResponseWriter, r *http.Request) {
	ap := &action.Action
	ap.Prepare(c, w, r)
	action.Log = action.Log.WithField("action", "NotFoundAction")
	ap.Execute(&action)
}

// ServeHTTPC is a method for web.Handler
func (action ParticipantsAction) ServeHTTPC(c web.C, w http.ResponseWriter, r *http.Request) {
	ap := &action.Action
	ap.Prepare(c, w, r)
	action.Log = action.Log.WithField("action", "ParticipantsAction")
	ap.Execute(&action)
}

// ServeHTTPC is a method for web.Handler
func (action PatchNotificationsAction) ServeHTTPC(c web.C, w http.ResponseWriter, r *http.Request) {
	ap := &action.Action
	ap.Prepare(c, w, r)
	action.Log = action.Log.WithField("action", "PatchNotificationsAction")
	ap.Execute(&action)
}

// ServeHTTPC is a method for web.Handler
func (action PendingTransactionCreateAction) ServeHTTPC(c web.C, w http.ResponseWriter, r *http.Request) {
	ap := &action.Action
	ap.Prepare(c, w, r)
	action.Log = action.Log.WithField("action", "PendingTransactionCreateAction")
	ap.Execute(&action)
}

// ServeHTTPC is a method for web.Handler
func (action PendingTransactionDeleteAction) ServeHTTPC(c web.C, w http.ResponseWriter, r *http.Request) {
	ap := &action.Action
	ap.Prepare(c, w, r)
	action.Log = action.Log.WithField("action", "PendingTransactionDeleteAction")
	ap.Execute(&action)
}

// ServeHTTPC is a method for web.Handler
func (action PendingTransactionIndexAction) ServeHTTPC(c web.C, w http.ResponseWriter, r *http.Request) {
	ap := &action.Action
	ap.Prepare(c, w, r)
	action.Log = action.Log.WithField("action", "PendingTransactionIndexAction")
	ap.Execute(&action)
}

// ServeHTTPC is a method for web.Handler
func (action PendingTransactionRejectAction) ServeHTTPC(c web.C, w http.ResponseWriter, r *http.Request) {
	ap := &action.Action
	ap.Prepare(c, w, r)
	action.Log = action.Log.WithField("action", "PendingTransactionRejectAction")
	ap.Execute(&action)
}

// ServeHTTPC is a method for web.Handler
func (action PendingTransactionShowAction) ServeHTTPC(c web.C, w http.ResponseWriter, r *http.Request) {
	ap := &action.Action
	ap.Prepare(c, w, r)
	action.Log = action.Log.WithField("action", "PendingTransactionShowAction")
	ap.Execute(&action)
}

// ServeHTTPC is a method for web.Handler
func (action PutDocumentAction) ServeHTTPC(c web.C, w http.ResponseWriter, r *http.Request) {
	ap := &action.Action
	ap.Prepare(c, w, r)
	action.Log = action.Log.WithField("action", "PutDocumentAction")
	ap.Execute(&action)
}

// ServeHTTPC is a method for web.Handler
func (action ResolveUserRecoveryRequestAction) ServeHTTPC(c web.C, w http.ResponseWriter, r *http.Request) {
	ap := &action.Action
	ap.Prepare(c, w, r)
	action.Log = action.Log.WithField("action", "ResolveUserRecoveryRequestAction")
	ap.Execute(&action)
}

// ServeHTTPC is a method for web.Handler
func (action UserProofOfIncomeApproveAction) ServeHTTPC(c web.C, w http.ResponseWriter, r *http.Request) {
	ap := &action.Action
	ap.Prepare(c, w, r)
	action.Log = action.Log.WithField("action", "UserProofOfIncomeApproveAction")
	ap.Execute(&action)
}

// ServeHTTPC is a method for web.Handler
func (action UserProofOfIncomeIndexAction) ServeHTTPC(c web.C, w http.ResponseWriter, r *http.Request) {
	ap := &action.Action
	ap.Prepare(c, w, r)
	action.Log = action.Log.WithField("action", "UserProofOfIncomeIndexAction")
	ap.Execute(&action)
}
