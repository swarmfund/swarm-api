package handlers

import (
	"context"
	"net/http"

	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/swarmfund/api/db2/api"
	"gitlab.com/swarmfund/api/internal/data"
	"gitlab.com/swarmfund/api/internal/hose"
	"gitlab.com/swarmfund/api/storage"
	"gitlab.com/swarmfund/go/doorman"
	"gitlab.com/swarmfund/go/xdrbuild"
	"gitlab.com/swarmfund/horizon-connector"
	"gitlab.com/tokend/keypair"
)

type ctxKey int

const (
	logCtxKey ctxKey = iota
	walletCtxKey
	emailTokensQCtxKey
	usersQCtxKey
	horizonConnectorCtxKey
	txSignerCtxKey
	txSourceCtxKey
	tfaQCtxKey
	doormanCtxKey
	storageCtxKey
	coreInfoCtxKey
	blobQCtxKey
	userBusDispatchCtxKey
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

func CtxStorage(s *storage.Connector) func(context.Context) context.Context {
	return func(ctx context.Context) context.Context {
		return context.WithValue(ctx, storageCtxKey, s)
	}
}

func Storage(r *http.Request) *storage.Connector {
	return r.Context().Value(storageCtxKey).(*storage.Connector)
}

func CtxCoreInfo(s data.CoreInfoI) func(context.Context) context.Context {
	return func(ctx context.Context) context.Context {
		return context.WithValue(ctx, coreInfoCtxKey, s)
	}
}

func CoreInfo(r *http.Request) data.CoreInfoI {
	return r.Context().Value(coreInfoCtxKey).(data.CoreInfoI)
}

func CtxBlobQ(q data.Blobs) func(context.Context) context.Context {
	return func(ctx context.Context) context.Context {
		return context.WithValue(ctx, blobQCtxKey, q)
	}
}

func BlobQ(r *http.Request) data.Blobs {
	return r.Context().Value(blobQCtxKey).(data.Blobs)
}

func CtxUserBusDispatch(dispatch hose.UserDispatch) func(context.Context) context.Context {
	return func(ctx context.Context) context.Context {
		return context.WithValue(ctx, userBusDispatchCtxKey, dispatch)
	}
}

func UserBusDispatch(r *http.Request, event hose.UserEvent) {
	dispatch := r.Context().Value(userBusDispatchCtxKey).(hose.UserDispatch)
	dispatch(event)
}

func CtxTransaction(source keypair.Address, signer keypair.Full) func(context.Context) context.Context {
	return func(ctx context.Context) context.Context {
		ctx = context.WithValue(ctx, txSourceCtxKey, source)
		ctx = context.WithValue(ctx, txSignerCtxKey, signer)
		return ctx
	}
}

func Transaction(r *http.Request) *xdrbuild.Transaction {
	info := CoreInfo(r)
	source := r.Context().Value(txSourceCtxKey).(keypair.Address)
	signer := r.Context().Value(txSourceCtxKey).(keypair.Full)
	return xdrbuild.
		NewBuilder(info.Passphrase(), info.TXExpire()).
		Transaction(source).
		Sign(signer)
}
