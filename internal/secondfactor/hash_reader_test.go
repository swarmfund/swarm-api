package secondfactor

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRequestHash(t *testing.T) {
	body := []byte(`spamegg`)

	server := func(handler http.HandlerFunc) *httptest.Server {
		middleware := HashMiddleware()
		return httptest.NewServer(middleware(handler))
	}

	t.Run("without hasher", func(t *testing.T) {
		r, err := http.NewRequest("PATCH", "/foo/bar", bytes.NewReader(body))
		defer func() {
			if rvr := recover(); rvr == nil {
				t.Fatal("panic expected")
			}
		}()

		if err != nil {
			t.Fatal(err)
		}

		RequestHash(r)
	})

	t.Run("not read body", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if rvr := recover(); rvr == nil {
					t.Fatal("panic expected")
				}
			}()
			RequestHash(r)
		})

		ts := server(handler)
		defer ts.Close()

		response, err := http.Post(ts.URL, "application/json", bytes.NewReader(body))
		if err != nil {
			t.Fatal(err)
		}
		response.Body.Close()

	})

	t.Run("multiple calls", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			_, err := ioutil.ReadAll(r.Body)
			if err != nil {
				t.Fatal(err)
			}
			a := RequestHash(r)
			b := RequestHash(r)
			if a != b {
				t.Fatal("hashes should match")
			}
			if a == "" {
				t.Fatal("hash should not be empty")
			}
		})

		ts := server(handler)
		defer ts.Close()

		response, err := http.Post(ts.URL, "application/json", bytes.NewReader(body))
		if err != nil {
			t.Fatal(err)
		}
		response.Body.Close()
	})

	t.Run("multiple requests", func(t *testing.T) {
		hashes := []string{}
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			_, err := ioutil.ReadAll(r.Body)
			if err != nil {
				t.Fatal(err)
			}
			hashes = append(hashes, RequestHash(r))
		})

		ts := server(handler)
		defer ts.Close()

		response, err := http.Post(ts.URL, "application/json", bytes.NewReader(body))
		if err != nil {
			t.Fatal(err)
		}
		response.Body.Close()

		response, err = http.Post(ts.URL, "application/json", bytes.NewReader(body))
		if err != nil {
			t.Fatal(err)
		}
		response.Body.Close()

		response, err = http.Post(ts.URL, "application/json", bytes.NewReader([]byte(`johndoe`)))
		if err != nil {
			t.Fatal(err)
		}
		response.Body.Close()

		if hashes[0] != hashes[1] {
			t.Fatal("hashes should match")
		}

		if hashes[1] == hashes[2] {
			t.Fatal("hashes should be different")
		}
	})
}
