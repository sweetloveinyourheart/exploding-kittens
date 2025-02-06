package consumerresetter

import eventing "github.com/sweetloveinyourheart/exploding-kittens/pkg/domain-eventing"

// ConsumerResetterHandler is used to check for an observer middleware in a chain.
type ConsumerResetterHandler interface {
	ResetConsumer() bool
}

type eventHandler struct {
	eventing.EventHandler
}

// NewMiddleware creates a new middleware that can be examined for reset status.
func NewMiddleware() func(eventing.EventHandler) eventing.EventHandler {
	return func(h eventing.EventHandler) eventing.EventHandler {
		return &eventHandler{
			EventHandler: h,
		}
	}
}

// ResetConsumer returns true if the consumer should be reset.
func (h *eventHandler) ResetConsumer() bool {
	return true
}

// InnerHandler implements MiddlewareChain
func (h *eventHandler) InnerHandler() eventing.EventHandler {
	return h.EventHandler
}
