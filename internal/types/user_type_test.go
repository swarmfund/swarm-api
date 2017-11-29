package types

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestUserType_Validate(t *testing.T) {
	cases := []struct {
		in       int
		expected error
	}{
		{0, ErrUserTypeInvalid},
		{1, nil},
		{2, ErrUserTypeInvalid},
	}

	for i, tc := range cases {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			tpe := UserType(tc.in)
			if err := tpe.Validate(); err != tc.expected {
				t.Fatalf("got %s expected %s", err, tc.expected)
			}
		})
	}
}

func TestUserType_UnmarshalJSON(t *testing.T) {
	cases := []struct {
		in       string
		expected UserType
		err      bool
	}{
		{`{"t":1}`, UserType(1), false},
		{`{"t":"2"}`, UserType(2), false},
		{`{"t":"individual"}`, UserType(0), true},
		{`{"t":{}}"`, UserType(0), true},
		{`{"t":[]}`, UserType(0), true},
		{`{"t":{yolo}}`, UserType(0), true},
	}

	for i, tc := range cases {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			var got struct {
				T UserType `json:"t"`
			}
			err := json.Unmarshal([]byte(tc.in), &got)
			if err != nil && !tc.err {
				t.Fatalf("expected nil error got %s", err)
			}
			if err == nil && tc.err {
				t.Fatal("expected error got nil")
			}
			if err == nil && tc.expected != got.T {
				t.Fatal("expected %d got %d", tc.expected, got)
			}
		})
	}
}
