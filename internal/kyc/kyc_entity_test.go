package kyc

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/swarmfund/api/internal/types"
)

func TestKYCEntity_UnmarshalJSON(t *testing.T) {
	data := []byte(`{
		"type": "individual",
		"attributes": {
			"first_name": "yo",
			"last_name": "ba"
		}
	}`)
	expected := KYCEntity{
		Type: types.KYCEntityTypeIndividual,
		Individual: &Individual{
			FirstName: "yo",
			LastName:  "ba",
		},
	}
	var entity KYCEntity
	if err := json.Unmarshal(data, &entity); err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, expected, entity)
}
