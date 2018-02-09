package resources

import "github.com/go-ozzo/ozzo-validation"

type KDFVersion struct {
	Version int `jsonapi:"primary,kdf"`
}

type KDFPlain struct {
	Data KDFPlainData `json:"data"`
}

func (r KDFPlain) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.Data, validation.Required),
	)
}

type KDFPlainData struct {
	Type string `json:"type"`
	ID   int    `json:"id,string"`
}

func (r KDFPlainData) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.Type, validation.In("kdf")),
		validation.Field(&r.ID, validation.Required),
	)
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
