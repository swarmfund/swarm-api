package handlers

import (
	"context"
	"net/http"

	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/swarmfund/api/blacklist"
	"gitlab.com/swarmfund/api/config"
	"gitlab.com/swarmfund/api/db2/api"
	"gitlab.com/swarmfund/api/internal/data"
	"gitlab.com/swarmfund/api/internal/hose"
	"gitlab.com/swarmfund/api/internal/salesforce"
	"gitlab.com/swarmfund/api/notificator"
	"gitlab.com/tokend/go/doorman"
	"gitlab.com/tokend/go/xdrbuild"
	"gitlab.com/tokend/horizon-connector"
)

type ctxKey int

const (
	logCtxKey ctxKey = iota
	walletCtxKey
	emailTokensQCtxKey
	usersQCtxKey
	horizonConnectorCtxKey
	txSignerCtxKey
	txBuilderCtxKey
	tfaQCtxKey
	doormanCtxKey
	storageCtxKey
	salesforceCtxKey
	coreInfoCtxKey
	blobQCtxKey
	notificatorCtxKey
	walletAdditionCtxKey
	domainApproverCtxKey
	logBusDispatchCtxKey
	auditLogsQCtxKey
)

func CtxWalletQ(q api.WalletQI) func(context.Context) context.Context {
	return func(ctx context.Context) context.Context {
		return context.WithValue(ctx, walletCtxKey, q)
	}
}

func WalletQ(r *http.Request) api.WalletQI {
	return r.Context().Value(walletCtxKey).(api.WalletQI).New()
}

func CtxLog(entry *logan.Entry) func(context.Context) context.Context {
	return func(ctx context.Context) context.Context {
		return context.WithValue(ctx, logCtxKey, entry)
	}
}

func Log(r *http.Request) *logan.Entry {
	return r.Context().Value(logCtxKey).(*logan.Entry)
}

func CtxEmailTokensQ(q data.EmailTokensQ) func(context.Context) context.Context {
	return func(ctx context.Context) context.Context {
		return context.WithValue(ctx, emailTokensQCtxKey, q)
	}
}

func EmailTokensQ(r *http.Request) data.EmailTokensQ {
	return r.Context().Value(emailTokensQCtxKey).(data.EmailTokensQ).New()
}

func CtxUsersQ(q api.UsersQI) func(context.Context) context.Context {
	return func(ctx context.Context) context.Context {
		return context.WithValue(ctx, usersQCtxKey, q)
	}
}

func UsersQ(r *http.Request) api.UsersQI {
	return r.Context().Value(usersQCtxKey).(api.UsersQI).New()
}

func CtxHorizon(q *horizon.Connector) func(context.Context) context.Context {
	return func(ctx context.Context) context.Context {
		return context.WithValue(ctx, horizonConnectorCtxKey, q)
	}
}

func Horizon(r *http.Request) *horizon.Connector {
	return r.Context().Value(horizonConnectorCtxKey).(*horizon.Connector)
}

func CtxTFAQ(q api.TFAQI) func(context.Context) context.Context {
	return func(ctx context.Context) context.Context {
		return context.WithValue(ctx, tfaQCtxKey, q)
	}
}

func TFAQ(r *http.Request) api.TFAQI {
	return r.Context().Value(tfaQCtxKey).(api.TFAQI).New()
}

func CtxDoorman(d doorman.Doorman) func(context.Context) context.Context {
	return func(ctx context.Context) context.Context {
		return context.WithValue(ctx, doormanCtxKey, d)
	}
}

func Doorman(r *http.Request, constraints ...doorman.SignerConstraint) error {
	d := r.Context().Value(doormanCtxKey).(doorman.Doorman)
	return d.Check(r, constraints...)
}

func CtxStorage(s data.Storage) func(context.Context) context.Context {
	return func(ctx context.Context) context.Context {
		return context.WithValue(ctx, storageCtxKey, s)
	}
}

func Storage(r *http.Request) data.Storage {
	return r.Context().Value(storageCtxKey).(data.Storage)
}

func CtxSalesforce(s *salesforce.Connector) func(context.Context) context.Context {
	return func(ctx context.Context) context.Context {
		return context.WithValue(ctx, salesforceCtxKey, s)
	}
}

func Salesforce(r *http.Request) *salesforce.Connector {
	return r.Context().Value(salesforceCtxKey).(*salesforce.Connector)
}

func CtxCoreInfo(s data.Info) func(context.Context) context.Context {
	return func(ctx context.Context) context.Context {
		return context.WithValue(ctx, coreInfoCtxKey, s)
	}
}

func CoreInfo(r *http.Request) *horizon.Info {
	info, err := r.Context().Value(coreInfoCtxKey).(data.Info).Info()
	if err != nil {
		//TODO handle error
		panic(err)
	}
	return info
}

func CtxBlobQ(q data.Blobs) func(context.Context) context.Context {
	return func(ctx context.Context) context.Context {
		return context.WithValue(ctx, blobQCtxKey, q)
	}
}

func BlobQ(r *http.Request) data.Blobs {
	return r.Context().Value(blobQCtxKey).(data.Blobs).New()
}

func CtxLogBusDispatch(dispatch hose.LogDispatch) func(context.Context) context.Context {
	return func(ctx context.Context) context.Context {
		return context.WithValue(ctx, logBusDispatchCtxKey, dispatch)
	}
}

func LogBusDispatch(r *http.Request, event hose.LogEvent) {
	dispatch := r.Context().Value(logBusDispatchCtxKey).(hose.LogDispatch)
	dispatch(event)
}

func CtxNotificator(notificator *notificator.Connector) func(context.Context) context.Context {
	return func(ctx context.Context) context.Context {
		return context.WithValue(ctx, notificatorCtxKey, notificator)
	}
}

func Notificator(r *http.Request) *notificator.Connector {
	return r.Context().Value(notificatorCtxKey).(*notificator.Connector)
}

func CtxTransaction(txbuilder data.Infobuilder) func(context.Context) context.Context {
	return func(ctx context.Context) context.Context {
		ctx = context.WithValue(ctx, txBuilderCtxKey, txbuilder)
		return ctx
	}

}

func Transaction(r *http.Request) *xdrbuild.Transaction {
	txbuilderbuilder := r.Context().Value(txBuilderCtxKey).(data.Infobuilder)
	info := r.Context().Value(coreInfoCtxKey).(data.Info)
	return txbuilderbuilder(info)
}

func CtxWallets(disableConfirm config.Wallets) func(context.Context) context.Context {
	return func(ctx context.Context) context.Context {
		return context.WithValue(ctx, walletAdditionCtxKey, disableConfirm)
	}
}

func Wallet(r *http.Request) config.Wallets {
	return r.Context().Value(walletAdditionCtxKey).(config.Wallets)
}

func CtxDomainApprover(approver *blacklist.Approver) func(context.Context) context.Context {
	return func(ctx context.Context) context.Context {
		return context.WithValue(ctx, domainApproverCtxKey, approver)
	}
}

func DomainApprover(r *http.Request) *blacklist.Approver {
	return r.Context().Value(domainApproverCtxKey).(*blacklist.Approver)
}

func CtxAuditLogs(q api.AuditLogQI) func(context.Context) context.Context {
	return func(ctx context.Context) context.Context {
		return context.WithValue(ctx, auditLogsQCtxKey, q)
	}
}

func AuditLogs(r *http.Request) api.AuditLogQI {
	return r.Context().Value(auditLogsQCtxKey).(api.AuditLogQI).New()
}
