package data

import (
	"errors"

	types2 "gitlab.com/swarmfund/api/internal/favorites/internal/types"
	"gitlab.com/swarmfund/api/internal/types"
)

var (
	ErrOwnerViolated = errors.New("owner contraint violated")
)

type Favorite struct {
	ID    int64               `db:"id"`
	Owner types.Address       `db:"owner"`
	Type  types2.FavoriteType `db:"type"`
	Key   string              `db:"key"`
}

type Favorites interface {
	Create(favorite Favorite) error
	Delete(owner types.Address, id int64) error

	// filter methods
	Page(uint64) Favorites
	ByType(types2.FavoriteType) Favorites
	Select() ([]Favorite, error)
}
