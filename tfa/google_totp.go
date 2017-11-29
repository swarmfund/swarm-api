package tfa

import (
	"encoding/json"

	"github.com/pkg/errors"
	"github.com/pquerna/otp"
	"github.com/pquerna/otp/totp"
)

type GoogleTOPT struct {
	Secret string `json:"secret"`
	key    *otp.Key
	id     int64
}

func NewTOTPBackend(issuer, account string) (Backend, error) {
	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      issuer,
		AccountName: account,
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed to generate key")
	}
	backend := GoogleTOPT{
		Secret: key.Secret(),
		key:    key,
	}
	return backend, nil
}

func NewTOTPFromDB(id int64, details []byte) (*GoogleTOPT, error) {
	var result GoogleTOPT
	if err := json.Unmarshal(details, &result); err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal details")
	}
	result.id = id
	return &result, nil
}

func (b GoogleTOPT) ID() int64 {
	return b.id
}

func (b GoogleTOPT) Attributes() map[string]interface{} {
	return map[string]interface{}{
		"secret": b.Secret,
		"seed":   b.key.String(),
	}
}

func (b GoogleTOPT) Deliver(bytes []byte) (map[string]interface{}, error) {
	return map[string]interface{}{}, nil
}

func (b GoogleTOPT) Verify(code string, _ string) (bool, error) {
	ok := totp.Validate(code, b.Secret)
	return ok, nil
	return false, nil
}
