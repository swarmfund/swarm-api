package handlers_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/swarmfund/api/internal/api/handlers"
	"gitlab.com/swarmfund/api/internal/types"
)

func TestNewBlobIndexRequest(t *testing.T) {
	params := map[string]string{"address": "GAHOOHZTJDHYMLV5HP3GSUPWUOERGEOAWB52NBSV2IKR2225SB3SW2QK"}
	request := handlers.RequestWithURLParams([]byte(``), params)

	query := request.URL.Query()
	query.Set("type", "2048")
	query.Set("TBE", "5000")

	request.URL.RawQuery = query.Encode()

	blobRequest, err := handlers.NewBlobIndexRequest(request)
	if err != nil {
		t.Fatal(err)
	}

	blobType := types.BlobTypeKYCForm

	expectedRequest := handlers.BlobIndexRequest{
		Address:       "GAHOOHZTJDHYMLV5HP3GSUPWUOERGEOAWB52NBSV2IKR2225SB3SW2QK",
		Relationships: map[string]string{"TBE": "5000"},
		Type:          &blobType,
	}

	assert.Equal(t, expectedRequest.Address, blobRequest.Address)
	assert.Equal(t, *expectedRequest.Type, *blobRequest.Type)
	assert.EqualValues(t, expectedRequest.Relationships, blobRequest.Relationships)
}

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

	blobType := uint64(types.BlobTypeKYCForm)

	expected := handlers.BlobIndexFilter{
		Page:          1,
		Address:       "GAHOOHZTJDHYMLV5HP3GSUPWUOERGEOAWB52NBSV2IKR2225SB3SW2QK",
		Type:          &blobType,
		Relationships: map[string]string{"TBE": "5000"},
	}

	assert.Equal(t, expected.Page, got.Page)
	assert.Equal(t, expected.Address, got.Address)
	assert.Equal(t, expected.Type, got.Type)
	//TODO add new feature to urlvalue
	//assert.EqualValues(t, expected.Relationships, got.Relationships)
}
