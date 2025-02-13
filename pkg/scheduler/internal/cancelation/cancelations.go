package cancelation

import (
	"context"
	"sync"
)

type Cancelations struct {
	mu          sync.Mutex
	cancelFuncs map[string]context.CancelFunc
}

func NewCancelations() *Cancelations {
	return &Cancelations{

		cancelFuncs: make(map[string]context.CancelFunc),
	}
}

func (c *Cancelations) Add(id string, fn context.CancelFunc) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.cancelFuncs[id] = fn
}

func (c *Cancelations) Delete(id string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	delete(c.cancelFuncs, id)
}

func (c *Cancelations) Get(id string) (fn context.CancelFunc, ok bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	fn, ok = c.cancelFuncs[id]
	return fn, ok
}
