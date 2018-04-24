package tfa

import (
	"encoding/json"

	"encoding/base64"

	"github.com/pkg/errors"
	"gitlab.com/swarmfund/api/internal/types"
	"gitlab.com/tokend/go/keypair"
)

type PasswordDetails struct {
	WalletID     string        `json:"wallet_id"`
	Salt         string        `json:"salt"`
	AccountID    types.Address `json:"account_id"`
	KeychainData string        `json:"keychain_data"`
}

type Password struct {
	Details PasswordDetails `json:"details"`
	id      int64
}

func NewPasswordBackend(details PasswordDetails) *Password {
	return &Password{
		Details: details,
	}
}

func NewPasswordFromDB(id int64, details []byte) (*Password, error) {
	backend := Password{}
	backend.id = id
	if err := json.Unmarshal(details, &backend); err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal details")
	}
	return &backend, nil
}

func (p Password) ID() int64 {
	return p.id
}

func (p Password) Attributes() map[string]interface{} {
	return map[string]interface{}{}
}

func (p Password) Meta() map[string]interface{} {
	return map[string]interface{}{
		"salt":          p.Details.Salt,
		"keychain_data": p.Details.KeychainData,
		"factor_type":   types.WalletFactorPassword,
	}
}

func (p Password) Verify(code string, token string) (bool, error) {
	kp, err := keypair.Parse(string(p.Details.AccountID))
	if err != nil {
		return false, errors.Wrap(err, "failed to parse keypair")
	}
	sign, err := base64.StdEncoding.DecodeString(code)
	if err != nil {
		return false, errors.Wrap(err, "failed to decode otp")
	}
	if err := kp.Verify([]byte(token), sign); err != nil {
		return false, nil
	}
	return true, nil
}
