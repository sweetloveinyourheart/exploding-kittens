package desk

import (
	"context"
	"fmt"
	"strings"

	"github.com/cockroachdb/errors"

	"github.com/sweetloveinyourheart/exploding-kittens/pkg/domain-eventing/common"
)

var ErrEventDataTypeMismatch = errors.New("event data type mismatch")

type AllEventsProjector interface {
	HandleDeskCreated(ctx context.Context, event common.Event, data *DeskCreated, entity *Desk) (*Desk, error)
	HandleDeskShuffled(ctx context.Context, event common.Event, data *DeskShuffled, entity *Desk) (*Desk, error)
	HandleCardsDiscarded(ctx context.Context, event common.Event, data *CardsDiscarded, entity *Desk) (*Desk, error)
}

type eventsProjector interface {
	handleDeskCreated(ctx context.Context, event common.Event, entity *Desk) (*Desk, error)
	handleDeskShuffled(ctx context.Context, event common.Event, entity *Desk) (*Desk, error)
	handleCardsDiscarded(ctx context.Context, event common.Event, entity *Desk) (*Desk, error)
}

// DeskProjector is an event handler for Projections in the Desk domain.
type DeskProjector struct {
	handler any
}

var _ eventsProjector = (*DeskProjector)(nil)

// AfterHandler is the interface that wraps the AfterHandleEvent method.
type AfterHandler interface {
	AfterHandleEvent(ctx context.Context, event common.Event, data any) error
}

// BeforeHandler is the interface that wraps the BeforeHandleEvent method.
type BeforeHandler interface {
	BeforeHandleEvent(ctx context.Context, event common.Event, data any) error
}

// AfterEntityHandler is the interface that wraps the AfterHandleEvent method.
type AfterEntityHandler interface {
	AfterHandleEvent(ctx context.Context, event common.Event, data any, entity *Desk) (*Desk, error)
}

// BeforeEntityHandler is the interface that wraps the BeforeHandleEvent method.
type BeforeEntityHandler interface {
	BeforeHandleEvent(ctx context.Context, event common.Event, data any, entity *Desk) (*Desk, error)
}

// NewDeskProjection creates a new DeskProjector instance with the supplied handler.
// Handlers that implement any of the Handle* methods will have those methods called when the corresponding event
// is received.
// Handlers that implement the BeforeHandleEvent method will have that method called before the event is processed.
// Handlers that implement the AfterHandleEvent method will have that method called after the event is processed.
// Handlers that wish to handle all events for the domain should implement the AllEventsHandler interface.
func NewDeskProjection(handler any) *DeskProjector {
	return &DeskProjector{
		handler: handler,
	}
}

// Project projects an event onto an entity.
func (p *DeskProjector) Project(ctx context.Context, event common.Event, entity *Desk) (updateEntity *Desk, err error) {
	if strings.EqualFold(string(event.AggregateType()), string(AggregateType)) {

		if handler, ok := p.handler.(BeforeHandler); ok {
			if err := handler.BeforeHandleEvent(ctx, event, event.Data()); err != nil {
				return nil, err
			}
		}

		if handler, ok := p.handler.(BeforeEntityHandler); ok {
			t, err := handler.BeforeHandleEvent(ctx, event, event.Data(), entity)
			if err != nil {
				return nil, err
			}
			entity = t
		}

		t, err := p.handleDeskEvent(ctx, event, entity)
		if err != nil {
			return nil, err
		}
		entity = t

		if handler, ok := p.handler.(AfterHandler); ok {
			if err := handler.AfterHandleEvent(ctx, event, event.Data()); err != nil {
				return nil, err
			}
		}

		if handler, ok := p.handler.(AfterEntityHandler); ok {
			t, err := handler.AfterHandleEvent(ctx, event, event.Data(), entity)
			if err != nil {
				return nil, err
			}
			entity = t
		}
	}

	return entity, nil
}

// HandleEvent processes the supplied event and implements the eventing.EventHandler interface.
func (p *DeskProjector) HandleEvent(ctx context.Context, event common.Event) error {
	if !strings.EqualFold(string(event.AggregateType()), string(AggregateType)) {
		return nil
	}

	if handler, ok := p.handler.(BeforeHandler); ok {
		if err := handler.BeforeHandleEvent(ctx, event, event.Data()); err != nil {
			return err
		}
	}

	if _, err := p.handleDeskEvent(ctx, event, nil); err != nil {
		return err
	}

	if handler, ok := p.handler.(AfterHandler); ok {
		return handler.AfterHandleEvent(ctx, event, event.Data())
	}

	return nil
}

// handleDeskEvent handles game events.
func (p *DeskProjector) handleDeskEvent(ctx context.Context, event common.Event, entity *Desk) (*Desk, error) {
	var eventHandler func(context.Context, common.Event, *Desk) (*Desk, error)

	switch event.EventType() {
	case EventTypeDeskCreated:
		eventHandler = p.handleDeskCreated
	case EventTypeDeskShuffled:
		eventHandler = p.handleDeskShuffled
	case EventTypeCardsDiscarded:
		eventHandler = p.handleCardsDiscarded
	default:
		if unregistered, ok := event.(common.UnregisteredEvent); !ok || !unregistered.Unregistered() {
			return nil, fmt.Errorf("unknown event type: %s", event.EventType())
		}
	}

	if entity, err := eventHandler(ctx, event, entity); err != nil {
		return nil, err
	} else {
		return entity, nil
	}
}

// handleDeskCreated handles game created events.
func (p *DeskProjector) handleDeskCreated(ctx context.Context, event common.Event, entity *Desk) (*Desk, error) {
	data, ok := event.Data().(*DeskCreated)
	if !ok {
		return nil, errors.WithStack(errors.Wrap(ErrEventDataTypeMismatch, "handleDeskCreated"))
	}

	if handler, ok := p.handler.(interface {
		HandleDeskCreated(ctx context.Context, event common.Event, data *DeskCreated, entity *Desk) (*Desk, error)
	}); ok {
		return handler.HandleDeskCreated(ctx, event, data, entity)
	}

	if handler, ok := p.handler.(interface {
		HandleDeskCreated(ctx context.Context, event common.Event, data *DeskCreated) error
	}); ok {
		return entity, handler.HandleDeskCreated(ctx, event, data)
	}

	return entity, nil
}

// handleDeskShuffled handles game shuffled events.
func (p *DeskProjector) handleDeskShuffled(ctx context.Context, event common.Event, entity *Desk) (*Desk, error) {
	data, ok := event.Data().(*DeskShuffled)
	if !ok {
		return nil, errors.WithStack(errors.Wrap(ErrEventDataTypeMismatch, "handleDeskShuffled"))
	}

	if handler, ok := p.handler.(interface {
		HandleDeskShuffled(ctx context.Context, event common.Event, data *DeskShuffled, entity *Desk) (*Desk, error)
	}); ok {
		return handler.HandleDeskShuffled(ctx, event, data, entity)
	}

	if handler, ok := p.handler.(interface {
		HandleDeskShuffled(ctx context.Context, event common.Event, data *DeskShuffled) error
	}); ok {
		return entity, handler.HandleDeskShuffled(ctx, event, data)
	}

	return entity, nil
}

// handleCardsDiscarded handles cards discarded events.
func (p *DeskProjector) handleCardsDiscarded(ctx context.Context, event common.Event, entity *Desk) (*Desk, error) {
	data, ok := event.Data().(*CardsDiscarded)
	if !ok {
		return nil, errors.WithStack(errors.Wrap(ErrEventDataTypeMismatch, "handleCardsDiscarded"))
	}

	if handler, ok := p.handler.(interface {
		HandleCardsDiscarded(ctx context.Context, event common.Event, data *CardsDiscarded, entity *Desk) (*Desk, error)
	}); ok {
		return handler.HandleCardsDiscarded(ctx, event, data, entity)
	}

	if handler, ok := p.handler.(interface {
		HandleCardsDiscarded(ctx context.Context, event common.Event, data *CardsDiscarded) error
	}); ok {
		return entity, handler.HandleCardsDiscarded(ctx, event, data)
	}

	return entity, nil
}
