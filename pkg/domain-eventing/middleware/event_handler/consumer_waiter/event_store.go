package consumerwaiter

import eventing "github.com/sweetloveinyourheart/exploding-kittens/pkg/domain-eventing"

type ConsumerWaitHandler interface {
	WaitForCurrentEvents() bool
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

// WaitForCurrentEvents returns true if the consumer should catch up to current events.
func (h *eventHandler) WaitForCurrentEvents() bool {
	return true
}

// InnerHandler implements MiddlewareChain
func (h *eventHandler) InnerHandler() eventing.EventHandler {
	return h.EventHandler
}
