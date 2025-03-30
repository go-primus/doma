package aggregatestore

import (
	"context"

	"github.com/go-primus/doma/runtime/core"
	"github.com/go-primus/doma/runtime/eventstore"
)

// VersionedAggregate is an interface representing a versioned aggregate created
// from events. It receives commands and generates events that are stored.
//
// The aggregate is created/loaded and saved by the Repository inside the
// Dispatcher. A domain specific aggregate can either implement the full interface,
// or more commonly embed *AggregateBase to take care of the common methods.
type VersionedAggregate interface {
	// Provides all the basic aggregate data.
	core.Aggregate

	// Provides events to persist and publish from the aggregate.
	eventstore.EventSource

	// AggregateVersion returns the version of the aggregate.
	AggregateVersion() int
	// SetAggregateVersion sets the version of the aggregate. It should only be
	// called after an event has been successfully applied, often by EH.
	SetAggregateVersion(int)

	// ApplyEvent applies an event on the aggregate by setting its values.
	// If there are no errors the version should be incremented by calling
	// IncrementVersion.
	ApplyEvent(context.Context, core.Event) error
}
