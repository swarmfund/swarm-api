package lorem

import "math/rand"

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

func Token() string {
	return randomString(20, letters)
}
