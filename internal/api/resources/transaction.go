package resources

import (
	. "github.com/go-ozzo/ozzo-validation"
)

type Transaction struct {
	Data TransactionData `json:"data"`
}

func (r Transaction) Validate() error {
	return ValidateStruct(&r,
		Field(&r.Data, Required),
	)
}

type TransactionData struct {
	Attributes TransactionAttributes `json:"attributes"`
}

func (r TransactionData) Validate() error {
	return ValidateStruct(&r,
		Field(&r.Attributes, Required),
	)
}

type TransactionAttributes struct {
	Envelope string `json:"envelope"`
}

func (r TransactionAttributes) Validate() error {
	return ValidateStruct(&r,
		Field(&r.Envelope, Required),
	)
}
