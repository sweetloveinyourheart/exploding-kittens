package tracing

import (
	"context"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"

	eventing "github.com/sweetloveinyourheart/exploding-kittens/pkg/domain-eventing"
	"github.com/sweetloveinyourheart/exploding-kittens/pkg/domain-eventing/common"
)

// EventStore is an eventing.EventStore that adds tracing with Open Tracing.
type EventStore struct {
	eventing.EventStore
	tracer trace.Tracer
}

// NewEventStore creates a new EventStore.
func NewEventStore(eventStore eventing.EventStore) *EventStore {
	if eventStore == nil {
		return nil
	}

	return &EventStore{
		EventStore: eventStore,
		tracer:     otel.Tracer(TracerName),
	}
}

// Save implements the Save method of the eventing.EventStore interface.
func (s *EventStore) Save(ctx context.Context, events []common.Event, originalVersion uint64) error {
	opName := "EventStore.Save"

	opts := []trace.SpanStartOption{
		trace.WithSpanKind(trace.SpanKindInternal),
	}
	if len(events) > 0 {
		attrs := []attribute.KeyValue{
			EventType(events[0].EventType().String()),
			AggregateType(events[0].AggregateType().String()),
			AggregateID(events[0].AggregateID()),
			AggregateVersion(int64(events[0].Version())),
		}
		opts = append(opts, trace.WithAttributes(attrs...))
	}

	ctx, span := s.tracer.Start(ctx, opName, opts...)

	err := s.EventStore.Save(ctx, events, originalVersion)
	if err != nil {
		span.RecordError(err)
	}

	span.End()

	return err
}

// Load implements the Load method of the eventing.EventStore interface.
func (s *EventStore) Load(ctx context.Context, id string) ([]common.Event, error) {
	opName := "EventStore.Load"

	opts := []trace.SpanStartOption{
		trace.WithSpanKind(trace.SpanKindInternal),
	}

	ctx, span := s.tracer.Start(ctx, opName, opts...)
	events, err := s.EventStore.Load(ctx, id)
	if err != nil {
		span.RecordError(err)
	}

	if len(events) > 0 {
		attrs := []attribute.KeyValue{
			EventType(events[0].EventType().String()),
			AggregateType(events[0].AggregateType().String()),
			AggregateID(events[0].AggregateID()),
			AggregateVersion(int64(events[0].Version())),
		}
		span.SetAttributes(attrs...)
	}

	span.End()

	return events, err
}
