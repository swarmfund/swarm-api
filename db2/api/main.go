package api

import "gitlab.com/swarmfund/api/db2"

type QInterface interface {
	GetRepo() *db2.Repo

	Users() UsersQI
	Recoveries() RecoveryQI
	HMAC() HMACQI
	AuthorizedDevice() AuthorizedDeviceQI
	TFA() TFAQI
	Wallet() WalletQI

	PendingTransactions() PendingTransactionsQI
	PendingTransactionByID(dest interface{}, id int64) error
	PendingTransactionByHash(hash string) (*PendingTransaction, error)
	PendingTransactionSigners() PendingTransactionSignersQI
	PenTXSub() PenTXSubQI

	//KYCTracker() KYCTrackerQI

	Notifications() NotificationsQI
}

// Q is a helper struct on which to hang common queries against a history
// portion of the horizon database.
type Q struct {
	*db2.Repo
}

func (q *Q) Wallet() WalletQI {
	return &WalletQ{
		parent: q,
		sql:    walletSelect,
	}
}

func (q *Q) GetRepo() *db2.Repo {
	return q.Repo
}
