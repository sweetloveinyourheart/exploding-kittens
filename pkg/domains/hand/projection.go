package hand

import (
	"context"
	"fmt"
	"strings"

	"github.com/cockroachdb/errors"

	"github.com/sweetloveinyourheart/exploding-kittens/pkg/domain-eventing/common"
)

var ErrEventDataTypeMismatch = errors.New("event data type mismatch")

type AllEventsProjector interface {
	HandleHandCreated(ctx context.Context, event common.Event, data *HandCreated, entity *Hand) (*Hand, error)
	HandleHandShuffled(ctx context.Context, event common.Event, data *HandShuffled, entity *Hand) (*Hand, error)
	HandleCardsAdded(ctx context.Context, event common.Event, data *CardsAdded, entity *Hand) (*Hand, error)
	HandleCardsRemoved(ctx context.Context, event common.Event, data *CardsRemoved, entity *Hand) (*Hand, error)
	HandleCardStolen(ctx context.Context, event common.Event, data *CardStolen, entity *Hand) (*Hand, error)
	HandleCardsPlayed(ctx context.Context, event common.Event, data *CardsPlayed, entity *Hand) (*Hand, error)
}

type eventsProjector interface {
	handleHandCreated(ctx context.Context, event common.Event, entity *Hand) (*Hand, error)
	handleHandShuffled(ctx context.Context, event common.Event, entity *Hand) (*Hand, error)
	handleCardsAdded(ctx context.Context, event common.Event, entity *Hand) (*Hand, error)
	handleCardsRemoved(ctx context.Context, event common.Event, entity *Hand) (*Hand, error)
	handleCardStolen(ctx context.Context, event common.Event, entity *Hand) (*Hand, error)
	handleCardsPlayed(ctx context.Context, event common.Event, entity *Hand) (*Hand, error)
}

// HandProjector is an event handler for Projections in the Hand domain.
type HandProjector struct {
	handler any
}

var _ eventsProjector = (*HandProjector)(nil)

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
	AfterHandleEvent(ctx context.Context, event common.Event, data any, entity *Hand) (*Hand, error)
}

// BeforeEntityHandler is the interface that wraps the BeforeHandleEvent method.
type BeforeEntityHandler interface {
	BeforeHandleEvent(ctx context.Context, event common.Event, data any, entity *Hand) (*Hand, error)
}

// NewHandProjection creates a new HandProjector instance with the supplied handler.
// Handlers that implement any of the Handle* methods will have those methods called when the corresponding event
// is received.
// Handlers that implement the BeforeHandleEvent method will have that method called before the event is processed.
// Handlers that implement the AfterHandleEvent method will have that method called after the event is processed.
// Handlers that wish to handle all events for the domain should implement the AllEventsHandler interface.
func NewHandProjection(handler any) *HandProjector {
	return &HandProjector{
		handler: handler,
	}
}

// Project projects an event onto an entity.
func (p *HandProjector) Project(ctx context.Context, event common.Event, entity *Hand) (updateEntity *Hand, err error) {
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

		t, err := p.handleHandEvent(ctx, event, entity)
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
func (p *HandProjector) HandleEvent(ctx context.Context, event common.Event) error {
	if !strings.EqualFold(string(event.AggregateType()), string(AggregateType)) {
		return nil
	}

	if handler, ok := p.handler.(BeforeHandler); ok {
		if err := handler.BeforeHandleEvent(ctx, event, event.Data()); err != nil {
			return err
		}
	}

	if _, err := p.handleHandEvent(ctx, event, nil); err != nil {
		return err
	}

	if handler, ok := p.handler.(AfterHandler); ok {
		return handler.AfterHandleEvent(ctx, event, event.Data())
	}

	return nil
}

// handleHandEvent handles game events.
func (p *HandProjector) handleHandEvent(ctx context.Context, event common.Event, entity *Hand) (*Hand, error) {
	var eventHandler func(context.Context, common.Event, *Hand) (*Hand, error)

	switch event.EventType() {
	case EventTypeHandCreated:
		eventHandler = p.handleHandCreated
	case EventTypeCardsPlayed:
		eventHandler = p.handleCardsPlayed
	case EventTypeHandShuffled:
		eventHandler = p.handleHandShuffled
	case EventTypeCardsAdded:
		eventHandler = p.handleCardsAdded
	case EventTypeCardsRemoved:
		eventHandler = p.handleCardsRemoved
	case EventTypeCardStolen:
		eventHandler = p.handleCardStolen
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

// handleHandCreated handles game created events.
func (p *HandProjector) handleHandCreated(ctx context.Context, event common.Event, entity *Hand) (*Hand, error) {
	data, ok := event.Data().(*HandCreated)
	if !ok {
		return nil, errors.WithStack(errors.Wrap(ErrEventDataTypeMismatch, "handleHandCreated"))
	}

	if handler, ok := p.handler.(interface {
		HandleHandCreated(ctx context.Context, event common.Event, data *HandCreated, entity *Hand) (*Hand, error)
	}); ok {
		return handler.HandleHandCreated(ctx, event, data, entity)
	}

	if handler, ok := p.handler.(interface {
		HandleHandCreated(ctx context.Context, event common.Event, data *HandCreated) error
	}); ok {
		return entity, handler.HandleHandCreated(ctx, event, data)
	}

	return entity, nil
}

// handleCardsPlayed handles game cards played events.
func (p *HandProjector) handleCardsPlayed(ctx context.Context, event common.Event, entity *Hand) (*Hand, error) {
	data, ok := event.Data().(*CardsPlayed)
	if !ok {
		return nil, errors.WithStack(errors.Wrap(ErrEventDataTypeMismatch, "handleCardsPlayed"))
	}

	if handler, ok := p.handler.(interface {
		HandleCardsPlayed(ctx context.Context, event common.Event, data *CardsPlayed, entity *Hand) (*Hand, error)
	}); ok {
		return handler.HandleCardsPlayed(ctx, event, data, entity)
	}

	if handler, ok := p.handler.(interface {
		HandleCardsPlayed(ctx context.Context, event common.Event, data *CardsPlayed) error
	}); ok {
		return entity, handler.HandleCardsPlayed(ctx, event, data)
	}

	return entity, nil
}

// handleHandShuffled handles game shuffled events.
func (p *HandProjector) handleHandShuffled(ctx context.Context, event common.Event, entity *Hand) (*Hand, error) {
	data, ok := event.Data().(*HandShuffled)
	if !ok {
		return nil, errors.WithStack(errors.Wrap(ErrEventDataTypeMismatch, "handleHandShuffled"))
	}

	if handler, ok := p.handler.(interface {
		HandleHandShuffled(ctx context.Context, event common.Event, data *HandShuffled, entity *Hand) (*Hand, error)
	}); ok {
		return handler.HandleHandShuffled(ctx, event, data, entity)
	}

	if handler, ok := p.handler.(interface {
		HandleHandShuffled(ctx context.Context, event common.Event, data *HandShuffled) error
	}); ok {
		return entity, handler.HandleHandShuffled(ctx, event, data)
	}

	return entity, nil
}

// handleCardsAdded handles game cards added events.
func (p *HandProjector) handleCardsAdded(ctx context.Context, event common.Event, entity *Hand) (*Hand, error) {
	data, ok := event.Data().(*CardsAdded)
	if !ok {
		return nil, errors.WithStack(errors.Wrap(ErrEventDataTypeMismatch, "handleCardsAdded"))
	}

	if handler, ok := p.handler.(interface {
		HandleCardsAdded(ctx context.Context, event common.Event, data *CardsAdded, entity *Hand) (*Hand, error)
	}); ok {
		return handler.HandleCardsAdded(ctx, event, data, entity)
	}

	if handler, ok := p.handler.(interface {
		HandleCardsAdded(ctx context.Context, event common.Event, data *CardsAdded) error
	}); ok {
		return entity, handler.HandleCardsAdded(ctx, event, data)
	}

	return entity, nil
}

// handleCardsRemoved handles game cards removed events.
func (p *HandProjector) handleCardsRemoved(ctx context.Context, event common.Event, entity *Hand) (*Hand, error) {
	data, ok := event.Data().(*CardsRemoved)
	if !ok {
		return nil, errors.WithStack(errors.Wrap(ErrEventDataTypeMismatch, "handleCardsRemoved"))
	}

	if handler, ok := p.handler.(interface {
		HandleCardsRemoved(ctx context.Context, event common.Event, data *CardsRemoved, entity *Hand) (*Hand, error)
	}); ok {
		return handler.HandleCardsRemoved(ctx, event, data, entity)
	}

	if handler, ok := p.handler.(interface {
		HandleCardsRemoved(ctx context.Context, event common.Event, data *CardsRemoved) error
	}); ok {
		return entity, handler.HandleCardsRemoved(ctx, event, data)
	}

	return entity, nil
}

// handleCardStolen handles game card stolen events.
func (p *HandProjector) handleCardStolen(ctx context.Context, event common.Event, entity *Hand) (*Hand, error) {
	data, ok := event.Data().(*CardStolen)
	if !ok {
		return nil, errors.WithStack(errors.Wrap(ErrEventDataTypeMismatch, "handleCardStolen"))
	}

	if handler, ok := p.handler.(interface {
		HandleCardStolen(ctx context.Context, event common.Event, data *CardStolen, entity *Hand) (*Hand, error)
	}); ok {
		return handler.HandleCardStolen(ctx, event, data, entity)
	}

	if handler, ok := p.handler.(interface {
		HandleCardStolen(ctx context.Context, event common.Event, data *CardStolen) error
	}); ok {
		return entity, handler.HandleCardStolen(ctx, event, data)
	}

	return entity, nil
}
