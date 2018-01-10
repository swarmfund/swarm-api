package kyc

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSyndicate_Validate(t *testing.T) {
	cases := []struct {
		name   string
		entity Syndicate
		valid  bool
	}{
		{
			"valid",
			Syndicate{
				Name: "Yoba Inc.",
			},
			true,
		},
		{
			"missing name",
			Syndicate{},
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
