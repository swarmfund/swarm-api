package tfa

import (
	"encoding/base64"
	"fmt"

	"gitlab.com/tokend/go/hash"
)

type Backend interface {
	ID() int64
	Attributes() map[string]interface{}
	Verify(code string, token string) (bool, error)
}

type HasMeta interface {
	Meta() map[string]interface{}
}

type TFA struct {
	OwnerID   int64
	OTPData   interface{}
	Token     string
	BackendID int64
}

func Token(owner int64, action string) string {
	hashSource := []byte(fmt.Sprintf("%d:%s", owner, action))
	tokenBytes := hash.Hash(hashSource)
	t := base64.StdEncoding.EncodeToString(tokenBytes[:])
	return t
}
