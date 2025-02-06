package aggregate

import (
	"context"
	"time"

	"github.com/cockroachdb/errors"
	"go.uber.org/zap"
	"k8s.io/apimachinery/pkg/util/wait"

	eventing "github.com/sweetloveinyourheart/exploding-kittens/pkg/domain-eventing"
	"github.com/sweetloveinyourheart/exploding-kittens/pkg/domain-eventing/aggregate"
	"github.com/sweetloveinyourheart/exploding-kittens/pkg/domain-eventing/common"
	"github.com/sweetloveinyourheart/exploding-kittens/pkg/domain-eventing/middleware/oplock"
	log "github.com/sweetloveinyourheart/exploding-kittens/pkg/logger"
)

// ErrNilAggregateStore is when a dispatcher is created with a nil aggregate store.
var ErrNilAggregateStore = errors.New("aggregate store is nil")

// CommandHandler dispatches commands to an aggregate.
//
// The dispatch process is as follows:
// 1. The handler receives a command.
// 2. An aggregate is created or loaded using an aggregate store.
// 3. The aggregate's command handler is called.
// 4. The aggregate stores events in response to the command.
// 5. The new events are stored in the event store.
// 6. The events are published on the event bus after a successful store.
type CommandHandler struct {
	t                 common.AggregateType
	store             eventing.AggregateStore
	namespaceDetached bool
	deadline          time.Duration
}

// NewCommandHandler creates a new CommandHandler for an aggregate type.
func NewCommandHandler(t common.AggregateType, store eventing.AggregateStore, options ...Option) (*CommandHandler, error) {
	if store == nil {
		return nil, ErrNilAggregateStore
	}

	h := &CommandHandler{
		t:     t,
		store: store,
	}

	for _, option := range options {
		option(h)
	}

	return h, nil
}

var _ eventing.CommandHandler = (*CommandHandler)(nil)

type Option func(*CommandHandler)

func WithDeadline(d time.Duration) Option {
	return func(h *CommandHandler) {
		h.deadline = d
	}
}

func WithDetachedNamespace() Option {
	return func(h *CommandHandler) {
		h.namespaceDetached = true
	}
}

type WrappedCommand interface {
	Wrapped() eventing.Command
}

type KnownVersion interface {
	KnownVersion() uint64
}

// HandleCommand handles a command with the registered aggregate.
// Returns ErrAggregateNotFound if no aggregate could be found.
func (h *CommandHandler) HandleCommand(ctx context.Context, cmd eventing.Command) error {
	_, err := h.HandleCommandEx(ctx, cmd)
	return err
}

// HandleCommandEx handles a command with the registered aggregate.
// Returns ErrAggregateNotFound if no aggregate could be found.
func (h *CommandHandler) HandleCommandEx(ctx context.Context, cmd eventing.Command) ([]common.Event, error) {
	select {
	case <-ctx.Done():
		return nil, errors.Errorf("command handler initial context done: %w", ctx.Err())
	default:
	}

	waitCtx := ctx

	if err := eventing.CheckCommand(cmd); err != nil {
		return nil, err
	}

	delay := &wait.Backoff{
		Cap: 5 * time.Second,
	}
	// Skip the first duration, which is always 0.
	//_ = delay.Step()
	override, _ := ctx.Value(DeadlineOverrideKey).(bool)
	var cancelFunc context.CancelFunc
	if h.deadline > 0 && !override {
		waitCtx, cancelFunc = context.WithDeadline(ctx, time.Now().Add(h.deadline))
		defer cancelFunc()
	}
	deadline, hasDeadline := waitCtx.Deadline()

	// Get the time remaining for the deadline.
	var initialRemaining time.Duration
	if hasDeadline {
		initialRemaining = time.Until(deadline)
	}

	expectVersion := uint64(0)
	if knownVersion, ok := cmd.(KnownVersion); ok {
		expectVersion = knownVersion.KnownVersion()
	}

	if wrapped, ok := cmd.(WrappedCommand); ok {
		cmd = wrapped.Wrapped()
	}

	if h.namespaceDetached {
		ctx = eventing.NewContextWithDetachedNamespace(ctx)
	}

	var lastErr error

	for {
		a, err := h.store.Load(ctx, h.t, cmd.AggregateID())
		if err != nil {
			return nil, err
		} else if a == nil {
			return nil, errors.WithStack(eventing.ErrAggregateNotFound)
		}

		if va, ok := a.(aggregate.VersionedAggregate); ok {
			if expectVersion > 0 && va.AggregateVersion() < expectVersion {
				if va.AggregateVersion() == 0 {
					return nil, errors.WithStack(eventing.ErrAggregateNotFound)
				}
				// Wait for the next try or cancellation.
				select {
				case <-time.After(delay.Step()):
				case <-waitCtx.Done():
					return nil, errors.Errorf("command handler version deadline exceeded: %w %s %s %s %s", waitCtx.Err(), h.t, cmd.AggregateID(), cmd.AggregateType(), cmd.CommandType())
				}
				continue
			}
		}

		if err = a.HandleCommand(ctx, cmd); err != nil {
			if errors.Is(err, eventing.ErrIncorrectEntityVersion) || errors.Is(err, eventing.ErrEntityNotFound) {
				// If there is no deadline, return whatever we have at this point.
				if !hasDeadline {
					log.Global().DebugContext(ctx, "command handler: failed to handle command, no deadline", zap.Bool("has_deadline", hasDeadline), zap.Time("deadline", deadline), zap.Duration("initial", initialRemaining), zap.Error(err))
					return nil, &eventing.AggregateError{Err: err}
				}

				if errors.Is(err, eventing.ErrIncorrectEntityVersion) {
					ctx = oplock.ContextWithOplock(ctx)
				}

				lastErr = err

				// Wait for the next try or cancellation.
				select {
				case <-time.After(delay.Step()):
				case <-waitCtx.Done():
					log.Global().DebugContext(ctx, "command handler: failed to handle command, deadline exceeded", zap.NamedError("last_error", lastErr), zap.Bool("has_deadline", hasDeadline), zap.Time("deadline", deadline), zap.Duration("initial", initialRemaining), zap.Error(err))
					if lastErr == nil {
						lastErr = waitCtx.Err()
					}
					return nil, errors.Errorf("command handler deadline exceeded: %w %s %s %s %s %s", lastErr, waitCtx.Err(), h.t, cmd.AggregateID(), cmd.AggregateType(), cmd.CommandType())
				}
			} else if err != nil {
				// Return any real error.
				log.Global().DebugContext(ctx, "command handler: failed to handle command", zap.NamedError("last_error", lastErr), zap.Bool("has_deadline", hasDeadline), zap.Time("deadline", deadline), zap.Duration("initial", initialRemaining), zap.Error(err))
				return nil, &eventing.AggregateError{Err: err}
			}
		} else {
			var events []common.Event
			events, err = h.store.Save(ctx, a)
			var ese *eventing.EventStoreError
			if err != nil && errors.As(err, &ese) {
				if !errors.Is(ese.Err, eventing.ErrIncorrectEntityVersion) &&
					!errors.Is(ese.Err, eventing.ErrEntityNotFound) &&
					!errors.Is(err, eventing.ErrIncorrectEntityVersion) &&
					!errors.Is(err, eventing.ErrEntityNotFound) {
					log.Global().DebugContext(ctx, "command handler: failed to save aggregate", zap.Bool("has_deadline", hasDeadline), zap.Time("deadline", deadline), zap.Duration("initial", initialRemaining), zap.Error(err))
					return nil, err
				}
				events = nil
				lastErr = err
			} else if errors.Is(err, eventing.ErrIncorrectEntityVersion) || errors.Is(err, eventing.ErrEntityNotFound) {
				// Try again for incorrect version or if the entity was not found.
				events = nil
				lastErr = err
			} else if err != nil {
				// Return any real error.
				log.Global().DebugContext(ctx, "command handler: failed to save aggregate", zap.Bool("has_deadline", hasDeadline), zap.Time("deadline", deadline), zap.Duration("initial", initialRemaining), zap.Error(err))
				return nil, err
			} else {
				// Return
				if len(events) == 0 {
					events = nil
				}
				return events, nil
			}

			// If there is no deadline, return whatever we have at this point.
			if !hasDeadline {
				log.Global().DebugContext(ctx, "command handler: failed to save aggregate, no deadline", zap.NamedError("last_error", lastErr), zap.Bool("has_deadline", hasDeadline), zap.Time("deadline", deadline), zap.Duration("initial", initialRemaining), zap.Error(err))
				return nil, err
			}

			// Wait for the next try or cancellation.
			select {
			case <-time.After(delay.Step()):
			case <-waitCtx.Done():
				log.Global().DebugContext(ctx, "command handler: failed to save aggregate, deadline exceeded", zap.NamedError("last_error", lastErr), zap.Bool("has_deadline", hasDeadline), zap.Time("deadline", deadline), zap.Duration("initial", initialRemaining), zap.Error(err))
				if lastErr == nil {
					lastErr = waitCtx.Err()
				}
				return nil, errors.Errorf("command handler save deadline exceeded: %w %s %s %s %s %s", lastErr, waitCtx.Err(), h.t, cmd.AggregateID(), cmd.AggregateType(), cmd.CommandType())
			}
		}
	}
}

type ContextKey string

var DeadlineOverrideKey = ContextKey("deadline_override")
