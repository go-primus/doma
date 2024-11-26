package doma

import (
	"context"

	"github.com/google/uuid"
)

// AggregateStore is responsible for loading and saving aggregates.
type AggregateStore interface {
	// Load loads the most recent version of an aggregate with a type and id.
	Load(context.Context, AggregateType, uuid.UUID) (Aggregate, error)

	// Save saves the uncommitted events for an aggregate.
	Save(context.Context, Aggregate) error
}
