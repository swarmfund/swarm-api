package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"fmt"

	"github.com/go-chi/chi"
	"github.com/google/jsonapi"
	"github.com/magiconair/properties/assert"
	"gitlab.com/swarmfund/api/db2/api/mocks"
	"gitlab.com/swarmfund/api/internal/api/middlewares"
	"gitlab.com/swarmfund/api/internal/api/resources"
)

func TestGetKDF(t *testing.T) {
	walletQ := mocks.WalletQI{}
	router := chi.NewRouter()
	router.Use(middlewares.Ctx(
		CtxWalletQ(&walletQ),
	))
	router.Get("/", GetKDF)

	ts := httptest.NewServer(router)
	defer ts.Close()

	get := func(email string) (resources.KDF, int) {
		url := ts.URL
		if email != "" {
			url += fmt.Sprintf("?email=%s", email)
		}
		response, err := http.Get(url)
		if err != nil {
			t.Fatal(err)
		}
		defer response.Body.Close()

		var got resources.KDF
		if response.StatusCode == 200 {
			if err := jsonapi.UnmarshalPayload(response.Body, &got); err != nil {
				t.Fatal(err)
			}
		}
		return got, response.StatusCode
	}

	t.Run("default", func(t *testing.T) {
		_, status := get("")
		if status != 200 {
			t.Fatalf("expected 200 got %d", status)
		}
	})

	t.Run("invalid email", func(t *testing.T) {
		walletQ.On("New").Return(&walletQ).Once()
		walletQ.On("ByEmail", "not@existing.com").Return(nil, nil).Once()
		defer walletQ.AssertExpectations(t)
		_, status := get("not@existing.com")
		assert.Equal(t, status, http.StatusNotFound)
	})

	//TODO FIX ME
	//t.Run("valid email", func(t *testing.T) {
	//	wallet := api2.Wallet{
	//		Salt: "salty",
	//		KDF:  1,
	//	}
	//	walletQ.On("New").Return(&walletQ).Once()
	//	walletQ.On("ByEmail", "do@exist.com").Return(&wallet, nil).Once()
	//	defer walletQ.AssertExpectations(t)
	//	got, status := get("do@exist.com")
	//	if status != 200 {
	//		t.Fatalf("expected 200 got %d", status)
	//	}
	//	if got.Salt != wallet.Salt {
	//		t.Fatalf("expected %s got %s", wallet.Salt, got.Salt)
	//	}
	//})
}
