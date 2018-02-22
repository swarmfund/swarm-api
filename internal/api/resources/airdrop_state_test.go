package resources

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsEmailAirdropEligible(t *testing.T) {
	cases := []struct {
		expected bool
		email    string
	}{
		{true, "foo@gmail.com"},
		{true, "so@ba.ka"},
		{false, "very@qq.com"},
		{false, "sorry@QQ.com"},
	}

	for i, tc := range cases {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			got := isEmailAirdropEligible(tc.email)
			assert.Equal(t, tc.expected, got)
		})
	}
}
