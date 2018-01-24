package handlers

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/swarmfund/api/internal/api/resources"
)

func TestNewChangeWalletIDRequest(t *testing.T) {
	cases := []struct {
		name     string
		err      bool
		walletID string
		body     string
		expected ChangeWalletIDRequest
	}{
		{
			"valid",
			false,
			"wallet-oi",
			`{
				"data": {
					"type": "wallet",
					"id": "isheebk4wi962pjpnvc1pp66r",
					"attributes": {
						"email": "fo@ob.ar",
						"account_id": "GAFB3PIYQEIJRA3U7ZNAI3Q7KZMT7754GJTHCDIIGIYMJ4SCEARYGDTM",
						"salt": "salt==",
						"keychain_data": "foo..bar"
					},
					"relationships": {
						"transaction": {
							"data": {
								"attributes": {
									"envelope": "AAA...AAAA"
								}
							}
						},
						"kdf": {
							"data": {
								"type": "kdf",
								"id": "1"
							}
						},
						"factor": {
							"data": {
								"type": "password",
								"attributes": {
									"account_id": "GDI54FYDBF2S6GEQJHBLS3HMIEYYKDLVT7YCCI33K5J6B4JTGNP77DEK",
									"keychain_data": "foo...bar",
									"salt": "salt"
								}
							}
						}
					}
				}
			}`,
			ChangeWalletIDRequest{
				CurrentWalletID: "wallet-oi",
				Wallet: resources.Wallet{
					Data: resources.WalletData{
						Type: "wallet",
						ID:   "isheebk4wi962pjpnvc1pp66r",
						Attributes: resources.WalletAttributes{
							AccountID:    "GAFB3PIYQEIJRA3U7ZNAI3Q7KZMT7754GJTHCDIIGIYMJ4SCEARYGDTM",
							Email:        "fo@ob.ar",
							KeychainData: "foo..bar",
							Salt:         "salt==",
						},
						Relationships: resources.WalletRelationships{
							Transaction: &resources.Transaction{
								Data: resources.TransactionData{
									Attributes: resources.TransactionAttributes{
										Envelope: "AAA...AAAA",
									},
								},
							},
							KDF: &resources.KDFPlain{
								Data: resources.KDFPlainData{
									Type: "kdf",
									ID:   1,
								},
							},
							Factor: &resources.PasswordFactor{
								Data: resources.PasswordFactorData{
									Type: "password",
									Attributes: resources.PasswordFactorAttributes{
										AccountID:    "GDI54FYDBF2S6GEQJHBLS3HMIEYYKDLVT7YCCI33K5J6B4JTGNP77DEK",
										KeychainData: "foo...bar",
										Salt:         "salt",
									},
								},
							},
						},
					},
				},
			},
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Run(tc.name, func(t *testing.T) {
				// TODO test url params
				r := RequestWithURLParams([]byte(tc.body), map[string]string{
					"wallet-id": tc.walletID,
				})
				got, err := NewChangeWalletIDRequest(r)
				if tc.err {
					assert.Error(t, err)
				} else {
					assert.NoError(t, err)
					assert.Equal(t, tc.expected, got)
				}
			})
		})
	}
}
