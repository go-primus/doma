package doma

// EventSource is a source of events, used for getting events for handling,
// storing, publishing etc. Mostly used in the aggregate stores.
type EventSource interface {
	// UncommittedEvents returns events that are not committed to the event store,
	// or handeled in other ways (depending on the caller).
	UncommittedEvents() []Event
	// ClearUncommittedEvents clears uncommitted events, used after they have been
	// committed to the event store or handled in other ways (depending on the caller).
	ClearUncommittedEvents()
}
