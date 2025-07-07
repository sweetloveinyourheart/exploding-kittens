package tracing

import (
	"context"
	"fmt"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"

	eventing "github.com/sweetloveinyourheart/exploding-kittens/pkg/domain-eventing"
	"github.com/sweetloveinyourheart/exploding-kittens/pkg/domain-eventing/common"
)

type EventHandlerMiddlewareOption func(*eventHandler)

// WithSpanKind sets the span kind of the tracing span.
func WithSpanKind(kind trace.SpanKind) EventHandlerMiddlewareOption {
	return func(h *eventHandler) {
		h.spanKind = kind
	}
}

// NewEventHandlerMiddleware returns an event handler middleware that adds tracing spans.
func NewEventHandlerMiddleware(options ...EventHandlerMiddlewareOption) eventing.EventHandlerMiddleware {
	return eventing.EventHandlerMiddleware(func(h eventing.EventHandler) eventing.EventHandler {
		return &eventHandler{EventHandler: h, tracer: otel.Tracer(TracerName), spanKind: trace.SpanKindInternal}
	})
}

type eventHandler struct {
	eventing.EventHandler
	tracer   trace.Tracer
	spanKind trace.SpanKind
}

// InnerHandler implements MiddlewareChain
func (h *eventHandler) InnerHandler() eventing.EventHandler {
	return h.EventHandler
}

// HandleEvent implements the HandleEvent method of the EventHandler.
func (h *eventHandler) HandleEvent(ctx context.Context, event common.Event) error {
	opName := fmt.Sprintf("%s.Event(%s)", h.HandlerType(), event.EventType())

	attrs := []attribute.KeyValue{
		EventType(event.EventType().String()),
		AggregateType(event.AggregateType().String()),
		AggregateID(event.AggregateID()),
		AggregateVersion(int64(event.Version())),
	}
	opts := []trace.SpanStartOption{
		trace.WithAttributes(attrs...),
		trace.WithSpanKind(h.spanKind),
	}
	newCtx, span := h.tracer.Start(ctx, opName, opts...)

	err := h.EventHandler.HandleEvent(newCtx, event)
	if err != nil {
		span.RecordError(err)
	}

	span.End()

	return err
}
