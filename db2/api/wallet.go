package api

type Wallet struct {
	Id                int64  `db:"id"`
	AccountID         string `db:"account_id"`
	CurrentAccountID  string `db:"current_account_id"`
	WalletId          string `db:"wallet_id"`
	Username          string `db:"email"`
	Salt              string `db:"salt"`
	KDF               int    `db:"kdf_id"`
	KeychainData      string `db:"keychain_data"`
	VerificationToken string `db:"verification_token"`
	// Verified comes from join on email_tokens and shows if wallet email was confirmed
	Verified bool `db:"verified"`
}
