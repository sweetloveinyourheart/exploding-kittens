package eventing

import (
	"context"

	"github.com/sweetloveinyourheart/exploding-kittens/pkg/domain-eventing/common"
)

// SimpleCommandHandler is an interface that all handlers of commands should implement.
type SimpleCommandHandler interface {
	HandleCommand(context.Context, Command) error
}

// CommandHandler is an interface that all handlers of commands should implement.
type CommandHandler interface {
	HandleCommand(context.Context, Command) error
	HandleCommandEx(context.Context, Command) ([]common.Event, error)
}

// CommandHandlerFunc is a function that can be used as a command handler.
type CommandHandlerFunc func(context.Context, Command) ([]common.Event, error)

// HandleCommand implements the HandleCommand method of the CommandHandler.
func (h CommandHandlerFunc) HandleCommand(ctx context.Context, cmd Command) error {
	_, err := h(ctx, cmd)
	return err
}

// HandleCommandEx implements the HandleCommand method of the CommandHandler.
func (h CommandHandlerFunc) HandleCommandEx(ctx context.Context, cmd Command) ([]common.Event, error) {
	return h(ctx, cmd)
}
