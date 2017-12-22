package lorem

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestULID(t *testing.T) {
	first := ULID()
	second := ULID()
	assert.NotEqual(t, first, second)
}
