package consumersync

import eventing "github.com/sweetloveinyourheart/exploding-kittens/pkg/domain-eventing"

// ConsumerSyncHandler is used to check for an observer middleware in a chain.
type ConsumerSyncHandler interface {
	SyncEvents() bool
}

type eventHandler struct {
	eventing.EventHandler
}

// NewMiddleware creates a new middleware that can be examined for wait status.
func NewMiddleware() func(eventing.EventHandler) eventing.EventHandler {
	return func(h eventing.EventHandler) eventing.EventHandler {
		return &eventHandler{
			EventHandler: h,
		}
	}
}

// SyncEvents returns true if the consumer process events in order
func (h *eventHandler) SyncEvents() bool {
	return true
}

// InnerHandler implements MiddlewareChain
func (h *eventHandler) InnerHandler() eventing.EventHandler {
	return h.EventHandler
}
