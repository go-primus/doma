package registry

import (
	"errors"
	"fmt"
	"sync"

	"github.com/go-primus/doma/runtime/core"
	"github.com/google/uuid"
)

var (
	// ErrAggregateNotFound is when no aggregate can be found.
	ErrAggregateNotFound = errors.New("aggregate not found")
	// ErrAggregateNotRegistered is when no aggregate factory was registered.
	ErrAggregateNotRegistered = errors.New("aggregate not registered")
)

// RegisterAggregate registers an aggregate factory for a type. The factory is
// used to create concrete aggregate types when loading from the database.
//
// An example would be:
//
//	RegisterAggregate(func(id UUID) Aggregate { return &MyAggregate{id} })
func RegisterAggregate(factory func(uuid.UUID) core.Aggregate) {
	// Check that the created aggregate matches the registered type.
	// TODO: Explore the use of reflect/gob for creating concrete types without
	// a factory func.
	aggregate := factory(uuid.New())
	if aggregate == nil {
		panic("eventhorizon: created aggregate is nil")
	}

	aggregateType := aggregate.AggregateType()
	if aggregateType == core.AggregateType("") {
		panic("eventhorizon: attempt to register empty aggregate type")
	}

	aggregatesMu.Lock()
	defer aggregatesMu.Unlock()

	if _, ok := aggregates[aggregateType]; ok {
		panic(fmt.Sprintf("eventhorizon: registering duplicate types for %q", aggregateType))
	}

	aggregates[aggregateType] = factory
}

// CreateAggregate creates an aggregate of a type with an ID using the factory
// registered with RegisterAggregate.
func CreateAggregate(aggregateType core.AggregateType, id uuid.UUID) (core.Aggregate, error) {
	aggregatesMu.RLock()
	defer aggregatesMu.RUnlock()

	if factory, ok := aggregates[aggregateType]; ok {
		return factory(id), nil
	}

	return nil, ErrAggregateNotRegistered
}

var aggregates = make(map[core.AggregateType]func(uuid.UUID) core.Aggregate)
var aggregatesMu sync.RWMutex
