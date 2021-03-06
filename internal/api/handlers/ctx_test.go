package handlers

import (
	"context"
	"net/http"
	"testing"

	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/swarmfund/api/db2/api"
	"gitlab.com/swarmfund/api/internal/data"
	"gitlab.com/tokend/horizon-connector"
)

func TestCtxLog(t *testing.T) {
	entry := logan.New()
	ctx := context.Background()
	request := &http.Request{}
	request = request.WithContext(CtxLog(entry)(ctx))
	got := Log(request)
	if got != entry {
		t.Fatalf("expected %#v got %#v", entry, got)
	}
}

func TestCtxEmailTokensQ(t *testing.T) {
	expected := &api.EmailTokensQ{}
	ctx := context.Background()
	request := &http.Request{}
	request = request.WithContext(CtxEmailTokensQ(expected)(ctx))
	got := EmailTokensQ(request)
	if _, ok := got.(data.EmailTokensQ); !ok {
		t.Fatalf("expected %T got %T", expected, got)
	}
}

func TestCtxHorizon(t *testing.T) {
	expected := &horizon.Connector{}
	ctx := context.Background()
	request := &http.Request{}
	request = request.WithContext(CtxHorizon(expected)(ctx))
	got := Horizon(request)
	if got != expected {
		t.Fatalf("expected %#v got %#v", expected, got)
	}
}
