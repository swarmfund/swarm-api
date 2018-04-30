package api

import (
	"encoding/json"

	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/swarmfund/api/db2/api"
	"gitlab.com/swarmfund/api/internal/hose"
	"gitlab.com/swarmfund/api/internal/types"
	"gitlab.com/tokend/go/xdr"
)

type KYCUpdateRequest struct {
	Requestor types.Address
}

type StateChecker struct {
	requests map[xdr.Uint64]KYCUpdateRequest
}

func NewCheckState() func(change xdr.LedgerEntryChange) *api.UserStateUpdate {
	s := StateChecker{
		requests: map[xdr.Uint64]KYCUpdateRequest{},
	}
	return s.checkState
}

func (s *StateChecker) checkState(change xdr.LedgerEntryChange) *api.UserStateUpdate {
	switch change.Type {
	case xdr.LedgerEntryChangeTypeCreated:
		entry := change.Created.Data
		switch entry.Type {
		case xdr.LedgerEntryTypeReviewableRequest:
			request := entry.ReviewableRequest
			switch entry.ReviewableRequest.Body.Type {
			case xdr.ReviewableRequestTypeUpdateKyc:
				// track kyc update requests once created,
				// so we could determine it's details once it deleted
				s.requests[request.RequestId] = KYCUpdateRequest{
					Requestor: types.Address(request.Requestor.Address()),
				}
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
			// relying on request tracking on create entry
			request, ok := s.requests[change.Removed.ReviewableRequest.RequestId]
			if ok {
				return &api.UserStateUpdate{
					Address: request.Requestor,
					State:   &types.DefaultUserState,
				}
			}
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
		var kycdata struct {
			BlobID string `json:"blob_id"`
		}
		if err := json.Unmarshal([]byte(kyc.KycData), &kycdata); err != nil {
			panic(errors.Wrap(err, "failed to unmarshal KYC data"))
		}
		return &api.UserStateUpdate{
			Address: types.Address(kyc.AccountId.Address()),
			KYCBlob: &kycdata.BlobID,
		}
	}
	return nil
}

func checkType(change xdr.LedgerEntryChange) *api.UserStateUpdate {
	var account *xdr.AccountEntry
	switch change.Type {
	case xdr.LedgerEntryChangeTypeUpdated:
		entry := change.Updated.Data
		switch entry.Type {
		case xdr.LedgerEntryTypeAccount:
			account = entry.Account
		}
	case xdr.LedgerEntryChangeTypeCreated:
		entry := change.Created.Data
		switch entry.Type {
		case xdr.LedgerEntryTypeAccount:
			account = entry.Account
		}
	}
	if account != nil {
		var tpe types.UserType
		switch account.AccountType {
		case xdr.AccountTypeNotVerified:
			tpe = types.UserTypeNotVerified
		case xdr.AccountTypeGeneral:
			tpe = types.UserTypeGeneral
		case xdr.AccountTypeSyndicate:
			tpe = types.UserTypeSyndicate
		case xdr.AccountTypeCommission, xdr.AccountTypeExchange, xdr.AccountTypeMaster, xdr.AccountTypeOperational:
			return nil
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
	return nil
}

func init() {
	appInit.Add("user-state-watcher", func(app *App) {
		entry := app.Config().Log().WithField("service", "user-state-watcher")
		mutators := []func(xdr.LedgerEntryChange) *api.UserStateUpdate{
			NewCheckState(), checkKYC, checkType,
		}
		app.txBus.Subscribe(func(event hose.TransactionEvent) {
			if event.Transaction == nil {
				return
			}
			for _, change := range event.Transaction.LedgerChanges() {
				for _, mutator := range mutators {
					update := func() *api.UserStateUpdate {
						defer func() {
							if rvr := recover(); rvr != nil {
								entry.WithRecover(rvr).WithFields(event.Transaction.GetLoganFields()).Error("mutator panicked")
							}
						}()
						return mutator(change)
					}()
					if update != nil {
						update.Timestamp = event.Transaction.CreatedAt
						err := app.apiQ.Users().SetState(*update)
						if err != nil {
							entry.WithError(err).WithFields(event.Transaction.GetLoganFields()).Error("failed to set user state")
						}
					}
				}
			}
		})
	}, "tx-watcher")
}
