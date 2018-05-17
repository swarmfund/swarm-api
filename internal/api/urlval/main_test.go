package urlval

import (
	"net/http"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDecode(t *testing.T) {
	type Filters struct {
		Page    uint64  `url:"page"`
		State   *uint64 `url:"state"`
		Boolean *bool   `url:"boolean"`
		PString *string `url:"pstring"`
	}
	uint := uint64(42)
	pstring := "Vasyl Lomachenko"
	var boolean = true
	cases := []struct {
		name     string
		values   url.Values
		expected Filters
	}{
		{"empty", url.Values{}, Filters{}},
		{"*uint64", url.Values{"state": {"42"}}, Filters{0, &uint, nil, nil}},
		{"uint64", url.Values{"page": {"42"}}, Filters{uint, nil, nil, nil}},
		{"bool", url.Values{"page": {"42"}, "boolean": {"true"}}, Filters{uint, nil, &boolean, nil}},
		{"*string", url.Values{"page": {"42"}, "pstring": {"Vasyl Lomachenko"}}, Filters{uint, nil, nil, &pstring}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var got Filters
			assert.NoError(t, Decode(tc.values, &got))
			assert.Equal(t, tc.expected, got)
		})
	}
}

func TestEncode(t *testing.T) {
	type Fields struct {
		Page    uint64  `url:"page"`
		State   *uint64 `url:"state"`
		Name    *string `url:"name"`
		Allowed *bool   `url:"allowed"`
	}

	state := uint64(42)
	name := "yoba"
	allowed := true
	cases := []struct {
		name     string
		fields   Fields
		expected FilterLinks
	}{
		{"page", Fields{Page: 10}, FilterLinks{"/users?page=10", "/users?page=11", "/users?page=9"}},
		{"first page", Fields{Page: 1}, FilterLinks{"/users?page=1", "/users?page=2", ""}},
		{"state", Fields{Page: 5, State: &state}, FilterLinks{"/users?page=5&state=42", "/users?page=6&state=42", "/users?page=4&state=42"}},
		{"state first page", Fields{Page: 1, State: &state}, FilterLinks{"/users?page=1&state=42", "/users?page=2&state=42", ""}},
		{"string encode", Fields{Page: 5, Name: &name}, FilterLinks{"/users?name=yoba&page=5", "/users?name=yoba&page=6", "/users?name=yoba&page=4"}},
		{"bool encode", Fields{Page: 5, Allowed: &allowed}, FilterLinks{"/users?allowed=true&page=5", "/users?allowed=true&page=6", "/users?allowed=true&page=4"}},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := Encode(&http.Request{URL: &url.URL{Path: "/users"}}, tc.fields)
			assert.Equal(t, tc.expected, got)
		})
	}

}
