package storage

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/swarmfund/api/internal/types"
)

func TestKeyMarshal(t *testing.T) {
	key := NewKey(2<<62-1, types.DocumentTypeAssetLogo)

	encoded, err := key.MarshalText()
	assert.NoError(t, err)
	assert.Len(t, encoded, 56)

	decoded := Key{}
	err = decoded.UnmarshalText(encoded)
	assert.NoError(t, err)
	assert.Equal(t, *key, decoded)
}
