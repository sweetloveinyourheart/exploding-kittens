package eventing

import (
	"context"

	"github.com/cockroachdb/errors"

	"github.com/sweetloveinyourheart/exploding-kittens/pkg/domain-eventing/common"
)

// ErrMissingEvent is when there is no event to be handled.
var ErrMissingEvent = errors.New("missing event")

// EventHandler is a handler of events. If registered on a bus as a handler only
// one handler of the same type will receive each event. If registered on a bus
// as an observer all handlers of the same type will receive each event.
type EventHandler interface {
	// HandlerType is the type of the handler.
	HandlerType() common.EventHandlerType

	// HandleEvent handles an event.
	HandleEvent(context.Context, common.Event) error
}
