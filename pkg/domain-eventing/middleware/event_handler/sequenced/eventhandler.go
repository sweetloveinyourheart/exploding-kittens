package sequenced

import eventing "github.com/sweetloveinyourheart/exploding-kittens/pkg/domain-eventing"

// SequencedHandler is used to check for an observer middleware in a chain.
type SequencedHandler interface {
	Sequenced() bool
}

type eventHandler struct {
	eventing.EventHandler
}

// NewMiddleware creates a new middleware that can be examined for ephemeral status.
// A handler can be ephemeral if it never cares about events created before the handler,
// or care about events that might occur when the handler is offline.
//
// Such handlers can be for instance handlers that create their initial state on startup
// but needs to update their internal state based on events as they happen.
//
// Marking a handler as ephemeral enables event publishers to optimize operations
// and clean up subscriptions when they are no longer needed.
func NewMiddleware() func(eventing.EventHandler) eventing.EventHandler {
	return func(h eventing.EventHandler) eventing.EventHandler {
		return &eventHandler{
			EventHandler: h,
		}
	}
}

// Sequenced returns true if the consumer should be sequenced.
func (h *eventHandler) Sequenced() bool {
	return true
}

// InnerHandler implements MiddlewareChain
func (h *eventHandler) InnerHandler() eventing.EventHandler {
	return h.EventHandler
}
