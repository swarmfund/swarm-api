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
	// Detached comes from join on organization_wallets and shows if wallet is
	// independent or used in organization flow
	Detached bool `db:"detached"`
	// OrganizationAddress comes from join on organization_wallets and is a
	// foreign key to users.address. If wallet has non nil OrganizationAddress
	// means it's connected to one and his signature should be checked against it
	OrganizationAddress *string `db:"organization_address"`
}
