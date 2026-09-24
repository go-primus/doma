package registry

import "errors"

var (
	ErrNotFound = errors.New("driver not found")
	ErrExists   = errors.New("driver already exists")
)

// Driver creates resources of type T.
type Driver interface {
	Name() string
}

// Registry manages drivers for a specific resource type.
type Registry interface {
	Register(name string, d Driver) error
	Get(name string) (Driver, bool)
	Names() []string
}

// Option is a driver option function.
type Option func(opts Options)

// Options holds driver options.
type Options map[string]any
