package doma

import (
	"errors"
	"fmt"
	"sync"

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
func RegisterAggregate(factory func(uuid.UUID) Aggregate) {
	// Check that the created aggregate matches the registered type.
	// TODO: Explore the use of reflect/gob for creating concrete types without
	// a factory func.
	aggregate := factory(uuid.New())
	if aggregate == nil {
		panic("eventhorizon: created aggregate is nil")
	}

	aggregateType := aggregate.AggregateType()
	if aggregateType == AggregateType("") {
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
func CreateAggregate(aggregateType AggregateType, id uuid.UUID) (Aggregate, error) {
	aggregatesMu.RLock()
	defer aggregatesMu.RUnlock()

	if factory, ok := aggregates[aggregateType]; ok {
		return factory(id), nil
	}

	return nil, ErrAggregateNotRegistered
}

var aggregates = make(map[AggregateType]func(uuid.UUID) Aggregate)
var aggregatesMu sync.RWMutex
