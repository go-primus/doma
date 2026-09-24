package kvstore

import (
	"errors"

	"github.com/go-primus/doma/internal/registry"
)

type Driver interface {
	registry.Driver
	Open(path string) (KvStore, error)
}

var registryx registry.Registry = registry.NewRegistry()

// func init() {
// 	registryx = registry.NewRegistry()
// }

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
		return nil, false
	}

	return driver, true
}

func Names() []string {
	return registryx.Names()
}

func Open(driverName, path string) (KvStore, error) {
	d, ok := Get(driverName)
	if !ok {
		return nil, errors.New("unsupported kvstore driver: " + driverName)
	}
	return d.Open(path)
}
