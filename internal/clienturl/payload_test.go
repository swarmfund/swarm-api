package clienturl

import (
	"encoding/base64"
	"testing"
)

func TestPayload_Encode(t *testing.T) {
	cases := []struct {
		name     string
		payload  Payload
		expected string
	}{
		{
			"simple",
			Payload{
				Status: 500,
				Type:   1,
				Meta: map[string]interface{}{
					"foo": "bar",
				},
			},
			`{"status":500,"type":1,"meta":{"foo":"bar"}}`,
		}, {
			"omitted status",
			Payload{
				Status: 200,
				Type:   1,
				Meta: map[string]interface{}{
					"foo": "bar",
				},
			},
			`{"type":1,"meta":{"foo":"bar"}}`,
		},
		{
			"empty meta",
			Payload{
				Status: 500,
				Type:   1,
				Meta:   map[string]interface{}{},
			},
			`{"status":500,"type":1}`,
		},
		{
			"nil meta",
			Payload{
				Status: 500,
				Type:   1,
			},
			`{"status":500,"type":1}`,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			encoded, err := tc.payload.Encode()
			if err != nil {
				t.Fatal(err)
			}

			json, err := base64.URLEncoding.WithPadding(base64.NoPadding).DecodeString(encoded)
			if err != nil {
				t.Fatal(err)
			}

			if string(json) != tc.expected {
				t.Fatalf("expected %s got %s", json, tc.expected)
			}

		})
	}
}
