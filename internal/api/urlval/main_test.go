package urlval

import (
	"fmt"
	"net/http"
	"net/url"
	"testing"
)

func TestPopulate(t *testing.T) {
	values := url.Values{}
	values.Add("page", "1")
	values.Add("state", "2")
	s := struct {
		Page  uint64  `url:"page"`
		State *uint64 `url:"state"`
	}{}
	Decode(values, &s)
}

func TestEncode(t *testing.T) {
	f := struct {
		Page uint64 `url:"page"`
	}{
		Page: 3,
	}
	links := Encode(&http.Request{URL: &url.URL{Path: "/users/"}}, f)
	fmt.Println(links)
}
