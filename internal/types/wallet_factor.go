package types

import "errors"

type WalletFactorType string

const (
	WalletFactorTOTP     = "totp"
	WalletFactorPassword = "password"
)

var (
	ErrInvalidWalletFactorType = errors.New("backend type is invalid")
	tfaBackendTypes            = map[WalletFactorType]bool{
		WalletFactorTOTP:     true,
		WalletFactorPassword: true,
	}
)

func (t WalletFactorType) Validate() error {
	if _, ok := tfaBackendTypes[t]; !ok {
		return ErrInvalidWalletFactorType
	}
	return nil
}
