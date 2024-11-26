package doma

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
)

// Snapshotable is an interface for creating and applying a Snapshot record.
type Snapshotable interface {
	CreateSnapshot() *Snapshot
	ApplySnapshot(snapshot *Snapshot)
}

// Snapshot is a recording of the state of an aggregate at a point in time
type Snapshot struct {
	Version       int
	AggregateType AggregateType
	Timestamp     time.Time
	State         interface{}
}

var snapshotDataFactories = make(map[AggregateType]func(uuid2 uuid.UUID) SnapshotData)

type SnapshotData interface{}

var snapshotDataFactoriesMu sync.RWMutex

var ErrSnapshotDataNotRegistered = errors.New("snapshot data not registered")

// RegisterSnapshotData registers an snapshot factory for a type. The factory is
// used to create concrete snapshot state type when unmarshalling.
//
// An example would be:
//
//	RegisterSnapshotData("aggregateType1", func() SnapshotData { return &MySnapshotData{} })
func RegisterSnapshotData(aggregateType AggregateType, factory func(id uuid.UUID) SnapshotData) {
	if aggregateType == AggregateType("") {
		panic("eventhorizon: attempt to register empty aggregate type")
	}

	snapshotDataFactoriesMu.Lock()
	defer snapshotDataFactoriesMu.Unlock()

	if _, ok := snapshotDataFactories[aggregateType]; ok {
		panic(fmt.Sprintf("eventhorizon: registering duplicate types for %q", aggregateType))
	}

	snapshotDataFactories[aggregateType] = factory
}

// CreateSnapshotData create a concrete instance using the registered snapshot factories.
func CreateSnapshotData(AggregateID uuid.UUID, aggregateType AggregateType) (SnapshotData, error) {
	snapshotDataFactoriesMu.RLock()
	defer snapshotDataFactoriesMu.RUnlock()

	if factory, ok := snapshotDataFactories[aggregateType]; ok {
		return factory(AggregateID), nil
	}

	return nil, ErrSnapshotDataNotRegistered
}
