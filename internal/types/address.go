package types

import (
	"errors"

	"gitlab.com/swarmfund/go/strkey"
)

var (
	ErrAddressInvalid = errors.New("address is invalid")
)

type Address string

func (a Address) Validate() error {
	_, err := strkey.Decode(strkey.VersionByteAccountID, string(a))
	if err != nil {
		return ErrAddressInvalid
	}
	return nil
}
