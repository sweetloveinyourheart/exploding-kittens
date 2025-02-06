package bus

import (
	"context"
	"sync"

	"github.com/cockroachdb/errors"
	eventing "github.com/sweetloveinyourheart/exploding-kittens/pkg/domain-eventing"
	"github.com/sweetloveinyourheart/exploding-kittens/pkg/domain-eventing/common"
)

var (
	// ErrHandlerAlreadySet is when SetHandler fails because a handler is already registered for a
	// command type.
	ErrHandlerAlreadySet = errors.New("handler is already set")
	// ErrHandlerNotFound is when HandleCommand or HandleCommandEx fails because no handler is
	// registered for the given command type.
	ErrHandlerNotFound = errors.New("no handlers for command")
)

// CommandHandler in the domain-eventing package is an eventing.CommandHandler that handles commands
// by routing them to the other eventing.CommandHandlers that are registered in it.
type CommandHandler struct {
	handlers   map[common.CommandType]eventing.CommandHandler
	handlersMu sync.RWMutex
}

// NewCommandHandler creates a CommandHandler with no eventing.CommandHandlers registered. Call
// SetHandler() to add handlers to it.
func NewCommandHandler() *CommandHandler {
	return &CommandHandler{
		handlers: make(map[common.CommandType]eventing.CommandHandler),
	}
}

// HandleCommand handles a command with a handler capable of handling it.
func (h *CommandHandler) HandleCommand(ctx context.Context, cmd eventing.Command) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	if err := eventing.CheckCommand(cmd); err != nil {
		return err
	}

	h.handlersMu.RLock()
	defer h.handlersMu.RUnlock()

	if handler, ok := h.handlers[cmd.CommandType()]; ok {
		return handler.HandleCommand(ctx, cmd)
	}

	return ErrHandlerNotFound
}

// HandleCommandEx handles a command with a handler capable of handling it.
func (h *CommandHandler) HandleCommandEx(ctx context.Context, cmd eventing.Command) ([]common.Event, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	if err := eventing.CheckCommand(cmd); err != nil {
		return nil, err
	}

	h.handlersMu.RLock()
	defer h.handlersMu.RUnlock()

	if handler, ok := h.handlers[cmd.CommandType()]; ok {
		return handler.HandleCommandEx(ctx, cmd)
	}

	return nil, ErrHandlerNotFound
}

// SetHandler adds/registers a handler for a specific command type. Only one handler can be set for
// a command type.
func (h *CommandHandler) SetHandler(handler eventing.CommandHandler, cmdType common.CommandType) error {
	h.handlersMu.Lock()
	defer h.handlersMu.Unlock()

	if _, ok := h.handlers[cmdType]; ok {
		return errors.Wrap(ErrHandlerAlreadySet, cmdType.String())
	}

	h.handlers[cmdType] = handler

	return nil
}
