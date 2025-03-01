package codec_json

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/cockroachdb/errors"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
	"go.opentelemetry.io/otel/trace"

	eventing "github.com/sweetloveinyourheart/exploding-kittens/pkg/domain-eventing"
	"github.com/sweetloveinyourheart/exploding-kittens/pkg/domain-eventing/common"
)

// EventCodec is a codec for marshaling and unmarshaling events
// to and from bytes in JSON format.
type EventCodec struct{}

// MarshalEvent marshals an event into bytes in JSON format.
func (c *EventCodec) MarshalEvent(ctx context.Context, event common.Event) ([]byte, error) {
	e := evt{
		EventType:     event.EventType(),
		Timestamp:     event.Timestamp(),
		AggregateType: event.AggregateType(),
		AggregateID:   event.AggregateID(),
		Version:       event.Version(),
		Metadata:      event.Metadata(),
		Context:       eventing.MarshalContext(ctx),
	}

	span := trace.SpanFromContext(ctx)
	span.SetAttributes(semconv.MessagingMessagePayloadSizeBytes(len(e.RawData)))

	// Marshal event data if there is any.
	if event.Data() != nil {
		if unregistered, ok := event.(common.UnregisteredEvent); ok && unregistered.Unregistered() {
			e.RawData = event.Data().([]byte)
		} else {
			var err error
			if e.RawData, err = json.Marshal(event.Data()); err != nil {
				return nil, errors.WithStack(fmt.Errorf("could not marshal event data: %w", err))
			}
		}
	}

	// Marshal the event (using JSON for now).
	b, err := json.Marshal(e)
	if err != nil {
		return nil, errors.WithStack(fmt.Errorf("could not marshal event: %w", err))
	}

	return b, nil
}

// UnmarshalEvent unmarshals an event from bytes in JSON format.
func (c *EventCodec) UnmarshalEvent(ctx context.Context, b []byte, options ...eventing.EventOption) (common.Event, context.Context, error) {
	// Decode the raw JSON event data.
	var e evt
	if err := json.Unmarshal(b, &e); err != nil {
		return nil, nil, errors.WithStack(fmt.Errorf("could not unmarshal event: %w", err))
	}

	// Build the event.

	opts := []eventing.EventOption{
		eventing.ForAggregate(
			e.AggregateType,
			e.AggregateID,
			e.Version,
		),
		eventing.WithMetadata(e.Metadata),
	}
	opts = append(opts, options...)

	event, err := eventing.NewEventFromRaw(
		ctx,
		e.EventType,
		e.RawData,
		e.Timestamp,
		opts...,
	)
	if err != nil {
		return nil, nil, errors.WithStack(fmt.Errorf("could not create event: %w", err))
	}

	// Unmarshal the context.
	ctx = eventing.UnmarshalContext(ctx, e.Context)

	return event, ctx, nil
}

// evt is the internal event used on the wire only.
type evt struct {
	EventType     common.EventType       `json:"event_type"`
	RawData       json.RawMessage        `json:"data,omitempty"`
	Timestamp     time.Time              `json:"timestamp"`
	AggregateType common.AggregateType   `json:"aggregate_type"`
	AggregateID   string                 `json:"aggregate_id"`
	Version       uint64                 `json:"version"`
	Metadata      map[string]interface{} `json:"metadata"`
	Context       map[string]interface{} `json:"context"`
}
