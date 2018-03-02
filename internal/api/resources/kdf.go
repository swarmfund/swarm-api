package resources

import (
	"github.com/go-ozzo/ozzo-validation"
	"gitlab.com/swarmfund/api/internal/data"
)

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
	Type       string        `json:"type"`
	ID         int           `json:"id,string"`
	Attributes KDFAttributes `json:"attributes"`
}

type KDFAttributes struct {
	Algorithm string  `json:"algorithm"`
	Bits      uint    `json:"bits"`
	N         float64 `json:"n"`
	R         uint    `json:"r"`
	P         uint    `json:"p"`
	Salt      string  `json:"salt,omitempty"`
}

func NewKDF(r data.KDF) KDF {
	return KDF{
		Type: "kdf",
		ID:   r.Version,
		Attributes: KDFAttributes{
			Algorithm: r.Algorithm,
			Bits:      r.Bits,
			N:         r.N,
			R:         r.R,
			P:         r.P,
			Salt:      r.Salt,
		},
	}
}
