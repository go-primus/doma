package memory

import (
	"context"
	"fmt"
	"sync"

	"github.com/google/uuid"
	"github.com/jinzhu/copier"

	eh "github.com/primus/primus/doma"
)

// EventStore is an eventhorizon.EventStore where all events are stored in
// memory and not persisted. Useful for testing and experimenting.
type EventStore struct {
	db   map[uuid.UUID]aggregateRecord
	dbMu sync.RWMutex
	// eventHandler eh.EventHandler
}

// NewEventStore creates a new EventStore using memory as storage.
func NewEventStore(options ...Option) (*EventStore, error) {
	s := &EventStore{
		db: map[uuid.UUID]aggregateRecord{},
	}

	for _, option := range options {
		if err := option(s); err != nil {
			return nil, fmt.Errorf("error while applying option: %v", err)
		}
	}

	return s, nil
}

// Option is an option setter used to configure creation.
type Option func(*EventStore) error

// Save implements the Save method of the eventhorizon.EventStore interface.
func (s *EventStore) Save(ctx context.Context, events []eh.Event, originalVersion int) error {
	if err := s.save(ctx, events, originalVersion); err != nil {
		return err
	}

	// Let the optional event handler handle the events. Aborts the transaction
	// in case of error.
	// if s.eventHandler != nil {
	// 	for _, e := range events {
	// 		if err := s.eventHandler.HandleEvent(ctx, e); err != nil {
	// 			return err
	// 		}
	// 	}
	// }

	return nil
}

// This method needs to be separate from the Save() method to not lock the mutex during publishing.
func (s *EventStore) save(ctx context.Context, events []eh.Event, originalVersion int) error {
	s.dbMu.Lock()
	defer s.dbMu.Unlock()

	if len(events) == 0 {
		return eh.ErrMissingEvents
	}

	dbEvents := make([]eh.Event, len(events))
	id := events[0].AggregateID()
	at := events[0].AggregateType()

	// Build all event records, with incrementing versions starting from the
	// original aggregate version.
	for i, event := range events {
		// Only accept events belonging to the same aggregate.
		if event.AggregateID() != id {
			return eh.ErrMismatchedEventAggregateIDs
		}

		if event.AggregateType() != at {
			return eh.ErrMismatchedEventAggregateTypes
		}

		// Only accept events that apply to the correct aggregate version.
		if event.Version() != originalVersion+i+1 {
			return eh.ErrIncorrectEventVersion
		}

		// Create the event record with timestamp.
		e, err := copyEvent(ctx, event)
		if err != nil {
			return fmt.Errorf("could not copy event: %w", err)
		}

		dbEvents[i] = e
	}

	// Either insert a new aggregate or append to an existing.
	if originalVersion == 0 {
		aggregate := aggregateRecord{
			AggregateID: id,
			Version:     len(dbEvents),
			Events:      dbEvents,
		}

		s.db[id] = aggregate
	} else {
		// Increment aggregate version on insert of new event record, and
		// only insert if version of aggregate is matching (ie not changed
		// since loading the aggregate).
		if aggregate, ok := s.db[id]; ok {
			if aggregate.Version != originalVersion {
				return eh.ErrEventConflictFromOtherSave
			}

			aggregate.Version += len(dbEvents)
			aggregate.Events = append(aggregate.Events, dbEvents...)

			s.db[id] = aggregate
		}
	}

	return nil
}

// Load implements the Load method of the eventhorizon.EventStore interface.
func (s *EventStore) Load(ctx context.Context, id uuid.UUID) ([]eh.Event, error) {
	return s.LoadFrom(ctx, id, 1)
}

// LoadFrom loads all events from version for the aggregate id from the store.
func (s *EventStore) LoadFrom(ctx context.Context, id uuid.UUID, version int) ([]eh.Event, error) {
	s.dbMu.RLock()
	defer s.dbMu.RUnlock()

	aggregate, ok := s.db[id]
	if !ok {
		return nil, eh.ErrAggregateNotFound
	}

	events := make([]eh.Event, len(aggregate.Events))

	for i, event := range aggregate.Events {
		if event.Version() < version {
			continue
		}

		e, err := copyEvent(ctx, event)
		if err != nil {
			return nil, fmt.Errorf("could not copy event: %w", err)
		}

		events[i] = e
	}

	return events, nil
}

type aggregateRecord struct {
	AggregateID uuid.UUID
	Version     int
	Events      []eh.Event
	// Snapshot    eh.Aggregate
}

// Close implements the Close method of the eventhorizon.EventStore interface.
func (s *EventStore) Close() error {
	return nil
}

// copyEvent duplicates an event.
func copyEvent(ctx context.Context, event eh.Event) (eh.Event, error) {
	var data eh.EventData

	// Copy data if there is any.
	if event.Data() != nil {
		var err error
		if data, err = eh.CreateEventData(event.EventType()); err != nil {
			return nil, fmt.Errorf("could not create event data: %w", err)
		}

		copier.Copy(data, event.Data())
	}

	return eh.NewEvent(
		event.EventType(),
		data,
		event.Timestamp(),
		eh.ForAggregate(
			event.AggregateType(),
			event.AggregateID(),
			event.Version(),
		),
		eh.WithMetadata(event.Metadata()),
	), nil
}
