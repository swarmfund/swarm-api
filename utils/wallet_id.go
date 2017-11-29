package utils

import (
	"encoding/base64"

	"encoding/hex"

	"github.com/pkg/errors"
)

func Base64ToHex(b64 string) (string, error) {
	bytes, err := base64.StdEncoding.DecodeString(b64)
	if err != nil {
		return "", errors.Wrap(err, "failed to decode b64")
	}
	return hex.EncodeToString(bytes), nil
}

func HexToBase64(h string) (string, error) {
	bytes, err := hex.DecodeString(h)
	if err != nil {
		return "", errors.Wrap(err, "failed to decode hex")
	}
	return base64.StdEncoding.EncodeToString(bytes), nil
}
