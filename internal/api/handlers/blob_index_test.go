package handlers_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/distributed_lab/ape/apeutil"
	"gitlab.com/swarmfund/api/internal/api/handlers"
)

func TestNewBlobIndexRequest(t *testing.T) {
	params := map[string]string{"address": "GAHOOHZTJDHYMLV5HP3GSUPWUOERGEOAWB52NBSV2IKR2225SB3SW2QK"}
	request := apeutil.RequestWithURLParams([]byte(``), params)

	query := request.URL.Query()

	//this values should set filter and then delete himself from query
	query.Set("page", "42")
	query.Set("type", "2048")

	//this should set in relationships
	query.Set("TBE", "5000")
	query.Set("NUMBER", "42")
	request.URL.RawQuery = query.Encode()

	got, err := handlers.NewBlobIndexRequest(request)
	if err != nil {
		t.Fatal(err)
	}

	blobType := uint64(2048)

	expected := handlers.BlobIndexRequest{
		Filter: handlers.BlobIndexFilter{
			Page: 42,
			Type: &blobType,
		},
		Address:       "GAHOOHZTJDHYMLV5HP3GSUPWUOERGEOAWB52NBSV2IKR2225SB3SW2QK",
		Relationships: map[string]string{"TBE": "5000", "NUMBER": "42"},
	}

	assert.EqualValues(t, expected, got)
}
