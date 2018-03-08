package data

import "time"

type EmailToken struct {
	ID         int64
	Token      string
	Email      string
	WalletID   string `db:"wallet_id"`
	Confirmed  bool
	LastSentAt *time.Time `db:"last_sent_at"`
}

//go:generate mockery -case underscore -name EmailTokensQ
type EmailTokensQ interface {
	New() EmailTokensQ
	// handlers:
	Create(walletID, token string, confirmed bool) error
	// Verify marks token as verified if provided arguments are valid.
	// bool shows if operation succeeded
	Verify(walletID, token string) (bool, error)
	Get(walletID string) (*EmailToken, error)
	MarkUnsent(tid int64) error
	// runner:
	GetUnsent() ([]EmailToken, error)
	GetUnconfirmed() ([]EmailToken, error)
	MarkSent(tid int64) error
}
