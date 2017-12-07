package handlers

import (
	"context"
	"net/http"

	"gitlab.com/swarmfund/api/db2/api"
	"gitlab.com/swarmfund/api/internal/data"
	"gitlab.com/swarmfund/api/log"
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
)

func CtxWalletQ(q api.WalletQI) func(context.Context) context.Context {
	return func(ctx context.Context) context.Context {
		return context.WithValue(ctx, walletCtxKey, q)
	}
}

func WalletQ(r *http.Request) api.WalletQI {
	return r.Context().Value(walletCtxKey).(api.WalletQI).New()
}

func CtxLog(entry *log.Entry) func(context.Context) context.Context {
	return func(ctx context.Context) context.Context {
		return context.WithValue(ctx, logCtxKey, entry)
	}
}

func Log(r *http.Request) *log.Entry {
	return r.Context().Value(logCtxKey).(*log.Entry)
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
