package handlers

import (
	"reflect"
	"testing"

	"net/http"
	"net/http/httptest"

	"bytes"

	"fmt"

	"github.com/go-chi/chi"
	"github.com/stretchr/testify/mock"
	"gitlab.com/swarmfund/api/db2/api"
	"gitlab.com/swarmfund/api/db2/api/mocks"
	"gitlab.com/swarmfund/api/internal/api/middlewares"
	"gitlab.com/swarmfund/api/internal/externalmocks"
	"gitlab.com/swarmfund/api/internal/types"
	doormanTypes "gitlab.com/swarmfund/go/doorman/types"
	"gitlab.com/swarmfund/go/signcontrol"
)

var (
	SignerConstraintAllow = doormanTypes.SignerConstraint(func(_ *http.Request) error {
		return nil
	})

	SignerConstraintNotAllowed = doormanTypes.SignerConstraint(func(_ *http.Request) error {
		return signcontrol.ErrNotAllowed
	})
)

func TestNewCreateBackendRequest(t *testing.T) {
	cases := []struct {
		name     string
		walletID string
		body     string
		err      bool
		expected CreateBackendRequest
	}{
		{
			"valid",
			"aaaa",
			`{"data": {"type": "totp"}}`,
			false,
			CreateBackendRequest{
				WalletID: "aaaa",
				Data: CreateBackendRequestData{
					Type: types.WalletFactorTOTP,
				},
			},
		},
		{
			"missing type",
			"aaaa",
			`{"data": {}}`,
			true,
			CreateBackendRequest{},
		},
		{
			"invalid type",
			"aaaa",
			`{"data": {"type":"foobar"}}`,
			true,
			CreateBackendRequest{},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			r := RequestWithURLParams([]byte(tc.body), map[string]string{
				"wallet-id": tc.walletID,
			})
			got, err := NewCreateBackendRequest(r)
			if err != nil && !tc.err {
				t.Fatalf("expected nil error got %s", err)
			}
			if err == nil && tc.err {
				t.Fatalf("expected error got nil")
			}
			if err == nil && !reflect.DeepEqual(got, tc.expected) {
				t.Fatal("expected %#v got %#v", tc.expected, got)
			}
		})
	}
}

func TestCreateTFABackend(t *testing.T) {
	walletQ := mocks.WalletQI{}
	doorman := externalmocks.Doorman{}
	tfaQ := mocks.TFAQI{}
	walletQ.On("New").Return(&walletQ)
	tfaQ.On("New").Return(&tfaQ)
	router := chi.NewRouter()
	router.Use(middlewares.Ctx(
		CtxWalletQ(&walletQ),
		CtxDoorman(&doorman),
		CtxTFAQ(&tfaQ),
	))
	router.Post("/{wallet-id}", CreateTFABackend)

	ts := httptest.NewServer(router)
	defer ts.Close()

	wallet := api.Wallet{
		WalletId:  "foobar",
		Username:  "fo@ob.ar",
		AccountID: "GDXTLOUXG26JETQVECI3QMVTCG6LBZ7MBQU46QB4PK4G2TMAAZJ32DPD",
	}

	post := func(walletID, body string) int {
		response, err := http.Post(fmt.Sprintf("%s/%s", ts.URL, walletID), "", bytes.NewReader([]byte(body)))
		if err != nil {
			t.Fatal(err)
		}
		defer response.Body.Close()

		return response.StatusCode
	}

	t.Run("not found", func(t *testing.T) {
		walletQ.On("ByWalletID", wallet.WalletId).Return(nil, nil).Once()
		defer walletQ.AssertExpectations(t)
		status := post(wallet.WalletId, `{
			"data": {
				"type": "totp"
			}
		}`)
		if status != 404 {
			t.Fatalf("expected %d got %d", 404, status)
		}
	})

	t.Run("check allowed", func(t *testing.T) {
		walletQ.On("ByWalletID", wallet.WalletId).Return(&wallet, nil).Once()
		doorman.On("SignerOf", wallet.AccountID).Return(SignerConstraintNotAllowed).Once()
		defer walletQ.AssertExpectations(t)
		defer doorman.AssertExpectations(t)
		status := post(wallet.WalletId, `{
			"data": {
				"type": "totp"
			}
		}`)
		if status != 401 {
			t.Fatalf("expected %d got %d", 401, status)
		}
	})

	t.Run("password not allowed", func(t *testing.T) {
		walletQ.On("ByWalletID", wallet.WalletId).Return(&wallet, nil).Once()
		doorman.On("SignerOf", wallet.AccountID).Return(SignerConstraintAllow).Once()
		defer walletQ.AssertExpectations(t)
		defer doorman.AssertExpectations(t)
		status := post(wallet.WalletId, `{
			"data": {
				"type": "password"
			}
		}`)
		if status != 409 {
			t.Fatalf("expected %d got %d", 409, status)
		}
	})

	t.Run("valid", func(t *testing.T) {
		var backendID int64 = 10
		walletQ.On("ByWalletID", wallet.WalletId).Return(&wallet, nil).Once()
		doorman.On("SignerOf", wallet.AccountID).Return(SignerConstraintAllow).Once()
		tfaQ.On("CreateBackend", wallet.WalletId, mock.Anything).Return(&backendID, nil).Once()
		defer walletQ.AssertExpectations(t)
		defer doorman.AssertExpectations(t)
		defer tfaQ.AssertExpectations(t)
		status := post(wallet.WalletId, `{
			"data": {
				"type": "totp"
			}
		}`)
		if status != 201 {
			t.Fatalf("expected %d got %d", 201, status)
		}
	})

	t.Run("multiple conflict", func(t *testing.T) {
		var backendID int64 = 10
		walletQ.On("ByWalletID", wallet.WalletId).Return(&wallet, nil).Twice()
		doorman.On("SignerOf", wallet.AccountID).Return(SignerConstraintAllow).Twice()
		tfaQ.On("CreateBackend", wallet.WalletId, mock.Anything).
			Return(&backendID, nil).Once()
		tfaQ.On("CreateBackend", wallet.WalletId, mock.Anything).
			Return(nil, api.ErrWalletBackendConflict).Once()
		defer walletQ.AssertExpectations(t)
		defer doorman.AssertExpectations(t)
		defer tfaQ.AssertExpectations(t)
		status := post(wallet.WalletId, `{
			"data": {
				"type": "totp"
			}
		}`)
		if status != 201 {
			t.Fatalf("expected %d got %d", 201, status)
		}
		status = post(wallet.WalletId, `{
			"data": {
				"type": "totp"
			}
		}`)
		if status != 409 {
			t.Fatalf("expected %d got %d", 409, status)
		}
	})

	// TODO test response struct for totp
}
