package lorem

import (
	"time"

	"crypto/rand"

	"github.com/oklog/ulid"
)

func ULID() string {
	return ulid.MustNew(ulid.Timestamp(time.Now()), rand.Reader).String()
}
