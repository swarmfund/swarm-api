package discourse

import (
	"net/url"
	"testing"
)

func TestConnector_CreateUser(t *testing.T) {
	endpoint, err := url.Parse("http://localhost:3000")
	if err != nil {
		t.Fatal(err)
	}
	connector := NewConnector(endpoint, "fo1", "fc8f661f1a1ed3444be10861fb5d9b64e1fdbf7f04813712c5253071ad101c24")
	err = connector.CreateUser(CreateUser{
		Email: "yoba@sobaka.com",
	})
	if err != nil {
		t.Fatal(err)
	}
}
