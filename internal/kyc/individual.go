package kyc

import "github.com/go-ozzo/ozzo-validation"

type Individual struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

func (e Individual) Validate() error {
	return validation.ValidateStruct(&e,
		validation.Field(&e.FirstName, validation.Required),
		validation.Field(&e.LastName, validation.Required),
	)
}
