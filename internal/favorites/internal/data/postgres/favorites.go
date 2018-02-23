package postgres

import (
	"github.com/lann/squirrel"
	"github.com/lib/pq"
	"github.com/pkg/errors"
	"gitlab.com/swarmfund/api/db2"
	"gitlab.com/swarmfund/api/internal/favorites/internal/data"
	"gitlab.com/swarmfund/api/internal/favorites/internal/types"
	types2 "gitlab.com/swarmfund/api/internal/types"
)

const (
	favoritesOwnerConstraint   = "favorites_users_fkey"
	favoritiesUniqueConstraint = "favorites_unique_per_user"
	tableFavorites             = "favorites"
	tableFavoritesAliased      = "favorites f"
	tableFavoritesLimit        = 1024
)

var (
	selectFavorite = squirrel.Select(
		"f.id",
		"f.owner",
		"f.type",
		"f.key").
		OrderBy("f.created_at asc").
		From(tableFavoritesAliased)
)

type Favorites struct {
	repo *db2.Repo
	sql  squirrel.SelectBuilder
}

func NewFavorites(repo *db2.Repo) *Favorites {
	return &Favorites{
		repo, selectFavorite,
	}
}

func (q *Favorites) Create(favorite data.Favorite) error {
	stmt := squirrel.Insert(tableFavorites).SetMap(map[string]interface{}{
		"type":  favorite.Type,
		"owner": favorite.Owner,
		"key":   favorite.Key,
	})
	_, err := q.repo.Exec(stmt)
	if err != nil {
		pqerr, ok := errors.Cause(err).(*pq.Error)
		if ok {
			// owner does not exist
			if pqerr.Constraint == favoritesOwnerConstraint {
				return data.ErrOwnerViolated
			}
			// already exists
			if pqerr.Constraint == favoritiesUniqueConstraint {
				// we already have that record, so why throw error
				return nil
			}
		}
	}
	return err
}

func (q *Favorites) Delete(owner types2.Address, id int64) error {
	stmt := squirrel.Delete(tableFavorites).Where("owner = ? and id = ?", owner, id)
	_, err := q.repo.Exec(stmt)
	return err
}

func (q *Favorites) Page(page uint64) data.Favorites {
	q.sql = q.sql.Offset(tableFavoritesLimit * (page - 1)).Limit(tableFavoritesLimit)
	return q
}

func (q *Favorites) ByType(tpe types.FavoriteType) data.Favorites {
	q.sql = q.sql.Where("f.type & ? != 0", tpe)
	return q
}

func (q *Favorites) Select() (result []data.Favorite, err error) {
	err = q.repo.Select(&result, q.sql)
	return result, err
}
