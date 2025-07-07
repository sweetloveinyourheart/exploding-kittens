package tracing

import (
	"context"
	"fmt"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	semconv "go.opentelemetry.io/otel/semconv/v1.20.0"
	"go.opentelemetry.io/otel/trace"

	eventing "github.com/sweetloveinyourheart/exploding-kittens/pkg/domain-eventing"
	"github.com/sweetloveinyourheart/exploding-kittens/pkg/domain-eventing/common"
)

// NewCommandHandlerMiddleware returns a new command handler middleware that adds tracing spans.
func NewCommandHandlerMiddleware() eventing.CommandHandlerMiddleware {
	tracer := otel.Tracer(TracerName)
	return eventing.CommandHandlerMiddleware(func(h eventing.CommandHandler) eventing.CommandHandler {
		return eventing.CommandHandlerFunc(func(ctx context.Context, cmd eventing.Command) ([]common.Event, error) {
			attrs := []attribute.KeyValue{
				semconv.MessagingSystem("memory"),
				attribute.String("operation", "command"),
				CommandType(cmd.CommandType().String()),
				AggregateType(cmd.AggregateType().String()),
				AggregateID(cmd.AggregateID()),
			}
			opts := []trace.SpanStartOption{
				trace.WithAttributes(attrs...),
				trace.WithSpanKind(trace.SpanKindInternal),
			}
			newCtx, span := tracer.Start(ctx, fmt.Sprintf("Command(%s)", cmd.CommandType()), opts...)

			result, err := h.HandleCommandEx(newCtx, cmd)
			if err != nil {
				span.RecordError(err)
			}
			span.End()

			return result, err
		})
	})
}
