package handlers

import (
	"context"
	"net/http"
	"testing"

	"bytes"

	"reflect"

	"github.com/go-chi/chi"
	"gitlab.com/swarmfund/api/internal/types"
)

// TODO move to ape
func RequestWithURLParams(body []byte, params map[string]string) *http.Request {
	rctx := chi.NewRouteContext()
	for key, value := range params {
		rctx.URLParams.Add(key, value)
	}
	r, _ := http.NewRequest("GET", "/", bytes.NewReader(body))
	r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))
	return r
}

func TestNewCreateUserRequest(t *testing.T) {
	cases := []struct {
		name     string
		address  string
		body     string
		err      bool
		expected CreateUserRequest
	}{
		{
			"valid",
			"GCSWI5EKDRNXBRQUY2M27CSTYQHST2S6ONLC5W5V2O4E6OTABR4CRORF",
			`{
				"data": {
					"attributes": {
						"type": 1
					}
				}
			}`,
			false,
			CreateUserRequest{
				Address: "GCSWI5EKDRNXBRQUY2M27CSTYQHST2S6ONLC5W5V2O4E6OTABR4CRORF",
				// FIXME I love jsonapi, see implementation for details
				Type:     1,
				UserType: types.UserType(1),
			},
		},
		{
			"string type",
			"GCSWI5EKDRNXBRQUY2M27CSTYQHST2S6ONLC5W5V2O4E6OTABR4CRORF",
			`{
				"data": {
					"attributes": {
						"type": "1"
					}
				}
			}`,
			true,
			CreateUserRequest{},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			r := RequestWithURLParams([]byte(tc.body), map[string]string{
				"address": tc.address,
			})
			got, err := NewCreateUserRequest(r)
			if err != nil && !tc.err {
				t.Fatalf("expected nil error got %s", err)
			}
			if err == nil && tc.err {
				t.Fatalf("expected error got nil")
			}
			if err == nil && !reflect.DeepEqual(got, tc.expected) {
				t.Fatal("expected %#v got #%v", tc.expected, got)
			}
		})
	}

}
