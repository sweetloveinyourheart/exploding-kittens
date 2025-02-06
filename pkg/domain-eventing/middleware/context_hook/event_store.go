package contexthook

import eventing "github.com/sweetloveinyourheart/exploding-kittens/pkg/domain-eventing"

type eventHandler struct {
	eventing.EventHandler
}

var _ eventing.EventHandler = (*eventHandler)(nil)

// NewMiddleware creates a new middleware that can be examined for options.
func NewMiddleware() func(eventing.EventHandler) eventing.EventHandler {
	return func(h eventing.EventHandler) eventing.EventHandler {
		return &eventHandler{
			EventHandler: h,
		}
	}
}
