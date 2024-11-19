package memory

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	eh "github.com/primus/primus/doma"
)

// Replace implements the Replace method of the eventhorizon.EventStore interface.
func (s *EventStore) Replace(ctx context.Context, event eh.Event) error {
	id := event.AggregateID()

	s.dbMu.RLock()

	aggregate, ok := s.db[id]
	if !ok {
		s.dbMu.RUnlock()

		return eh.ErrAggregateNotFound
	}
	s.dbMu.RUnlock()

	// Create the event record for the Database.
	e, err := copyEvent(ctx, event)
	if err != nil {
		return fmt.Errorf("could not copy event: %w", err)
	}

	// Find the event to replace.
	idx := -1

	for i, e := range aggregate.Events {
		if e.Version() == event.Version() {
			idx = i

			break
		}
	}

	if idx == -1 {
		return eh.ErrEventNotFound
	}

	// Replace event.
	s.dbMu.Lock()
	defer s.dbMu.Unlock()

	aggregate.Events[idx] = e

	return nil
}

// RenameEvent implements the RenameEvent method of the eventhorizon.EventStore interface.
func (s *EventStore) RenameEvent(ctx context.Context, from, to eh.EventType) error {
	s.dbMu.Lock()
	defer s.dbMu.Unlock()

	updated := map[uuid.UUID]aggregateRecord{}

	for id, aggregate := range s.db {
		events := make([]eh.Event, len(aggregate.Events))

		for i, e := range aggregate.Events {
			if e.EventType() == from {
				// Rename any matching event by duplicating.
				events[i] = eh.NewEvent(
					to,
					e.Data(),
					e.Timestamp(),
					eh.ForAggregate(
						e.AggregateType(),
						e.AggregateID(),
						e.Version(),
					),
					eh.WithMetadata(e.Metadata()),
				)
			}
		}

		aggregate.Events = events

		updated[id] = aggregate
	}

	for id, aggregate := range updated {
		s.db[id] = aggregate
	}

	return nil
}
