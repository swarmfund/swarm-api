package data

//go:generate mockery -case underscore -name CoreInfoI

type CoreInfoI interface {
	GetMasterAccountID() string
	GetPassphrase() string
	GetTXExpire() int64
}
