package middlewares

import (
	"bytes"
	"testing"

	"net/http"
	"net/http/httptest"

	"strings"

	"github.com/sirupsen/logrus"
	"gitlab.com/swarmfund/api/log"
)

func TestLogger(t *testing.T) {
	var buf bytes.Buffer
	logger := logrus.New()
	logger.Out = &buf
	entry := log.Entry{}
	entry.Logger = logger

	middleware := Logger(&entry)
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(345)
	})

	ts := httptest.NewServer(middleware(handler))
	defer ts.Close()

	response, err := http.Get(ts.URL)
	if err != nil {
		t.Fatal(err)
	}
	defer response.Body.Close()

	out := buf.String()
	lines := strings.Count(out, "\n")
	if lines != 2 {
		t.Errorf("expected 2 lines of output got %d", lines)
	}

	if strings.Count(out, "345") != 1 {
		t.Error("expected to contain status code")
	}
}
