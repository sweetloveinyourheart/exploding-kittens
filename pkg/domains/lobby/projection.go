package lobby

import (
	"context"
	"fmt"
	"strings"

	"github.com/cockroachdb/errors"

	"github.com/sweetloveinyourheart/exploding-kittens/pkg/domain-eventing/common"
)

var ErrEventDataTypeMismatch = errors.New("event data type mismatch")

type AllEventsProjector interface {
	HandleLobbyCreated(ctx context.Context, event common.Event, data *LobbyCreated, entity *Lobby) (*Lobby, error)
}

type eventsProjector interface {
	handleLobbyCreated(ctx context.Context, event common.Event, entity *Lobby) (*Lobby, error)
}

// LobbyProjector is an event handler for Projections in the Lobby domain.
type LobbyProjector struct {
	handler any
}

var _ eventsProjector = (*LobbyProjector)(nil)

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
	AfterHandleEvent(ctx context.Context, event common.Event, data any, entity *Lobby) (*Lobby, error)
}

// BeforeEntityHandler is the interface that wraps the BeforeHandleEvent method.
type BeforeEntityHandler interface {
	BeforeHandleEvent(ctx context.Context, event common.Event, data any, entity *Lobby) (*Lobby, error)
}

// NewLobbyProjection creates a new LobbyProjector instance with the supplied handler.
// Handlers that implement any of the Handle* methods will have those methods called when the corresponding event
// is received.
// Handlers that implement the BeforeHandleEvent method will have that method called before the event is processed.
// Handlers that implement the AfterHandleEvent method will have that method called after the event is processed.
// Handlers that wish to handle all events for the domain should implement the AllEventsHandler interface.
func NewLobbyProjection(handler any) *LobbyProjector {
	return &LobbyProjector{
		handler: handler,
	}
}

// Project projects an event onto an entity.
func (p *LobbyProjector) Project(ctx context.Context, event common.Event, entity *Lobby) (updateEntity *Lobby, err error) {
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

		t, err := p.handleLobbyEvent(ctx, event, entity)
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
func (p *LobbyProjector) HandleEvent(ctx context.Context, event common.Event) error {
	if !strings.EqualFold(string(event.AggregateType()), string(AggregateType)) {
		return nil
	}

	if handler, ok := p.handler.(BeforeHandler); ok {
		if err := handler.BeforeHandleEvent(ctx, event, event.Data()); err != nil {
			return err
		}
	}

	if _, err := p.handleLobbyEvent(ctx, event, nil); err != nil {
		return err
	}

	if handler, ok := p.handler.(AfterHandler); ok {
		return handler.AfterHandleEvent(ctx, event, event.Data())
	}

	return nil
}

// handleLobbyEvent handles lobby events.
func (p *LobbyProjector) handleLobbyEvent(ctx context.Context, event common.Event, entity *Lobby) (*Lobby, error) {
	var eventHandler func(context.Context, common.Event, *Lobby) (*Lobby, error)

	switch event.EventType() {
	case EventTypeLobbyCreated:
		eventHandler = p.handleLobbyCreated
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

// handleLobbyCreated handles lobby created events.
func (p *LobbyProjector) handleLobbyCreated(ctx context.Context, event common.Event, entity *Lobby) (*Lobby, error) {
	data, ok := event.Data().(*LobbyCreated)
	if !ok {
		return nil, errors.WithStack(errors.Wrap(ErrEventDataTypeMismatch, "handleLobbyCreated"))
	}

	if handler, ok := p.handler.(interface {
		HandleLobbyCreated(ctx context.Context, event common.Event, data *LobbyCreated, entity *Lobby) (*Lobby, error)
	}); ok {
		return handler.HandleLobbyCreated(ctx, event, data, entity)
	}

	if handler, ok := p.handler.(interface {
		HandleLobbyCreated(ctx context.Context, event common.Event, data *LobbyCreated) error
	}); ok {
		return entity, handler.HandleLobbyCreated(ctx, event, data)
	}

	return entity, nil
}
