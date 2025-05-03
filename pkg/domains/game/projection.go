package game

import (
	"context"
	"fmt"
	"strings"

	"github.com/cockroachdb/errors"

	"github.com/sweetloveinyourheart/exploding-kittens/pkg/domain-eventing/common"
)

var ErrEventDataTypeMismatch = errors.New("event data type mismatch")

type AllEventsProjector interface {
	HandleGameCreated(ctx context.Context, event common.Event, data *GameCreated, entity *Game) (*Game, error)
	HandleGameInitialized(ctx context.Context, event common.Event, data *GameInitialized, entity *Game) (*Game, error)
	HandleCardPlayed(ctx context.Context, event common.Event, data *CardPlayed, entity *Game) (*Game, error)
}

type eventsProjector interface {
	handleGameCreated(ctx context.Context, event common.Event, entity *Game) (*Game, error)
	handleGameInitialized(ctx context.Context, event common.Event, entity *Game) (*Game, error)
	handleCardPlayed(ctx context.Context, event common.Event, entity *Game) (*Game, error)
}

// GameProjector is an event handler for Projections in the Game domain.
type GameProjector struct {
	handler any
}

var _ eventsProjector = (*GameProjector)(nil)

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
	AfterHandleEvent(ctx context.Context, event common.Event, data any, entity *Game) (*Game, error)
}

// BeforeEntityHandler is the interface that wraps the BeforeHandleEvent method.
type BeforeEntityHandler interface {
	BeforeHandleEvent(ctx context.Context, event common.Event, data any, entity *Game) (*Game, error)
}

// NewGameProjection creates a new GameProjector instance with the supplied handler.
// Handlers that implement any of the Handle* methods will have those methods called when the corresponding event
// is received.
// Handlers that implement the BeforeHandleEvent method will have that method called before the event is processed.
// Handlers that implement the AfterHandleEvent method will have that method called after the event is processed.
// Handlers that wish to handle all events for the domain should implement the AllEventsHandler interface.
func NewGameProjection(handler any) *GameProjector {
	return &GameProjector{
		handler: handler,
	}
}

// Project projects an event onto an entity.
func (p *GameProjector) Project(ctx context.Context, event common.Event, entity *Game) (updateEntity *Game, err error) {
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

		t, err := p.handleGameEvent(ctx, event, entity)
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
func (p *GameProjector) HandleEvent(ctx context.Context, event common.Event) error {
	if !strings.EqualFold(string(event.AggregateType()), string(AggregateType)) {
		return nil
	}

	if handler, ok := p.handler.(BeforeHandler); ok {
		if err := handler.BeforeHandleEvent(ctx, event, event.Data()); err != nil {
			return err
		}
	}

	if _, err := p.handleGameEvent(ctx, event, nil); err != nil {
		return err
	}

	if handler, ok := p.handler.(AfterHandler); ok {
		return handler.AfterHandleEvent(ctx, event, event.Data())
	}

	return nil
}

// handleGameEvent handles game events.
func (p *GameProjector) handleGameEvent(ctx context.Context, event common.Event, entity *Game) (*Game, error) {
	var eventHandler func(context.Context, common.Event, *Game) (*Game, error)

	switch event.EventType() {
	case EventTypeGameCreated:
		eventHandler = p.handleGameCreated
	case EventTypeGameInitialized:
		eventHandler = p.handleGameInitialized
	case EventTypeCardPlayed:
		eventHandler = p.handleCardPlayed
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

// handleGameCreated handles game created events.
func (p *GameProjector) handleGameCreated(ctx context.Context, event common.Event, entity *Game) (*Game, error) {
	data, ok := event.Data().(*GameCreated)
	if !ok {
		return nil, errors.WithStack(errors.Wrap(ErrEventDataTypeMismatch, "handleGameCreated"))
	}

	if handler, ok := p.handler.(interface {
		HandleGameCreated(ctx context.Context, event common.Event, data *GameCreated, entity *Game) (*Game, error)
	}); ok {
		return handler.HandleGameCreated(ctx, event, data, entity)
	}

	if handler, ok := p.handler.(interface {
		HandleGameCreated(ctx context.Context, event common.Event, data *GameCreated) error
	}); ok {
		return entity, handler.HandleGameCreated(ctx, event, data)
	}

	return entity, nil
}

// handleGameInitialized handles game initialized events.
func (p *GameProjector) handleGameInitialized(ctx context.Context, event common.Event, entity *Game) (*Game, error) {
	data, ok := event.Data().(*GameInitialized)
	if !ok {
		return nil, errors.WithStack(errors.Wrap(ErrEventDataTypeMismatch, "handleGameInitialized"))
	}

	if handler, ok := p.handler.(interface {
		HandleGameInitialized(ctx context.Context, event common.Event, data *GameInitialized, entity *Game) (*Game, error)
	}); ok {
		return handler.HandleGameInitialized(ctx, event, data, entity)
	}

	if handler, ok := p.handler.(interface {
		HandleGameInitialized(ctx context.Context, event common.Event, data *GameInitialized) error
	}); ok {
		return entity, handler.HandleGameInitialized(ctx, event, data)
	}

	return entity, nil
}

// handleCardPlayed handles card played events.
func (p *GameProjector) handleCardPlayed(ctx context.Context, event common.Event, entity *Game) (*Game, error) {
	data, ok := event.Data().(*CardPlayed)
	if !ok {
		return nil, errors.WithStack(errors.Wrap(ErrEventDataTypeMismatch, "handleCardPlayed"))
	}

	if handler, ok := p.handler.(interface {
		HandleCardPlayed(ctx context.Context, event common.Event, data *CardPlayed, entity *Game) (*Game, error)
	}); ok {
		return handler.HandleCardPlayed(ctx, event, data, entity)
	}

	if handler, ok := p.handler.(interface {
		HandleCardPlayed(ctx context.Context, event common.Event, data *CardPlayed) error
	}); ok {
		return entity, handler.HandleCardPlayed(ctx, event, data)
	}

	return entity, nil
}
