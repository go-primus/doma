package db

import (
	"fmt"

	"github.com/go-primus/doma/internal/registry"
	"github.com/go-primus/doma/pkg/db/drivers/mysql"
	"github.com/go-primus/doma/pkg/db/drivers/postgres"
	"github.com/go-primus/doma/pkg/db/drivers/sqlite"
	"gorm.io/gorm"
)

var registryx registry.Registry

func init() {
	registryx = registry.NewRegistry()
	Register(&mysql.Driver{})
	Register(&postgres.Driver{})
	Register(&sqlite.Driver{})
}

func Register(d Driver) {
	registryx.Register(d.Name(), d)
}

func Get(name string) (Driver, bool) {

	d, ok := registryx.Get(name)
	if !ok {
		return nil, false
	}

	driver, ok := d.(Driver)
	if !ok {
		return driver, ok
	}
	return driver, true

}

func Names() []string {
	return registryx.Names()
}

func Open(dbType string, dsn string) (*gorm.DB, error) {
	d, ok := Get(dbType)
	if !ok {
		return nil, fmt.Errorf("unsupported database driver: %s", dbType)
	}
	if err := d.Initialize(dsn); err != nil {
		return nil, err
	}
	return d.Connect(dsn)
}
