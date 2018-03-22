package handlers_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/swarmfund/api/internal/api/handlers"
	"gitlab.com/swarmfund/api/internal/types"
)

func TestNewBlobIndexFilter(t *testing.T) {
	params := map[string]string{"address": "GAHOOHZTJDHYMLV5HP3GSUPWUOERGEOAWB52NBSV2IKR2225SB3SW2QK"}
	request := handlers.RequestWithURLParams([]byte(``), params)

	query := request.URL.Query()

	query.Set("type", "2048")
	query.Set("TBE", "5000")

	request.URL.RawQuery = query.Encode()

	got, err := handlers.NewBlobIndexFilter(request)
	if err != nil {
		t.Fatal(err)
	}

	blobType := types.BlobTypeKYCForm

	expected := handlers.BlobIndexFilter{
		Page:          1,
		Address:       "GAHOOHZTJDHYMLV5HP3GSUPWUOERGEOAWB52NBSV2IKR2225SB3SW2QK",
		Type:          &blobType,
		Relationships: map[string]string{"TBE": "5000"},
	}

	assert.EqualValues(t, expected, got)
}
