package lorem

import (
	"math/rand"

	"encoding/base64"

	"github.com/pkg/errors"
)

var (
	letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	digits  = "0123456789"
)

func randomString(n int, source string) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = source[rand.Intn(len(source))]
	}
	return string(b)
}

func Salt() string {
	r := make([]byte, 30)
	_, err := rand.Read(r)
	if err != nil {
		panic(errors.Wrap(err, "failed to read random"))
	}
	return base64.StdEncoding.EncodeToString(r)
}

func Token() string {
	return randomString(20, letters)
}
