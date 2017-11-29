package middlewares

import (
	"context"
	"net/http"
	"testing"

	"fmt"
	"io/ioutil"
	"net/http/httptest"
)

func TestCtx(t *testing.T) {
	key := "a"
	value := "value"
	middleware := Ctx(func(ctx context.Context) context.Context {
		return context.WithValue(ctx, key, value)
	})
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(fmt.Sprintf("%#v", r.Context().Value(key))))
	})

	ts := httptest.NewServer(middleware(handler))
	defer ts.Close()

	response, err := http.Get(ts.URL)
	if err != nil {
		t.Fatal(err)
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Printf("%s", body)
	// Output: value
}
