package resources

import (
	. "github.com/go-ozzo/ozzo-validation"
)

type Transaction struct {
	Envelope string `json:"envelope" jsonapi:"attr,envelope"`
}

func (r Transaction) Validate() error {
	return ValidateStruct(&r,
		Field(&r.Envelope, Required),
	)
}
