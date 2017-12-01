package api

import (
	"net/http"

	"gitlab.com/swarmfund/api/internal/api/handlers"
	"gitlab.com/swarmfund/api/internal/secondfactor"
	"gitlab.com/swarmfund/api/pentxsub"
	"gitlab.com/swarmfund/api/render/problem"
	"gitlab.com/swarmfund/go/xdr"
	"gitlab.com/swarmfund/horizon-connector"
)

// PendingTransactionCreateAction submits a transaction to the stellar-core network
// on behalf of the requesting client.
type PendingTransactionCreateAction struct {
	Action

	TX string

	Envelope xdr.TransactionEnvelope
	Source   string

	Rendered bool
	Resource interface{}
}

// JSON format action handler
func (action *PendingTransactionCreateAction) JSON() {
	action.Do(
		action.ValidateBodyType,
		action.loadTX,
		action.marshalEnvelope,
		action.checkSource,
		action.checkTFA,
		action.submitTX,
	)
}

func (action *PendingTransactionCreateAction) loadTX() {
	action.TX = action.GetNonEmptyString("tx")
}

func (action *PendingTransactionCreateAction) marshalEnvelope() {
	err := xdr.SafeUnmarshalBase64(action.TX, &action.Envelope)
	if err != nil {
		action.Err = action.malformedP()
		return
	}
}

func (action *PendingTransactionCreateAction) checkSource() {
	action.Source = action.Envelope.Tx.SourceAccount.Address()
	for _, operation := range action.Envelope.Tx.Operations {
		opSource := operation.SourceAccount
		if opSource != nil && opSource.Address() != action.Source {
			// transaction has messed up sources and I'm not in the mood to handle it.
			// also most probably no one will ever use it with good intentions.
			action.Err = action.malformedP()
			return
		}
	}
}

func (action *PendingTransactionCreateAction) checkTFA() {
	user, err := action.APIQ().Users().ByAddress(action.Source)
	if err != nil {
		action.Log.WithError(err).Error("failed to get user")
		action.Err = &problem.ServerError
		return
	}

	if user == nil {
		return
	}

	for _, operation := range action.Envelope.Tx.Operations {
		switch operation.Body.Type {
		case xdr.OperationTypeManageOffer, xdr.OperationTypePayment, xdr.OperationTypeManageForfeitRequest:
			signer, err := action.App.pendingSubmitter.GetSigner(action.TX)
			if err != nil {
				action.Log.WithError(err).Error("failed to get tx signer")
				action.Err = &problem.ServerError
				return
			}
			if signer == nil {
				// nothing we can do about it here
				// so just passing it along
				return
			}
			wallet, err := action.APIQ().Wallet().ByCurrentAccountID(signer.AccountID)
			if err != nil {
				action.Log.WithError(err).Error("failed to get wallet")
				action.Err = &problem.ServerError
				return
			}

			if wallet == nil {
				return
			}

			if err := secondfactor.NewConsumer(action.APIQ().TFA()).Consume(action.R, wallet); err != nil {
				handlers.RenderFactorConsumeError(action.W, action.R, err)
				action.Rendered = true
				return
			}
		}
	}
}

func (action *PendingTransactionCreateAction) submitTX() {
	if action.Rendered {
		return
	}
	result, err := action.App.pendingSubmitter.Submit(action.TX)
	if err != nil {
		switch err {
		case pentxsub.ErrTooManyOps:
			action.Err = &problem.AdminOperationsRestrictionViolated
			return
		default:
			serr, ok := err.(horizon.SubmitError)
			if ok {
				action.W.WriteHeader(serr.ResponseCode())
				action.W.Write(serr.ResponseBody())
				return
			}
			action.Log.WithError(err).Error("failed to submit tx")
			action.Err = &problem.ServerError
			return
		}
	}
	action.W.Write(result)
}

func (action *PendingTransactionCreateAction) malformedP() *problem.P {
	return &problem.P{
		Type:   "transaction_malformed",
		Title:  "Transaction Malformed",
		Status: http.StatusBadRequest,
		Detail: "API thinks your transaction is invalid in some way",
		Extras: map[string]interface{}{
			"envelope_xdr": action.TX,
		},
	}
}
