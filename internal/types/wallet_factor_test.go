package types

import (
	"fmt"
	"testing"
)

func TestWalletFactorType_Validate(t *testing.T) {
	cases := []struct {
		in  string
		err error
	}{
		{"totp", nil},
		{"password", nil},
		{"not-a-factor", ErrInvalidWalletFactorType},
	}

	for i, tc := range cases {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			f := WalletFactorType(tc.in)
			err := f.Validate()
			if err != tc.err {
				t.Fatalf("expected %s got %s", tc.err, err)
			}
		})
	}
}
