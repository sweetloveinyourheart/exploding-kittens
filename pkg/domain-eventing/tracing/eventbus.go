package tracing

import (
	"context"

	"go.opentelemetry.io/otel/trace"

	eventing "github.com/sweetloveinyourheart/exploding-kittens/pkg/domain-eventing"
	"github.com/sweetloveinyourheart/exploding-kittens/pkg/domain-eventing/common"
)

// EventBus is an event bus wrapper that adds tracing.
type EventBus struct {
	eventing.EventBus
	h eventing.EventHandler
}

// NewEventBus creates a EventBus.
func NewEventBus(eventBus eventing.EventBus) *EventBus {
	return &EventBus{
		EventBus: eventBus,
		// Wrap the eventing.EventHandler part of the bus with tracing middleware,
		// set as producer to set the correct tags.
		h: eventing.UseEventHandlerMiddleware(eventBus, NewEventHandlerMiddleware(WithSpanKind(trace.SpanKindProducer))),
	}
}

// HandleEvent implements the HandleEvent method of the eventing.EventHandler interface.
func (b *EventBus) HandleEvent(ctx context.Context, event common.Event) error {
	return b.h.HandleEvent(ctx, event)
}

// AddHandler implements the AddHandler method of the eventing.EventBus interface.
func (b *EventBus) AddHandler(ctx context.Context, m eventing.EventMatcher, h eventing.EventHandler) error {
	if h == nil {
		return eventing.ErrMissingHandler
	}

	// Wrap the handlers in tracing middleware.
	h = eventing.UseEventHandlerMiddleware(h, NewEventHandlerMiddleware(WithSpanKind(trace.SpanKindConsumer)))

	return b.EventBus.AddHandler(ctx, m, h)
}

func (b *EventBus) InnerBus() eventing.EventBus {
	return b.EventBus
}
