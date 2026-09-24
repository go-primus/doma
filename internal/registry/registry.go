package registry

import (
	"sync"
)

type registry struct {
	mu      sync.RWMutex
	drivers map[string]Driver
}

func (r *registry) Register(name string, d Driver) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.drivers[name]; ok {
		return ErrExists
	}
	r.drivers[name] = d
	return nil
}

func (r *registry) Get(name string) (Driver, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	d, ok := r.drivers[name]
	return d, ok
}

func (r *registry) Names() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()
	names := make([]string, 0, len(r.drivers))
	for n := range r.drivers {
		names = append(names, n)
	}
	return names
}

func NewRegistry() Registry {
	return &registry{
		drivers: make(map[string]Driver),
	}
}
