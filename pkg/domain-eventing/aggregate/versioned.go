package aggregate

import (
	"context"

	eventing "github.com/sweetloveinyourheart/exploding-kittens/pkg/domain-eventing"
	"github.com/sweetloveinyourheart/exploding-kittens/pkg/domain-eventing/common"
)

// VersionedAggregate is an interface representing a versioned aggregate created
// from events. It receives commands and generates events that are stored.
//
// The aggregate is created/loaded and saved by the Repository inside the
// Dispatcher. A domain specific aggregate can either implement the full interface,
// or more commonly embed *AggregateBase to take care of the common methods.
type VersionedAggregate interface {
	// Provides all the basic aggregate data.
	eventing.Aggregate

	// Provides events to persist and publish from the aggregate.
	eventing.EventSource

	// AggregateVersion returns the version of the aggregate.
	AggregateVersion() uint64
	// SetAggregateVersion sets the version of the aggregate. It should only be
	// called after an event has been successfully applied, often by EH.
	SetAggregateVersion(uint64)

	// ApplyEvent applies an event on the aggregate by setting its values.
	// If there are no errors the version should be incremented by calling
	// IncrementVersion.
	ApplyEvent(context.Context, common.Event) error
}
