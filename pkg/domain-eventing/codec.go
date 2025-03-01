package eventing

import (
	"context"

	"github.com/sweetloveinyourheart/exploding-kittens/pkg/domain-eventing/common"
)

// EventCodec is a codec for marshaling and unmarshaling events to and from bytes.
type EventCodec interface {
	// MarshalEvent marshals an event and the supported parts of context into bytes.
	MarshalEvent(context.Context, common.Event) ([]byte, error)
	// UnmarshalEvent unmarshals an event and supported parts of context from bytes.
	UnmarshalEvent(context.Context, []byte, ...EventOption) (common.Event, context.Context, error)
}

// CommandCodec is a codec for marshaling and unmarshaling commands to and from bytes.
type CommandCodec interface {
	// MarshalCommand marshals a command and the supported parts of context into bytes.
	MarshalCommand(context.Context, Command) ([]byte, error)
	// UnmarshalCommand unmarshals a command and supported parts of context from bytes.
	UnmarshalCommand(context.Context, []byte) (Command, context.Context, error)
}
