package aggregatestore

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/primus/primus/doma"
)

type AggregateStore struct {
	store doma.EventStore
}

var (
	// ErrInvalidEventStore is when a dispatcher is created with a nil event store.
	ErrInvalidEventStore = errors.New("invalid event store")
	// ErrAggregateNotVersioned is when the aggregate does not implement the VersionedAggregate interface.
	ErrAggregateNotVersioned = errors.New("aggregate is not versioned")
	// ErrMismatchedEventType occurs when loaded events from ID does not match aggregate type.
	ErrMismatchedEventType = errors.New("mismatched event type and aggregate type")
)

func NewAggregateStore(store doma.EventStore) (*AggregateStore, error) {

	if store == nil {
		return nil, ErrInvalidEventStore
	}

	return &AggregateStore{
		store: store,
	}, nil
}

func (r *AggregateStore) Load(ctx context.Context, aggregateType doma.AggregateType, id uuid.UUID) (doma.Aggregate, error) {

	agg, err := doma.CreateAggregate(aggregateType, id)
	if err != nil {
		return nil, err
	}

	a, ok := agg.(VersionedAggregate)
	if !ok {
		return nil, ErrAggregateNotVersioned
	}

	events, err := r.store.Load(ctx, a.EntityID())
	if err != nil && !errors.Is(err, doma.ErrAggregateNotFound) {
		return nil, err
	}

	if err := r.applyEvents(ctx, a, events); err != nil {
		return nil, err
	}

	return agg, nil
}

func (r *AggregateStore) Save(ctx context.Context, agg doma.Aggregate) error {
	a, ok := agg.(VersionedAggregate)
	if !ok {
		return ErrAggregateNotVersioned
	}
	events := a.UncommittedEvents()
	if len(events) == 0 {
		return nil
	}

	if err := r.store.Save(ctx, events, a.AggregateVersion()); err != nil {
		return err
	}
	a.ClearUncommittedEvents()

	if err := r.applyEvents(ctx, a, events); err != nil {
		return err
	}
	return nil

}

func (r *AggregateStore) applyEvents(ctx context.Context, a VersionedAggregate, events []doma.Event) error {

	for _, event := range events {
		if event.AggregateType() != a.AggregateType() {
			return ErrMismatchedEventType
		}

		if err := a.ApplyEvent(ctx, event); err != nil {
			return fmt.Errorf("could not apply event %s: %w", event, err)
		}
		a.SetAggregateVersion(event.Version())
	}

	return nil
}
