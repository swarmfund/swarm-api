package types

import (
	"fmt"
	"testing"
)

func TestAddress_Validate(t *testing.T) {
	cases := []struct {
		in       string
		expected error
	}{
		{"FFF", ErrAddressInvalid},
		// seed is not address
		{"SAUHKDQ4RFHH2MA4QGDWEYFXU6NZKWNFABONV3UW54X564XJDW7QJ5RX", ErrAddressInvalid},
		{"GD5UKTLHHNEFDSXL6EGVJ3O6QEMO5MXPVWF7E7BMAEO5WN3S3QBF6JZJ", nil},
	}

	for i, tc := range cases {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			addr := Address(tc.in)
			if err := addr.Validate(); err != tc.expected {
				t.Fatalf("expected %s got %s", tc.expected, err)
			}
		})
	}
}
