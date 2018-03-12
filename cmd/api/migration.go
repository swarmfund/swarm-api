package main

import (
	"database/sql"

	"github.com/pkg/errors"
	"gitlab.com/swarmfund/api/db2"
)

type Migrator func(*sql.DB, db2.MigrateDir, int) (int, error)

func migrateDB(direction string, count int, dbConnectionURL string, migrator Migrator) (int, error) {
	db, err := sql.Open("postgres", dbConnectionURL)
	if err != nil {
		return 0, errors.Wrap(err, "failed to open database")
	}

	applied, err := migrator(db, db2.MigrateDir(direction), count)
	return applied, errors.Wrap(err, "failed to apply migrations")
}
