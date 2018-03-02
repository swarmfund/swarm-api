package data

type KDF struct {
	Version   int     `db:"version"`
	Algorithm string  `db:"algorithm"`
	Bits      uint    `db:"bits"`
	N         float64 `db:"n"`
	R         uint    `db:"r"`
	P         uint    `db:"p"`
	Salt      string  `db:"salt"`
}

type WalletKDF struct {
	// Wallet foreign key on wallet email
	Wallet string `db:"wallet"`
	// Version foreign key on kdf version
	Version int `db:"version"`
	// Salt used to encrypt keychain data
	Salt string `db:"salt"`
}
