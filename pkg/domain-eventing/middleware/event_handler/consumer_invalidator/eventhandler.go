package consumerinvalidator

import (
	"context"

	eventing "github.com/sweetloveinyourheart/exploding-kittens/pkg/domain-eventing"
	"github.com/sweetloveinyourheart/exploding-kittens/pkg/domain-eventing/common"
)

// InvalidatorHandler is used to check for an invalidator observer middleware in a chain.
type InvalidatorHandler interface {
	IsInvalidatorHandler() bool
}

type eventHandler struct {
	eventing.EventHandler
	invalidatorFunc func(ctx context.Context, event common.Event)
	eventMatcher    eventing.EventMatcher
}

func NewMiddleware(eventMatcher eventing.EventMatcher, invalidatorFunc func(ctx context.Context, event common.Event)) func(eventing.EventHandler) eventing.EventHandler {
	return func(h eventing.EventHandler) eventing.EventHandler {
		return &eventHandler{
			EventHandler:    h,
			eventMatcher:    eventMatcher,
			invalidatorFunc: invalidatorFunc,
		}
	}
}

func (h *eventHandler) HandleEvent(ctx context.Context, event common.Event) error {
	err := h.EventHandler.HandleEvent(ctx, event)
	if err != nil {
		return err
	}

	if h.eventMatcher.Match(event) {
		h.invalidatorFunc(ctx, event)
	}
	return nil
}

// IsInvalidatorHandler returns true if the handler is an invalidator.
func (h *eventHandler) IsInvalidatorHandler() bool {
	return true
}

// InnerHandler implements MiddlewareChain
func (h *eventHandler) InnerHandler() eventing.EventHandler {
	return h.EventHandler
}
