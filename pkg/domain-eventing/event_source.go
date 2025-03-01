package eventing

import "github.com/sweetloveinyourheart/exploding-kittens/pkg/domain-eventing/common"

// EventSource is a source of events, used for getting events for handling,
// storing, publishing etc. Mostly used in the aggregate stores.
type EventSource interface {
	// UncommittedEvents returns events that are not committed to the event store,
	// or handled in other ways (depending on the caller).
	UncommittedEvents() []common.Event
	// ClearUncommittedEvents clears uncommitted events, used after they have been
	// committed to the event store or handled in other ways (depending on the caller).
	ClearUncommittedEvents()
}
