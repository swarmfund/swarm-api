package resources

import (
	. "github.com/go-ozzo/ozzo-validation"
)

type PasswordFactor struct {
	ID           int64  `json:"id" jsonapi:"primary,password"`
	Salt         string `json:"salt" jsonapi:"attr,salt"`
	AccountID    string `json:"account_id" jsonapi:"attr,account_id"`
	KeychainData string `json:"keychain_data" jsonapi:"attr,keychain_data"`
}

func (r PasswordFactor) Validate() error {
	return ValidateStruct(&r,
		Field(&r.Salt, Required),
		// TODO address validation
		Field(&r.AccountID, Required),
		Field(&r.KeychainData, Required),
	)
}
