package resources

import (
	"fmt"

	"github.com/go-ozzo/ozzo-validation"
	"gitlab.com/swarmfund/api/internal/favorites/internal/data"
	"gitlab.com/swarmfund/api/internal/favorites/internal/types"
	"gitlab.com/swarmfund/go/xdr"
)

type Favorite struct {
	Data FavoriteData `json:"data"`
}

func NewFavorite(record data.Favorite) Favorite {
	return Favorite{
		Data: FavoriteData{
			Type: record.Type,
			ID:   fmt.Sprintf("%d", record.ID),
			Attributes: FavoriteAttributes{
				Key: record.Key,
			},
		},
	}
}

func (r Favorite) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.Data, validation.Required),
	)
}

type FavoriteData struct {
	Type       types.FavoriteType `json:"type"`
	ID         string             `json:"id"`
	Attributes FavoriteAttributes `json:"attributes"`
}

func (r FavoriteData) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.Type, validation.Required),
	)
}

type FavoriteAttributes struct {
	Key string `json:"key"`
}

func (r FavoriteAttributes) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.Key, validation.Required),
	)

	xdr.REviRequesTyp
	xdr.OperationType()
}
