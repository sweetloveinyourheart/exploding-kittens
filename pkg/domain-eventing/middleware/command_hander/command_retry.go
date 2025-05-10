package retry

import (
	"context"
	"time"

	"github.com/avast/retry-go"

	eventing "github.com/sweetloveinyourheart/exploding-kittens/pkg/domain-eventing"
	"github.com/sweetloveinyourheart/exploding-kittens/pkg/domain-eventing/common"
	"github.com/sweetloveinyourheart/exploding-kittens/pkg/domain-eventing/middleware/oplock"
)

// NewCommandHandlerMiddleware returns a new command handler middleware that adds tracing spans.
func NewCommandHandlerMiddleware(options ...retry.Option) eventing.CommandHandlerMiddleware {
	options = append(options, retry.Attempts(3), retry.MaxDelay(5*time.Second), retry.LastErrorOnly(true))
	return eventing.CommandHandlerMiddleware(func(h eventing.CommandHandler) eventing.CommandHandler {
		return eventing.CommandHandlerFunc(func(ctx context.Context, cmd eventing.Command) ([]common.Event, error) {
			var result []common.Event
			err := retry.Do(func() error {
				var err error
				result, err = h.HandleCommandEx(ctx, cmd)
				if err != nil {
					ctx = oplock.ContextWithOplock(ctx)
				}
				return err
			}, options...)
			if err != nil {
				return nil, err
			}
			return result, nil
		})
	})
}
