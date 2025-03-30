package eventstore

import (
	"context"
	"errors"

	"github.com/go-primus/doma/runtime/core"
	"github.com/google/uuid"
)

// EventStore is an interface for an event sourcing event store.
type EventStore interface {
	// Save appends all events in the event stream to the store.
	Save(ctx context.Context, events []core.Event, originalVersion int) error

	// Load loads all events for the aggregate id from the store.
	Load(context.Context, uuid.UUID) ([]core.Event, error)

	// LoadFrom loads all events from version for the aggregate id from the store.
	LoadFrom(ctx context.Context, id uuid.UUID, version int) ([]core.Event, error)

	// Close closes the EventStore.
	Close() error
}

// SnapshotStore is an interface for snapshot store.
type SnapshotStore interface {
	LoadSnapshot(ctx context.Context, id uuid.UUID) (*core.Snapshot, error)
	SaveSnapshot(ctx context.Context, id uuid.UUID, snapshot core.Snapshot) error
}

var (
	// Missing events for save operation.
	ErrMissingEvents = errors.New("missing events")
	// Events in the same save operation is for different aggregate IDs.
	ErrMismatchedEventAggregateIDs = errors.New("mismatched event aggregate IDs")
	// Events in the same save operation is for different aggregate types.
	ErrMismatchedEventAggregateTypes = errors.New("mismatched event aggregate types")
	// Events in the same operation have non-serial versions or is not matching the original version.
	ErrIncorrectEventVersion = errors.New("incorrect event version")
	// Other events has been saved for this aggregate since the operation started.
	ErrEventConflictFromOtherSave = errors.New("event conflict from other save")
	// No matching event could be found (for maintenance operations etc).
	ErrEventNotFound = errors.New("event not found")
)

// EventStoreOperation is the operation done when an error happened.
type EventStoreOperation string

const (
	// Errors during loading of events.
	EventStoreOpLoad = "load"
	// Errors during saving of events.
	EventStoreOpSave = "save"
	// Errors during replacing of events.
	EventStoreOpReplace = "replace"
	// Errors during renaming of event types.
	EventStoreOpRename = "rename"
	// Errors during clearing of the event store.
	EventStoreOpClear = "clear"

	// Errors during loading of snapshot.
	EventStoreOpLoadSnapshot = "load_snapshot"
	// Errors during saving of snapshot.
	EventStoreOpSaveSnapshot = "save_snapshot"
)
