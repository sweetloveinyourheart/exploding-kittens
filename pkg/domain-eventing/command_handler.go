package eventing

import (
	"context"

	"github.com/sweetloveinyourheart/exploding-kittens/pkg/domain-eventing/common"
)

// CommandHandler is an interface that all handlers of commands should implement.
type CommandHandler interface {
	HandleCommand(context.Context, Command) error
	HandleCommandEx(context.Context, Command) ([]common.Event, error)
}
