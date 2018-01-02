package kyc

import "github.com/go-ozzo/ozzo-validation"

type Syndicate struct {
	Name string `json:"name"`
}

func (e Syndicate) Validate() error {
	return validation.ValidateStruct(&e,
		validation.Field(&e.Name, validation.Required),
	)
}
