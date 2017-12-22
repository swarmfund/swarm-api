package kyc

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/swarmfund/api/internal/types"
)

func TestEntity_Validate(t *testing.T) {
	cases := []struct {
		name   string
		entity Entity
		valid  bool
	}{
		{
			"valid",
			Entity{
				Type: types.KYCEntityTypeIndividual,
			},
			true,
		},
		{
			"invalid type",
			Entity{
				Type: 42,
			},
			false,
		},
		{
			"missing type",
			Entity{},
			false,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.entity.Validate()
			if tc.valid {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
			}
		})
	}
}
