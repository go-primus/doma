package core

import "github.com/google/uuid"

type Entity interface {
	EntityID() uuid.UUID
}

type Versionable interface {
	AggregateVersion() int
}
