package handlers

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/distributed_lab/ape/apeutil"
	"gitlab.com/swarmfund/api/internal/favorites/internal/resources"
	"gitlab.com/swarmfund/api/internal/favorites/internal/types"
	types2 "gitlab.com/swarmfund/api/internal/types"
)

func TestNewCreateFavoriteRequest(t *testing.T) {
	addr := types2.Address("GDQVAMAYFD44LR7QYS34NEY5TSKBSYRZHVXCLCCD2M2PXSLH4JTQJSQZ")
	cases := []struct {
		name     string
		address  types2.Address
		data     string
		err      bool
		expected CreateFavoriteRequest
	}{
		{
			"valid",
			addr,
			`{
				"data": {
					"type": "asset_pair"
				}
			}`,
			false,
			CreateFavoriteRequest{
				Owner: &addr,
				Favorite: resources.Favorite{
					Data: resources.FavoriteData{
						Type: types.FavoriteTypeAssetPair,
					},
				},
			},
		},
		{
			"valid w/o address",
			"",
			`{
				"data": {
					"type": "asset_pair"
				}
			}`,
			false,
			CreateFavoriteRequest{
				Favorite: resources.Favorite{
					Data: resources.FavoriteData{
						Type: types.FavoriteTypeAssetPair,
					},
				},
			},
		},
		{
			"missing type",
			"GDQVAMAYFD44LR7QYS34NEY5TSKBSYRZHVXCLCCD2M2PXSLH4JTQJSQZ",
			`{
				"data": {
				}
			}`,
			true,
			CreateFavoriteRequest{},
		},
		{
			"invalid type",
			"GDQVAMAYFD44LR7QYS34NEY5TSKBSYRZHVXCLCCD2M2PXSLH4JTQJSQZ",
			`{
				"data": {
					"type": "foobar"
				}
			}`,
			true,
			CreateFavoriteRequest{},
		},
		{
			"invalid address",
			"notvalidaddress",
			`{
				"data": {
					"type": "asset_pair"
				}
			}`,
			true,
			CreateFavoriteRequest{},
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			r := apeutil.RequestWithURLParams([]byte(tc.data), map[string]string{"address": string(tc.address)})
			got, err := NewCreateFavoriteRequest(r)
			if tc.err {
				assert.Error(t, err)
			} else {
				assert.Equal(t, tc.expected, got)
			}
		})
	}

}
