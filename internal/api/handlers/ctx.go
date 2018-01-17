package handlers

import (
	"context"
	"net/http"

	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/swarmfund/api/coreinfo"
	"gitlab.com/swarmfund/api/db2/api"
	"gitlab.com/swarmfund/api/internal/data"
	"gitlab.com/swarmfund/api/internal/hose"
	"gitlab.com/swarmfund/api/notificator"
	"gitlab.com/swarmfund/api/storage"
	"gitlab.com/swarmfund/go/doorman"
	"gitlab.com/swarmfund/go/keypair"
	"gitlab.com/swarmfund/horizon-connector"
)

type ctxKey int

const (
	logCtxKey ctxKey = iota
	walletCtxKey
	emailTokensQCtxKey
	usersQCtxKey
	horizonConnectorCtxKey
	accountManagerKPCtxKey
	tfaQCtxKey
	doormanCtxKey
	storageCtxKey
	coreInfoCtxKey
	blobQCtxKey
	userBusDispatchCtxKey
	authorizedDeviceCtxKey
	notificatorCtxKey
)

func CtxWalletQ(q api.WalletQI) func(context.Context) context.Context {
	return func(ctx context.Context) context.Context {
		return context.WithValue(ctx, walletCtxKey, q)
	}
}

func WalletQ(r *http.Request) api.WalletQI {
	return r.Context().Value(walletCtxKey).(api.WalletQI).New()
}

func CtxAuthorizedDeviceQ(q api.AuthorizedDeviceQI) func(context.Context) context.Context {
	return func(ctx context.Context) context.Context {
		return context.WithValue(ctx, authorizedDeviceCtxKey, q)
	}
}

func AuthorizedDeviceQ(r *http.Request) api.AuthorizedDeviceQI {
	return r.Context().Value(walletCtxKey).(api.AuthorizedDeviceQI).New()
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

func CtxAccountManagerKP(kp keypair.KP) func(context.Context) context.Context {
	return func(ctx context.Context) context.Context {
		return context.WithValue(ctx, accountManagerKPCtxKey, kp)
	}
}

func AccountManagerKP(r *http.Request) keypair.KP {
	return r.Context().Value(accountManagerKPCtxKey).(keypair.KP)
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

func CtxCoreInfo(s *coreinfo.Connector) func(context.Context) context.Context {
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

func CtxNotificator(conn notificator.ConnectorI) func(context.Context) context.Context {
	return func(ctx context.Context) context.Context {
		return context.WithValue(ctx, notificatorCtxKey, conn)
	}
}

func Notificator(r *http.Request) notificator.ConnectorI {
	return r.Context().Value(notificatorCtxKey).(notificator.ConnectorI)
}
