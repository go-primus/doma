package aggregatestore

import (
	"context"

	"github.com/go-primus/doma/runtime/core"
	"github.com/google/uuid"
)

// AggregateStore is responsible for loading and saving aggregates.
type AggregateStore interface {
	// Load loads the most recent version of an aggregate with a type and id.
	Load(context.Context, core.AggregateType, uuid.UUID) (core.Aggregate, error)

	// Save saves the uncommitted events for an aggregate.
	Save(context.Context, core.Aggregate) error
}
