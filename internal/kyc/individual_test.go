package kyc

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIndividual_Validate(t *testing.T) {
	cases := []struct {
		name   string
		entity Individual
		valid  bool
	}{
		{
			"valid",
			Individual{
				FirstName: "John",
				LastName:  "Doe",
			},
			true,
		},
		{
			"missing first name",
			Individual{
				LastName: "Doe",
			},
			false,
		},
		{
			"missing last name",
			Individual{
				LastName: "Doe",
			},
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
