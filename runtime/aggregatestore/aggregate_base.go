package aggregatestore

import (
	"time"

	"github.com/go-primus/doma/runtime/core"
	"github.com/google/uuid"
)

type AggregateBase struct {
	id     uuid.UUID
	t      core.AggregateType
	v      int
	events []core.Event
}

// NewAggregateBase creates an aggregate.
func NewAggregateBase(t core.AggregateType, id uuid.UUID) *AggregateBase {
	return &AggregateBase{
		id: id,
		t:  t,
	}
}

// EntityID implements the EntityID method of the eh.Entity and eh.Aggregate interface.
func (a *AggregateBase) EntityID() uuid.UUID {
	return a.id
}

// AggregateType implements the AggregateType method of the eh.Aggregate interface.
func (a *AggregateBase) AggregateType() core.AggregateType {
	return a.t
}

// AggregateVersion implements the AggregateVersion method of the Aggregate interface.
func (a *AggregateBase) AggregateVersion() int {
	return a.v
}

// SetAggregateVersion implements the SetAggregateVersion method of the Aggregate interface.
func (a *AggregateBase) SetAggregateVersion(v int) {
	a.v = v
}

// UncommittedEvents implements the UncommittedEvents method of the eh.EventSource
// interface.
func (a *AggregateBase) UncommittedEvents() []core.Event {
	return a.events
}

// ClearUncommittedEvents implements the ClearUncommittedEvents method of the eh.EventSource
// interface.
func (a *AggregateBase) ClearUncommittedEvents() {
	a.events = nil
}

// AppendEvent appends an event for later retrieval by Events().
func (a *AggregateBase) AppendEvent(t core.EventType, data core.EventData, timestamp time.Time, options ...core.EventOption) core.Event {
	options = append(options, core.ForAggregate(
		a.AggregateType(),
		a.EntityID(),
		a.AggregateVersion()+len(a.events)+1), // TODO: This will probably not work with a global version.
	)
	e := core.NewEvent(t, data, timestamp, options...)
	a.events = append(a.events, e)

	return e
}
