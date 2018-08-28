package main

import (
	"database/sql"

	"github.com/pkg/errors"
	"gitlab.com/swarmfund/api/db2"
)

type Migrator func(*sql.DB, db2.MigrateDir, int) (int, error)

func migrateDB(direction string, count int, db *sql.DB, migrator Migrator) (int, error) {
	applied, err := migrator(db, db2.MigrateDir(direction), count)
	return applied, errors.Wrap(err, "failed to apply migrations")
}
