package ephemeral

import eventing "github.com/sweetloveinyourheart/exploding-kittens/pkg/domain-eventing"

// EphemeralHandler is used to check for an ephemeral observer middleware in a chain.
type EphemeralHandler interface {
	IsEphemeralHandler() bool
}

type eventHandler struct {
	eventing.EventHandler
}

// NewMiddleware creates a new middleware that can be examined for ephemeral status.
// A handler can be ephemeral if it never cares about events created before the handler
// was created or about events that occur when the handler is offline.
//
// Such handlers can be for instance handlers that create their initial state on startup
// but need to update their internal state based on events as they happen.
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

// IsEphemeralHandler returns true if the handler should be ephemeral if possible.
func (h *eventHandler) IsEphemeralHandler() bool {
	return true
}

// InnerHandler implements MiddlewareChain
func (h *eventHandler) InnerHandler() eventing.EventHandler {
	return h.EventHandler
}
