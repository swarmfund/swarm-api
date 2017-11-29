package secondfactor

import (
	"context"
	"crypto/sha1"
	"hash"
	"io"
	"io/ioutil"
	"net/http"

	"bytes"
	"encoding/hex"
	"fmt"
)

type ctxKey int

const (
	_ ctxKey = iota
	hasherCtxKey
)

var (
	// DefaultHash hash.Hash implementation used for HashReader
	DefaultHash = sha1.New
	// DefaultNilHash is a sum of nil input, used to check if body were read or not
	DefaultNilHash = DefaultHash().Sum(nil)
)

type HashReader struct {
	io.Reader
	hash.Hash
}

func NewHashReader(r io.Reader, hash hash.Hash) *HashReader {
	return &HashReader{io.TeeReader(r, hash), hash}
}

// HashMiddleware set request body hash to context
func HashMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			hash := DefaultHash()
			hasher := NewHashReader(r.Body, hash)
			r.Body = ioutil.NopCloser(hasher)
			ctx := context.WithValue(r.Context(), hasherCtxKey, hasher)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// RequestHash returns request hash generated from url and body
// request body should be read before getting hash
// HashReader is expected in context
func RequestHash(r *http.Request) string {
	// get hasher from request context
	ctxvalue := r.Context().Value(hasherCtxKey)
	if ctxvalue == nil {
		panic("hasher expected in ctx")
	}
	hasher := ctxvalue.(*HashReader)

	hash := hasher.Sum(nil)

	// naive check of was body read (at least partially) or not
	if r.ContentLength > 0 && bytes.Equal(DefaultNilHash, hash) {
		panic("body should be read before hashing")
	}

	target := fmt.Sprintf("%s:%s:%x", r.Method, r.URL.Path, hash)
	raw := sha1.Sum([]byte(target))
	return hex.EncodeToString(raw[:])
}
