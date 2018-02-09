package api

import "gitlab.com/swarmfund/api/internal/types"

type Wallet struct {
	Id                int64         `db:"id"`
	AccountID         types.Address `db:"account_id"`
	CurrentAccountID  types.Address `db:"current_account_id"`
	WalletId          string        `db:"wallet_id"`
	Username          string        `db:"email"`
	Salt              string        `db:"salt"`
	KDF               int           `db:"kdf_id"`
	KeychainData      string        `db:"keychain_data"`
	VerificationToken string        `db:"verification_token"`
	// Verified comes from join on email_tokens and shows if wallet email was confirmed
	Verified bool `db:"verified"`
	// RecoveryAddress account recovery key, comes from join on recoveries
	RecoveryAddress *types.Address `db:"recovery_address"`
	// RecoveryWalletID comes from join on recoveries
	RecoveryWalletID *string `db:"recovery_wallet_id"`
	RecoverySalt     *string `db:"recovery_salt"`
}
