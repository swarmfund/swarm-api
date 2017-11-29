package middlewares

import (
	"net/http"
	"testing"

	"fmt"
	"io/ioutil"
	"net/http/httptest"
)

func TestRecover(t *testing.T) {
	middleware := Recover(func(w http.ResponseWriter, r *http.Request, rvr interface{}) {
		fmt.Println(rvr)
	})
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		panic("test")
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
	// Output: test
}
