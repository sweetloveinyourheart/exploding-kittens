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
	HandleCardsReceived(ctx context.Context, event common.Event, data *CardsReceived, entity *Hand) (*Hand, error)
	HandleCardsGiven(ctx context.Context, event common.Event, data *CardsGiven, entity *Hand) (*Hand, error)
	HandleCardsPlayed(ctx context.Context, event common.Event, data *CardsPlayed, entity *Hand) (*Hand, error)
}

type eventsProjector interface {
	handleHandCreated(ctx context.Context, event common.Event, entity *Hand) (*Hand, error)
	handleHandShuffled(ctx context.Context, event common.Event, entity *Hand) (*Hand, error)
	handleCardsReceived(ctx context.Context, event common.Event, entity *Hand) (*Hand, error)
	handleCardsGiven(ctx context.Context, event common.Event, entity *Hand) (*Hand, error)
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
	case EventTypeCardsReceived:
		eventHandler = p.handleCardsReceived
	case EventTypeCardsGiven:
		eventHandler = p.handleCardsGiven
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

// handleCardsReceived handles game cards received events.
func (p *HandProjector) handleCardsReceived(ctx context.Context, event common.Event, entity *Hand) (*Hand, error) {
	data, ok := event.Data().(*CardsReceived)
	if !ok {
		return nil, errors.WithStack(errors.Wrap(ErrEventDataTypeMismatch, "handleCardsReceived"))
	}

	if handler, ok := p.handler.(interface {
		HandleCardsReceived(ctx context.Context, event common.Event, data *CardsReceived, entity *Hand) (*Hand, error)
	}); ok {
		return handler.HandleCardsReceived(ctx, event, data, entity)
	}

	if handler, ok := p.handler.(interface {
		HandleCardsReceived(ctx context.Context, event common.Event, data *CardsReceived) error
	}); ok {
		return entity, handler.HandleCardsReceived(ctx, event, data)
	}

	return entity, nil
}

// handleCardsGiven handles game cards given events.
func (p *HandProjector) handleCardsGiven(ctx context.Context, event common.Event, entity *Hand) (*Hand, error) {
	data, ok := event.Data().(*CardsGiven)
	if !ok {
		return nil, errors.WithStack(errors.Wrap(ErrEventDataTypeMismatch, "handleCardsGiven"))
	}

	if handler, ok := p.handler.(interface {
		HandleCardsGiven(ctx context.Context, event common.Event, data *CardsGiven, entity *Hand) (*Hand, error)
	}); ok {
		return handler.HandleCardsGiven(ctx, event, data, entity)
	}

	if handler, ok := p.handler.(interface {
		HandleCardsGiven(ctx context.Context, event common.Event, data *CardsGiven) error
	}); ok {
		return entity, handler.HandleCardsGiven(ctx, event, data)
	}

	return entity, nil
}
