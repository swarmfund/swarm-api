package assets

import (
	"database/sql"

	"github.com/pkg/errors"
	"github.com/rubenv/sql-migrate"
	"gitlab.com/swarmfund/api/db2"
)

type MigrationsLoader struct {
	source *migrate.AssetMigrationSource
}

func NewMigrationsLoader() *MigrationsLoader {
	return &MigrationsLoader{}
}

func (l *MigrationsLoader) loadDir(dir string) error {
	l.source = &migrate.AssetMigrationSource{
		Asset:    Asset,
		AssetDir: AssetDir,
		Dir:      dir,
	}
	return nil
}

// Migrate performs schema migration.  Migrations can occur in one of three
// ways:
//
// - up: migrations are performed from the currently installed version upwards.
// If count is 0, all unapplied migrations will be run.
//
// - down: migrations are performed from the current version downard. If count
// is 0, all applied migrations will be run in a downard direction.
//
// - redo: migrations are first ran downard `count` times, and then are ran
// upward back to the current version at the start of the process. If count is
// 0, a count of 1 will be assumed.
func (l *MigrationsLoader) Migrate(db *sql.DB, dir db2.MigrateDir, count int) (int, error) {
	switch dir {
	case db2.MigrateUp:
		return migrate.ExecMax(db, "postgres", l.source, migrate.Up, count)
	case db2.MigrateDown:
		return migrate.ExecMax(db, "postgres", l.source, migrate.Down, count)
	case db2.MigrateRedo:

		if count == 0 {
			count = 1
		}

		down, err := migrate.ExecMax(db, "postgres", l.source, migrate.Down, count)
		if err != nil {
			return down, err
		}

		return migrate.ExecMax(db, "postgres", l.source, migrate.Up, down)
	default:
		return 0, errors.New("Invalid migration direction")
	}
}
