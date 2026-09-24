package eventbus

import (
	"fmt"

	"github.com/go-primus/doma/internal/registry"
)

type Driver interface {
	Name() string
	New(...Option) (EventBus, error)
}

var registryx registry.Registry = registry.NewRegistry()

type Option func(opts Options)
type Options map[string]any

func (opts Options) GetString(key string) string {
	val := opts[key]
	str, ok := val.(string)
	if !ok {
		return ""
	}
	return str
}

func init() {
	registryx.Register("emit", &emitDriver{})
	registryx.Register("nats", &natsDriver{})
}

func NewBus(busType string, opts ...Option) (EventBus, error) {

	d, ok := registryx.Get(busType)
	if !ok {
		return nil, fmt.Errorf("unsupported event bus driver: %s", busType)
	}

	driver, ok := d.(Driver)
	if !ok {
		return nil, fmt.Errorf("driver does not implement Driver interface")
	}

	return driver.New(opts...)

}
