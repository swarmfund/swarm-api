package resources

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWalletAttributes_Validate(t *testing.T) {
	cases := []struct {
		name  string
		valid bool
		data  WalletAttributes
	}{
		{
			"valid",
			true,
			WalletAttributes{
				AccountID:    "GDS5E2I74OO3554HFKNZP57SUJFFBTC7YRUNOELUC5NTPD54BRZJ2T6F",
				Email:        "fo@ob.ar",
				KeychainData: "key..chain",
				Salt:         "salty==",
			},
		},
		{
			"invalid account id",
			false,
			WalletAttributes{
				AccountID:    "GDS5EOELUC5NTPD54BRZJ2T6F",
				Email:        "fo@ob.ar",
				KeychainData: "key..chain",
				Salt:         "salty==",
			},
		},
		{
			"missing email",
			false,
			WalletAttributes{
				AccountID:    "GDS5E2I74OO3554HFKNZP57SUJFFBTC7YRUNOELUC5NTPD54BRZJ2T6F",
				KeychainData: "key..chain",
				Salt:         "salty==",
			},
		},
		{
			"missing email",
			false,
			WalletAttributes{
				AccountID:    "GDS5E2I74OO3554HFKNZP57SUJFFBTC7YRUNOELUC5NTPD54BRZJ2T6F",
				KeychainData: "key..chain",
				Salt:         "salty==",
			},
		},
		{
			"missing keychain",
			false,
			WalletAttributes{
				AccountID: "GDS5E2I74OO3554HFKNZP57SUJFFBTC7YRUNOELUC5NTPD54BRZJ2T6F",
				Email:     "fo@ob.ar",
				Salt:      "salty==",
			},
		},
		{
			"missing salt",
			false,
			WalletAttributes{
				AccountID:    "GDS5E2I74OO3554HFKNZP57SUJFFBTC7YRUNOELUC5NTPD54BRZJ2T6F",
				Email:        "fo@ob.ar",
				KeychainData: "key..chain",
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.data.Validate()
			if tc.valid {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
			}
		})
	}
}

func TestWalletRelationships_Validate(t *testing.T) {
	cases := []struct {
		name  string
		valid bool
		data  WalletRelationships
	}{
		{
			"valid empty",
			true,
			WalletRelationships{},
		},
		{
			"invalid transaction",
			false,
			WalletRelationships{
				Transaction: &Transaction{},
			},
		},
		{
			"invalid kdf",
			false,
			WalletRelationships{
				KDF: &KDFPlain{},
			},
		},
		{
			"invalid factor",
			false,
			WalletRelationships{
				Factor: &Wallet{},
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.data.Validate()
			if tc.valid {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
			}
		})
	}
}
