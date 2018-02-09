package data

//go:generate mockery -case underscore -name CoreInfoI

type CoreInfoI interface {
	GetMasterAccountID() string
	Passphrase() string
	TXExpire() int64
}
