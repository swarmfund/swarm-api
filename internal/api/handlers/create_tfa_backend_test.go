package handlers

import (
	"reflect"
	"testing"

	"net/http"
	"net/http/httptest"

	"bytes"

	"fmt"

	"github.com/go-chi/chi"
	"github.com/magiconair/properties/assert"
	"github.com/stretchr/testify/mock"
	"gitlab.com/swarmfund/api/db2/api"
	"gitlab.com/swarmfund/api/db2/api/mocks"
	"gitlab.com/swarmfund/api/internal/api/middlewares"
	"gitlab.com/swarmfund/api/internal/externalmocks"
	"gitlab.com/swarmfund/api/internal/secondfactor"
	"gitlab.com/swarmfund/api/internal/types"
	"gitlab.com/swarmfund/go/doorman"
	"gitlab.com/swarmfund/go/keypair"
	"gitlab.com/swarmfund/go/signcontrol"
)

type TestClient struct {
	t      *testing.T
	ts     *httptest.Server
	signer keypair.KP
}

func Client(t *testing.T, ts *httptest.Server) *TestClient {
	return &TestClient{
		t:  t,
		ts: ts,
	}
}

func (c *TestClient) RandomSigner() *TestClient {
	c.signer, _ = keypair.Random()
	return c
}

func (c *TestClient) Signer(signer keypair.KP) *TestClient {
	c.signer = signer
	return c
}

func (c *TestClient) Do(method, path, body string) *http.Response {
	c.t.Helper()
	request, err := http.NewRequest(method, fmt.Sprintf("%s/%s", c.ts.URL, path), bytes.NewReader([]byte(body)))
	if err != nil {
		c.t.Fatal(err)
	}

	if c.signer != nil {
		if err := signcontrol.SignRequest(request, c.signer); err != nil {
			c.t.Fatal(err)
		}
	}

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		c.t.Fatal(err)
	}
	return response
}

func (c *TestClient) Post(path, body string) *http.Response {
	return c.Do("POST", path, body)
}

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
				t.Fatalf("expected %#v got %#v", tc.expected, got)
			}
		})
	}
}

func TestCreateTFABackend(t *testing.T) {
	walletQ := mocks.WalletQI{}
	walletQ.On("New").Return(&walletQ)

	tfaQ := mocks.TFAQI{}
	tfaQ.On("New").Return(&tfaQ)

	accountQ := externalmocks.AccountQ{}
	doormanM := doorman.New(
		false, &accountQ,
	)

	router := chi.NewRouter()
	router.Use(
		secondfactor.HashMiddleware(),
		middlewares.Ctx(
			CtxWalletQ(&walletQ),
			CtxDoorman(doormanM),
			CtxTFAQ(&tfaQ),
		))
	router.Post("/{wallet-id}", CreateTFABackend)

	ts := httptest.NewServer(router)
	defer ts.Close()

	signer, err := keypair.Random()
	if err != nil {
		t.Fatal(err)
	}

	wallet := api.Wallet{
		WalletId:  "foobar",
		Username:  "fo@ob.ar",
		AccountID: signer.Address(),
	}

	t.Run("not found", func(t *testing.T) {
		walletQ.On("ByWalletID", wallet.WalletId).Return(nil, nil).Once()
		defer walletQ.AssertExpectations(t)
		resp := Client(t, ts).Post(wallet.WalletId, `{
			"data": {
				"type": "totp"
			}
		}`)
		assert.Equal(t, resp.StatusCode, 404)
	})

	t.Run("check not allowed", func(t *testing.T) {
		walletQ.On("ByWalletID", wallet.WalletId).Return(&wallet, nil).Once()
		accountQ.On("Signers", wallet.AccountID).Return(nil, nil).Once()
		defer walletQ.AssertExpectations(t)
		resp := Client(t, ts).RandomSigner().Post(wallet.WalletId, `{
			"data": {
				"type": "totp"
			}
		}`)
		assert.Equal(t, resp.StatusCode, 401)
	})

	t.Run("password not allowed", func(t *testing.T) {
		walletQ.On("ByWalletID", wallet.WalletId).Return(&wallet, nil).Once()
		defer walletQ.AssertExpectations(t)
		resp := Client(t, ts).Signer(signer).Post(wallet.WalletId, `{
			"data": {
				"type": "password"
			}
		}`)
		assert.Equal(t, resp.StatusCode, 409)
	})

	t.Run("valid", func(t *testing.T) {
		var backendID int64 = 10
		walletQ.On("ByWalletID", wallet.WalletId).Return(&wallet, nil).Once()
		tfaQ.On("CreateBackend", wallet.WalletId, mock.Anything).Return(&backendID, nil).Once()
		tfaQ.On("Consume", mock.Anything).Return(true, nil).Once()
		defer walletQ.AssertExpectations(t)
		defer tfaQ.AssertExpectations(t)
		resp := Client(t, ts).Signer(signer).Post(wallet.WalletId, `{
			"data": {
				"type": "totp"
			}
		}`)
		assert.Equal(t, resp.StatusCode, 201)
	})

	t.Run("multiple conflict", func(t *testing.T) {
		walletQ.On("ByWalletID", wallet.WalletId).Return(&wallet, nil).Once()
		tfaQ.On("CreateBackend", wallet.WalletId, mock.Anything).
			Return(nil, api.ErrWalletBackendConflict).Once()
		tfaQ.On("Consume", mock.Anything).Return(true, nil).Once()
		defer walletQ.AssertExpectations(t)
		defer tfaQ.AssertExpectations(t)
		resp := Client(t, ts).Signer(signer).Post(wallet.WalletId, `{
			"data": {
				"type": "totp"
			}
		}`)
		assert.Equal(t, resp.StatusCode, 409)
	})

	// TODO test not verified factor
	// TODO test response struct for totp
}
