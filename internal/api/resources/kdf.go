package resources

type KDFVersion struct {
	Version int `jsonapi:"primary,kdf"`
}

type KDF struct {
	Version   int     `jsonapi:"primary,kdf"`
	Algorithm string  `jsonapi:"attr,algorithm"`
	Bits      uint    `jsonapi:"attr,bits"`
	N         float64 `jsonapi:"attr,n"`
	R         uint    `jsonapi:"attr,r"`
	P         uint    `jsonapi:"attr,p"`
	Salt      string  `jsonapi:"attr,salt,omitempty"`
}
