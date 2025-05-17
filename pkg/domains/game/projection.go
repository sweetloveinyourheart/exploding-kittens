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
	HandleTurnStarted(ctx context.Context, event common.Event, data *TurnStarted, entity *Game) (*Game, error)
	HandleTurnFinished(ctx context.Context, event common.Event, data *TurnFinished, entity *Game) (*Game, error)
	HandleTurnReversed(ctx context.Context, event common.Event, data *TurnReversed, entity *Game) (*Game, error)
	HandleCardsPlayed(ctx context.Context, event common.Event, data *CardsPlayed, entity *Game) (*Game, error)
	HandleActionCreated(ctx context.Context, event common.Event, data *ActionCreated, entity *Game) (*Game, error)
	HandleActionExecuted(ctx context.Context, event common.Event, data *ActionExecuted, entity *Game) (*Game, error)
	HandleAffectedPlayerSelected(ctx context.Context, event common.Event, data *AffectedPlayerSelected, entity *Game) (*Game, error)
}

type eventsProjector interface {
	handleGameCreated(ctx context.Context, event common.Event, entity *Game) (*Game, error)
	handleGameInitialized(ctx context.Context, event common.Event, entity *Game) (*Game, error)
	handleTurnStarted(ctx context.Context, event common.Event, entity *Game) (*Game, error)
	handleTurnFinished(ctx context.Context, event common.Event, entity *Game) (*Game, error)
	handleTurnReversed(ctx context.Context, event common.Event, entity *Game) (*Game, error)
	handleCardsPlayed(ctx context.Context, event common.Event, entity *Game) (*Game, error)
	handleActionCreated(ctx context.Context, event common.Event, entity *Game) (*Game, error)
	handleActionExecuted(ctx context.Context, event common.Event, entity *Game) (*Game, error)
	handleAffectedPlayerSelected(ctx context.Context, event common.Event, entity *Game) (*Game, error)
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
	case EventTypeTurnStarted:
		eventHandler = p.handleTurnStarted
	case EventTypeTurnFinished:
		eventHandler = p.handleTurnFinished
	case EventTypeTurnReversed:
		eventHandler = p.handleTurnReversed
	case EventTypeActionCreated:
		eventHandler = p.handleActionCreated
	case EventTypeAffectedPlayerSelected:
		eventHandler = p.handleAffectedPlayerSelected
	case EventTypeActionExecuted:
		eventHandler = p.handleActionExecuted
	case EventTypeCardsPlayed:
		eventHandler = p.handleCardsPlayed
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

// handleTurnStarted handles turn started events.
func (p *GameProjector) handleTurnStarted(ctx context.Context, event common.Event, entity *Game) (*Game, error) {
	data, ok := event.Data().(*TurnStarted)
	if !ok {
		return nil, errors.WithStack(errors.Wrap(ErrEventDataTypeMismatch, "handleTurnStarted"))
	}

	if handler, ok := p.handler.(interface {
		HandleTurnStarted(ctx context.Context, event common.Event, data *TurnStarted, entity *Game) (*Game, error)
	}); ok {
		return handler.HandleTurnStarted(ctx, event, data, entity)
	}

	if handler, ok := p.handler.(interface {
		HandleTurnStarted(ctx context.Context, event common.Event, data *TurnStarted) error
	}); ok {
		return entity, handler.HandleTurnStarted(ctx, event, data)
	}

	return entity, nil
}

// handleTurnFinished handles turn finished events.
func (p *GameProjector) handleTurnFinished(ctx context.Context, event common.Event, entity *Game) (*Game, error) {
	data, ok := event.Data().(*TurnFinished)
	if !ok {
		return nil, errors.WithStack(errors.Wrap(ErrEventDataTypeMismatch, "handleTurnFinished"))
	}

	if handler, ok := p.handler.(interface {
		HandleTurnFinished(ctx context.Context, event common.Event, data *TurnFinished, entity *Game) (*Game, error)
	}); ok {
		return handler.HandleTurnFinished(ctx, event, data, entity)
	}

	if handler, ok := p.handler.(interface {
		HandleTurnFinished(ctx context.Context, event common.Event, data *TurnFinished) error
	}); ok {
		return entity, handler.HandleTurnFinished(ctx, event, data)
	}

	return entity, nil
}

// handleTurnReversed handles turn reversed events.
func (p *GameProjector) handleTurnReversed(ctx context.Context, event common.Event, entity *Game) (*Game, error) {
	data, ok := event.Data().(*TurnReversed)
	if !ok {
		return nil, errors.WithStack(errors.Wrap(ErrEventDataTypeMismatch, "handleTurnReversed"))
	}

	if handler, ok := p.handler.(interface {
		HandleTurnReversed(ctx context.Context, event common.Event, data *TurnReversed, entity *Game) (*Game, error)
	}); ok {
		return handler.HandleTurnReversed(ctx, event, data, entity)
	}

	if handler, ok := p.handler.(interface {
		HandleTurnReversed(ctx context.Context, event common.Event, data *TurnReversed) error
	}); ok {
		return entity, handler.HandleTurnReversed(ctx, event, data)
	}

	return entity, nil
}

// handleCardsPlayed handles cards played events.
func (p *GameProjector) handleCardsPlayed(ctx context.Context, event common.Event, entity *Game) (*Game, error) {
	data, ok := event.Data().(*CardsPlayed)
	if !ok {
		return nil, errors.WithStack(errors.Wrap(ErrEventDataTypeMismatch, "handleCardsPlayed"))
	}

	if handler, ok := p.handler.(interface {
		HandleCardsPlayed(ctx context.Context, event common.Event, data *CardsPlayed, entity *Game) (*Game, error)
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

// handleActionCreated handles action created events.
func (p *GameProjector) handleActionCreated(ctx context.Context, event common.Event, entity *Game) (*Game, error) {
	data, ok := event.Data().(*ActionCreated)
	if !ok {
		return nil, errors.WithStack(errors.Wrap(ErrEventDataTypeMismatch, "handleActionCreated"))
	}

	if handler, ok := p.handler.(interface {
		HandleActionCreated(ctx context.Context, event common.Event, data *ActionCreated, entity *Game) (*Game, error)
	}); ok {
		return handler.HandleActionCreated(ctx, event, data, entity)
	}

	if handler, ok := p.handler.(interface {
		HandleActionCreated(ctx context.Context, event common.Event, data *ActionCreated) error
	}); ok {
		return entity, handler.HandleActionCreated(ctx, event, data)
	}

	return entity, nil
}

// handleAffectedPlayerSelected handles affected player selected events.
func (p *GameProjector) handleAffectedPlayerSelected(ctx context.Context, event common.Event, entity *Game) (*Game, error) {
	data, ok := event.Data().(*AffectedPlayerSelected)
	if !ok {
		return nil, errors.WithStack(errors.Wrap(ErrEventDataTypeMismatch, "handleAffectedPlayerSelected"))
	}

	if handler, ok := p.handler.(interface {
		HandleAffectedPlayerSelected(ctx context.Context, event common.Event, data *AffectedPlayerSelected, entity *Game) (*Game, error)
	}); ok {
		return handler.HandleAffectedPlayerSelected(ctx, event, data, entity)
	}

	if handler, ok := p.handler.(interface {
		HandleAffectedPlayerSelected(ctx context.Context, event common.Event, data *AffectedPlayerSelected) error
	}); ok {
		return entity, handler.HandleAffectedPlayerSelected(ctx, event, data)
	}

	return entity, nil
}

// handleActionExecuted handles action executed events.
func (p *GameProjector) handleActionExecuted(ctx context.Context, event common.Event, entity *Game) (*Game, error) {
	data, ok := event.Data().(*ActionExecuted)
	if !ok {
		return nil, errors.WithStack(errors.Wrap(ErrEventDataTypeMismatch, "handleActionExecuted"))
	}

	if handler, ok := p.handler.(interface {
		HandleActionExecuted(ctx context.Context, event common.Event, data *ActionExecuted, entity *Game) (*Game, error)
	}); ok {
		return handler.HandleActionExecuted(ctx, event, data, entity)
	}

	if handler, ok := p.handler.(interface {
		HandleActionExecuted(ctx context.Context, event common.Event, data *ActionExecuted) error
	}); ok {
		return entity, handler.HandleActionExecuted(ctx, event, data)
	}

	return entity, nil
}
