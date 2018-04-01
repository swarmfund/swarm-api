package handlers

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/swarmfund/api/internal/types"
)

func TestNewCreateBlobRequest(t *testing.T) {
	data := `{
		"data": {
			"type": "kyc_form",
			"attributes": {
				"value": "foobar"
			}
		}
	}`

	t.Run("nil address is valid", func(t *testing.T) {
		r := RequestWithURLParams([]byte(data), map[string]string{})
		_, err := NewCreateBlobRequest(r)
		assert.NoError(t, err)
	})

	t.Run("invalid address", func(t *testing.T) {
		r := RequestWithURLParams([]byte(data), map[string]string{
			"address": "GINVALIDADDRESS",
		})
		_, err := NewCreateBlobRequest(r)
		assert.Error(t, err)
	})

	t.Run("valid address", func(t *testing.T) {
		address := "GB22CHWJQSDQ4VIP7A6QHDIFMTT2FCZKTOJGHCK56UKWR6VFVUDJ5RCR"
		r := RequestWithURLParams([]byte(data), map[string]string{
			"address": address,
		})
		got, err := NewCreateBlobRequest(r)
		assert.NoError(t, err)
		assert.NotNil(t, got.Address)
		assert.EqualValues(t, address, *got.Address)
	})

	t.Run("blob", func(t *testing.T) {
		r := RequestWithURLParams([]byte(data), map[string]string{})
		got, err := NewCreateBlobRequest(r)
		assert.NoError(t, err)
		assert.EqualValues(t, types.BlobTypeKYCForm, got.Data.Type)
		assert.EqualValues(t, "foobar", got.Data.Attributes.Value)
	})
}
