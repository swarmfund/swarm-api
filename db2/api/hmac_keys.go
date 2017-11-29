package api

import (
	"database/sql"

	sq "github.com/lann/squirrel"
)

type HMACQI interface {
	GetSecret(public string) (*string, error)
}

func (q *Q) HMAC() HMACQI {
	return &HMACQ{
		parent: q,
	}
}

type HMACQ struct {
	parent *Q
}

func (q *HMACQ) GetSecret(public string) (*string, error) {
	var result string
	stmt := sq.Select("secret").From("hmac_keys").Where("public = ?", public)
	err := q.parent.Get(&result, stmt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &result, err
}
