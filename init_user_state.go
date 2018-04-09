package api

import (
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/swarmfund/api/db2/api"
	"gitlab.com/swarmfund/api/internal/hose"
	"gitlab.com/swarmfund/api/internal/types"
	"gitlab.com/swarmfund/go/xdr"
)

func checkState(change xdr.LedgerEntryChange) *api.UserStateUpdate {
	switch change.Type {
	case xdr.LedgerEntryChangeTypeCreated:
		entry := change.Created.Data
		switch entry.Type {
		case xdr.LedgerEntryTypeReviewableRequest:
			request := entry.ReviewableRequest
			switch entry.ReviewableRequest.Body.Type {
			case xdr.ReviewableRequestTypeUpdateKyc:
				state := types.UserStateWaitingForApproval
				return &api.UserStateUpdate{
					Address: types.Address(request.Requestor.Address()),
					State:   &state,
				}
			}
		}
	case xdr.LedgerEntryChangeTypeUpdated:
		entry := change.Updated.Data
		switch entry.Type {
		case xdr.LedgerEntryTypeReviewableRequest:
			request := entry.ReviewableRequest
			switch entry.ReviewableRequest.Body.Type {
			case xdr.ReviewableRequestTypeUpdateKyc:
				var state types.UserState
				if request.RejectReason != "" {
					// reject reason means KYC request were rejected
					state = types.UserStateRejected
				} else {
					// no reject reasons means request is pending
					state = types.UserStateWaitingForApproval
				}
				return &api.UserStateUpdate{
					Address: types.Address(request.Requestor.Address()),
					State:   &state,
				}
			}
		}
	case xdr.LedgerEntryChangeTypeRemoved:
		switch change.Removed.Type {
		case xdr.LedgerEntryTypeReviewableRequest:
			// TODO track requests state locally
		}
	}
	return nil
}

func checkKYC(change xdr.LedgerEntryChange) *api.UserStateUpdate {
	var kyc *xdr.AccountKycEntry
	switch change.Type {
	case xdr.LedgerEntryChangeTypeCreated:
		entry := change.Created.Data
		switch entry.Type {
		case xdr.LedgerEntryTypeAccountKyc:
			kyc = entry.AccountKyc
		}
	case xdr.LedgerEntryChangeTypeUpdated:
		entry := change.Updated.Data
		switch entry.Type {
		case xdr.LedgerEntryTypeAccountKyc:
			kyc = entry.AccountKyc
		}
	}
	if kyc != nil {
		return &api.UserStateUpdate{
			Address: types.Address(kyc.AccountId.Address()),
			KYCBlob: (*string)(&kyc.KycData),
		}
	}
	return nil
}

func checkType(change xdr.LedgerEntryChange) *api.UserStateUpdate {
	switch change.Type {
	case xdr.LedgerEntryChangeTypeUpdated:
		entry := change.Updated.Data
		switch entry.Type {
		case xdr.LedgerEntryTypeAccount:
			account := entry.Account
			var tpe types.UserType
			switch account.AccountType {
			case xdr.AccountTypeNotVerified:
				tpe = types.UserTypeNotVerified
			case xdr.AccountTypeGeneral:
				tpe = types.UserTypeGeneral
			case xdr.AccountTypeSyndicate:
				tpe = types.UserTypeSyndicate
			default:
				panic(errors.From(errors.New("unexpected account type"), logan.F{
					"account": account.AccountId.Address(),
				}))
			}
			return &api.UserStateUpdate{
				Address: types.Address(account.AccountId.Address()),
				Type:    &tpe,
			}
		}
	}
	return nil
}

func init() {
	appInit.Add("user-state-watcher", func(app *App) {
		entry := app.Config().Log().WithField("service", "user-state-watcher")
		mutators := []func(xdr.LedgerEntryChange) *api.UserStateUpdate{
			checkState, checkKYC, checkType,
		}
		app.txBus.Subscribe(func(event hose.TransactionEvent) {
			if event.Transaction == nil {
				return
			}
			for _, change := range event.Transaction.LedgerChanges() {
				for _, mutator := range mutators {
					// TODO defer on mutator call
					if update := mutator(change); update != nil {
						update.Timestamp = event.Transaction.CreatedAt
						err := app.apiQ.Users().SetState(*update)
						if err != nil {
							entry.WithError(err).Error("failed to set user state")
						}
					}
				}
			}
		})
	}, "tx-watcher")
}
