package middlewares

import (
	"bytes"
	"testing"

	"net/http"
	"net/http/httptest"

	"strings"

	"time"

	"gitlab.com/distributed_lab/logan/v3"
)

func TestLogger(t *testing.T) {
	var buf bytes.Buffer
	logger := logan.New().Out(&buf)

	middleware := Logger(logger, 1*time.Hour)
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
