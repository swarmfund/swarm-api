package types

import (
	"fmt"
	"testing"

	"encoding/json"

	"github.com/stretchr/testify/assert"
)

func TestUserType_Validate(t *testing.T) {
	cases := []struct {
		in       int
		expected error
	}{
		{0, ErrUserTypeInvalid},
		{1, nil},
		{2, nil},
		{3, ErrUserTypeInvalid},
		{4, ErrUserTypeInvalid},
	}

	for i, tc := range cases {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			//tpe := UserType(tc.in)
			assert.EqualValues(t, UserType(tc.in).Validate(), tc.expected)
		})
	}
}

func TestUserType_UnmarshalJSON(t *testing.T) {
	cases := []struct {
		name     string
		in       string
		expected UserType
		err      bool
	}{
		{"int", `1`, 1, true},
		{"int as string", `"2"`, 0, true},
		{"string", `"syndicate"`, 2, false},
		{"object", `{}`, 0, true},
		{"array", `[]`, 0, true},
		{"invalid json", `}{`, 0, true},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var got UserType
			err := json.Unmarshal([]byte(tc.in), &got)
			if err != nil && !tc.err {
				t.Fatalf("expected nil error got %s", err)
			}
			if err == nil && tc.err {
				t.Fatal("expected error got nil")
			}
			if err == nil && tc.expected != got {
				t.Fatalf("expected %d got %s", tc.expected, got)
			}
		})
	}
}
