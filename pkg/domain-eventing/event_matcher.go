package eventing

import "github.com/sweetloveinyourheart/exploding-kittens/pkg/domain-eventing/common"

// EventMatcher matches, for example on event types, aggregate types etc.
type EventMatcher interface {
	// Match returns true if the matcher matches an event.
	Match(common.Event) bool
}
