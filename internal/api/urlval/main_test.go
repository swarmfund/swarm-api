package urlval

import (
	"net/http"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDecode(t *testing.T) {
	type Filters struct {
		Page  uint64  `url:"page"`
		State *uint64 `url:"state"`
	}
	uint := uint64(42)
	cases := []struct {
		name     string
		values   url.Values
		expected Filters
	}{
		{"empty", map[string][]string{}, Filters{}},
		{"*uint64", map[string][]string{"state": {"42"}}, Filters{0, &uint}},
		{"uint64", map[string][]string{"page": {"42"}}, Filters{uint, nil}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var got Filters
			Decode(tc.values, &got)
			assert.Equal(t, tc.expected, got)
		})
	}
}

func TestEncode(t *testing.T) {
	type Fields struct {
		Page  uint64  `url:"page"`
		State *uint64 `url:"state"`
	}
	cases := []struct {
		name     string
		fields   Fields
		expected FilterLinks
	}{
		{"page", Fields{10, nil}, FilterLinks{"/users?page=10", "/users?page=11", "/users?page=9"}},
		{"first page", Fields{1, nil}, FilterLinks{"/users?page=1", "/users?page=2", "/users?"}},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := Encode(&http.Request{URL: &url.URL{Path: "/users"}}, tc.fields)
			assert.Equal(t, got, tc.expected)
		})
	}
}

func TestDecodeWithValue(t *testing.T) {
	type Filter struct {
		Page          uint64            `url:"page"`
		Relationships map[string]string `url:"relationships"`
	}

	var urlValues url.Values = map[string][]string{"page": {"23"}, "relationships": {"U CAN'T SET ME!!!"}}
	expected := Filter{Page: 23}

	var got Filter
	if err := DecodeWithValues(&urlValues, &got); err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, expected.Page, got.Page)
	//value have been deleted from map
	assert.Equal(t, "", urlValues.Get("unvalid"))
	//check unsupported value
	assert.Equal(t, "U CAN'T SET ME!!!", urlValues.Get("relationships"))
}
