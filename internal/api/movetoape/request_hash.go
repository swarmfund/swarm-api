package movetoape

import (
	"fmt"
	"net/http"

	"encoding/hex"

	"gitlab.com/tokend/go/hash"
)

// TODO move to ape

func RequestHash(r *http.Request) string {
	raw := hash.Hash([]byte(
		fmt.Sprintf("%s:%s:%d", r.Method, r.URL.Path, r.ContentLength)))
	return hex.EncodeToString(raw[:])
}
