package resources

import (
	"fmt"

	"github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"gitlab.com/swarmfund/api/internal/favorites/internal/data"
	"gitlab.com/swarmfund/api/internal/favorites/internal/types"
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
		validation.Field(&r.Attributes, validation.Required),
	)
}

type FavoriteAttributes struct {
	// Email is a record "owner" for guest-by-email flow
	Email *string `json:"email,omitempty"`
	Key   string  `json:"key"`
}

func (r FavoriteAttributes) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.Email, is.Email),
		validation.Field(&r.Key, validation.Required),
	)
}
