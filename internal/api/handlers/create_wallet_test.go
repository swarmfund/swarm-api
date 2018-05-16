package handlers

import (
	"net/http"
	"strings"
	"testing"

	"context"

	"github.com/stretchr/testify/assert"
	"gitlab.com/swarmfund/api/blacklist"
	"gitlab.com/swarmfund/api/internal/api/resources"
)

func TestNewCreateWalletRequest(t *testing.T) {
	cases := []struct {
		name     string
		body     string
		err      bool
		expected CreateWalletRequest
	}{
		{
			"valid",
			`{
				"data": {
					"type": "wallet",
					"id": "388108095960430b80554ac3efb6807a9f286854033aca47f6f466094ab50876",
					"attributes": {
						"account_id": "GD6PPS6VCAN5AN52N2BSUJQTKW2T22AERA6HHI33VBW67T5GCDFWTVET",
						"email": "fo@ob.ar",
						"salt": "wAWBYjv5eSDVZsjY1suFFA==",
						"keychain_data": "eyJJViI6IjN...I6ImdjbSJ9"
					},
					"relationships": {
						"kdf": {
							"data": {
								"type": "kdf",
								"id": "1"
							}
						},
						"recovery": {
							"data": {
								"type": "recovery",
								"id": "recoverywalletid",
								"attributes": {
									"account_id": "GD6PPS6VCAN5AN52N2BSUJQTKW2T22AERA6HHI33VBW67T5GCDFWTVET",
									"keychain_data": "recoverykeychaindata",
									"salt": "recoverysalt=="
								}
							}
						},
						"factor": {
							"data": {
								"type": "password",
								"attributes": {
									"account_id": "GDI54FYDBF2S6GEQJHBLS3HMIEYYKDLVT7YCCI33K5J6B4JTGNP77DEK",
									"keychain_data": "foo..bar",
									"salt": "salt"
								}
							}
						}
					}
				}
			}`,
			false,
			CreateWalletRequest{
				Wallet: resources.Wallet{
					Data: resources.WalletData{
						Type: "wallet",
						ID:   "388108095960430b80554ac3efb6807a9f286854033aca47f6f466094ab50876",
						Attributes: resources.WalletAttributes{
							AccountID:    "GD6PPS6VCAN5AN52N2BSUJQTKW2T22AERA6HHI33VBW67T5GCDFWTVET",
							Email:        "fo@ob.ar",
							Salt:         "wAWBYjv5eSDVZsjY1suFFA==",
							KeychainData: "eyJJViI6IjN...I6ImdjbSJ9",
						},
						Relationships: resources.WalletRelationships{
							Recovery: &resources.PasswordFactor{
								Data: resources.PasswordFactorData{
									Type: "recovery",
									ID:   "recoverywalletid",
									Attributes: resources.PasswordFactorAttributes{
										AccountID:    "GD6PPS6VCAN5AN52N2BSUJQTKW2T22AERA6HHI33VBW67T5GCDFWTVET",
										KeychainData: "recoverykeychaindata",
										Salt:         "recoverysalt==",
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
										Salt:         "salt",
										KeychainData: "foo..bar",
									},
								},
							},
						},
					},
				},
			},
		},
		{
			"missing recovery account id",
			`{
				"data": {
					"type": "wallet",
					"id": "388108095960430b80554ac3efb6807a9f286854033aca47f6f466094ab50876",
					"attributes": {
						"account_id": "GD6PPS6VCAN5AN52N2BSUJQTKW2T22AERA6HHI33VBW67T5GCDFWTVET",
						"email": "fo@ob.ar",
						"salt": "wAWBYjv5eSDVZsjY1suFFA==",
						"keychain_data": "eyJJViI6IjN...I6ImdjbSJ9"
					},
					"relationships": {
						"kdf": {
							"data": {
								"type": "kdf",
								"id": "1"
							}
						},
						"recovery": {
							"data": {
								"type": "recovery",
								"id": "recoverywalletid",
								"attributes": {
									"keychain_data": "recoverykeychaindata",
									"salt": "recoverysalt=="
								}
							}
						},
						"factor": {
							"data": {
								"type": "password",
								"attributes": {
									"account_id": "GDI54FYDBF2S6GEQJHBLS3HMIEYYKDLVT7YCCI33K5J6B4JTGNP77DEK",
									"keychain_data": "foo..bar",
									"salt": "salt"
								}
							}
						}
					}
				}
			}`,
			true,
			CreateWalletRequest{},
		},
		{
			"missing kdf",
			`{
				"data": {
					"type": "wallet",
					"id": "388108095960430b80554ac3efb6807a9f286854033aca47f6f466094ab50876",
					"attributes": {
						"account_id": "GD6PPS6VCAN5AN52N2BSUJQTKW2T22AERA6HHI33VBW67T5GCDFWTVET",
						"email": "fo@ob.ar",
						"salt": "wAWBYjv5eSDVZsjY1suFFA==",
						"keychain_data": "eyJJViI6IjN...I6ImdjbSJ9"
					},
					"relationships": {
						"factor": {
							"data": {
								"type": "password",
								"attributes": {
									"account_id": "GDI54FYDBF2S6GEQJHBLS3HMIEYYKDLVT7YCCI33K5J6B4JTGNP77DEK",
									"keychain_data": "foo..bar",
									"salt": "salt"
								}
							}
						}
					}
				}
			}`,
			true,
			CreateWalletRequest{},
		},
		{
			"missing factor",
			`{
					"data": {
						"type": "wallet",
						"id": "388108095960430b80554ac3efb6807a9f286854033aca47f6f466094ab50876",
						"attributes": {
							"account_id": "GD6PPS6VCAN5AN52N2BSUJQTKW2T22AERA6HHI33VBW67T5GCDFWTVET",
							"email": "fo@ob.ar",
							"salt": "wAWBYjv5eSDVZsjY1suFFA==",
							"keychain_data": "eyJJViI6IjN...I6ImdjbSJ9"
						},
						"relationships": {
							"kdf": {
								"data": {
									"type": "kdf",
									"id": "1"
								}
							}
						}
					}
				}`,
			true,
			CreateWalletRequest{},
		},
		{
			name: "forbidden domain name for email",
			body: `{
				"data": {
					"type": "wallet",
					"id": "388108095960430b80554ac3efb6807a9f286854033aca47f6f466094ab50876",
					"attributes": {
						"account_id": "GD6PPS6VCAN5AN52N2BSUJQTKW2T22AERA6HHI33VBW67T5GCDFWTVET",
						"email": "fo@bar.com",
						"salt": "wAWBYjv5eSDVZsjY1suFFA==",
						"keychain_data": "eyJJViI6IjN...I6ImdjbSJ9"
					},
					"relationships": {
						"kdf": {
							"data": {
								"type": "kdf",
								"id": "1"
							}
						},
						"recovery": {
							"data": {
								"type": "recovery",
								"id": "recoverywalletid",
								"attributes": {
									"account_id": "GD6PPS6VCAN5AN52N2BSUJQTKW2T22AERA6HHI33VBW67T5GCDFWTVET",
									"keychain_data": "recoverykeychaindata",
									"salt": "recoverysalt=="
								}
							}
						},
						"factor": {
							"data": {
								"type": "password",
								"attributes": {
									"account_id": "GDI54FYDBF2S6GEQJHBLS3HMIEYYKDLVT7YCCI33K5J6B4JTGNP77DEK",
									"keychain_data": "foo..bar",
									"salt": "salt"
								}
							}
						}
					}
				}
			}`,
			err:      true,
			expected: CreateWalletRequest{},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			r, err := http.NewRequest("", "", strings.NewReader(tc.body))
			if err != nil {
				t.Fatal(err)
			}

			ctxFunc := CtxDomainApprover(blacklist.NewApprover("bar.com", "foo.com"))
			ctx := context.TODO()
			ctx = ctxFunc(ctx)
			req := r.WithContext(ctx)

			got, err := NewCreateWalletRequest(req)
			if tc.err {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expected, got)
			}
		})
	}
}
