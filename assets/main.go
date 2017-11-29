package assets

import (
	"gitlab.com/swarmfund/api/log"
)

//go:generate go-bindata -ignore .+\.go$ -pkg assets -o bindata.go ./...
//go:generate gofmt -w bindata.go

const (
	templatesDir  = "templates"
	enumsDir      = "enums"
	migrationsDir = "migrations"
)

var (
	Templates  *TemplatesLoader
	Enums      *EnumsLoader
	Migrations *MigrationsLoader
)

type AssetFn func(name string) ([]byte, error)
type AssetDirFn func(name string) ([]string, error)

func init() {
	Templates = NewTemplatesLoader()
	if err := Templates.loadDir(templatesDir); err != nil {
		log.WithField("service", "load-templates").WithError(err).Fatal("failed to load templates")
		return
	}
}

func init() {
	Enums = NewEnumsLoader()
	if err := Enums.loadDir(enumsDir); err != nil {
		log.WithField("service", "load-enums").WithError(err).Fatal("failed to load enums")
		return
	}
}

func init() {
	Migrations = NewMigrationsLoader()
	if err := Migrations.loadDir(migrationsDir); err != nil {
		log.WithField("service", "load-migrations").WithError(err).Fatal("failed to load migrations")
		return
	}
}
