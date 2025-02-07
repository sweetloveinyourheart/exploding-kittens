package projector

import (
	"context"
	"fmt"
	"time"

	"github.com/cockroachdb/errors"
	"github.com/gofrs/uuid"

	eventing "github.com/sweetloveinyourheart/exploding-kittens/pkg/domain-eventing"
	"github.com/sweetloveinyourheart/exploding-kittens/pkg/domain-eventing/common"
)

// Projector is a projector of events onto models.
type Projector[T any, PT eventing.GenericEntity[T]] interface {
	// ProjectorType returns the type of the projector.
	ProjectorType() common.ProjectorType

	// Project projects an event onto a model and returns the updated model or
	// an error.
	Project(context.Context, common.Event, *T) (*T, error)
}

var (
	// ErrModelNotSet is when a model factory is not set on the EventHandler.
	ErrModelNotSet = errors.New("model not set")
	// ErrModelRemoved is when a model has been removed.
	ErrModelRemoved = errors.New("model removed")
	// Returned if the model has not incremented its version as predicted.
	ErrIncorrectProjectedEntityVersion = errors.New("incorrect projected entity version")
)

// Error is an error in the projector.
type Error struct {
	// Err is the error that happened when projecting the event.
	Err error
	// Projector is the projector where the error happened.
	Projector string
	// Event is the event being projected.
	Event common.Event
	// EntityID of related operation.
	EntityID string
}

// Error implements the Error method of the errors.Error interface.
func (e *Error) Error() string {
	str := "projector '" + e.Projector + "': "

	if e.Err != nil {
		str += e.Err.Error()
	} else {
		str += "unknown error"
	}

	if e.EntityID != "" && e.EntityID != uuid.Nil.String() {
		str += fmt.Sprintf(", Entity(%s)", e.EntityID)
	}

	if e.Event != nil {
		str += ", " + e.Event.String()
	}

	return str
}

// Unwrap implements the errors.Unwrap method.
func (e *Error) Unwrap() error {
	return e.Err
}

// Cause implements the github.com/cockroachdb/errors Unwrap method.
func (e *Error) Cause() error {
	return e.Unwrap()
}

// EventHandler is a CQRS projection handler to run a Projector implementation.
type EventHandler[T any, PT eventing.GenericEntity[T]] struct {
	projector      Projector[T, PT]
	repo           eventing.ReadWriteRepo[T, PT]
	useWait        bool
	useRetryOnce   bool
	entityLookupFn func(common.Event) string
}

type _ent struct{}

func (_ *_ent) EntityID() string {
	return uuid.Nil.String()
}

var _ common.Entity = (*_ent)(nil)

var _ = eventing.EventHandler(&EventHandler[_ent, *_ent]{})

// NewEventHandler creates a new EventHandler.
func NewEventHandler[T any, PT eventing.GenericEntity[T]](projector Projector[T, PT], repo eventing.ReadWriteRepo[T, PT], options ...Option[T, PT]) *EventHandler[T, PT] {
	h := &EventHandler[T, PT]{
		projector:      projector,
		repo:           repo,
		entityLookupFn: defaultEntityLookupFn,
	}

	for _, option := range options {
		option(h)
	}

	return h
}

// Option is an option setter used to configure creation.
type Option[T any, PT eventing.GenericEntity[T]] func(*EventHandler[T, PT])

// WithWait adds waiting for the correct version when projecting.
func WithWait[T any, PT eventing.GenericEntity[T]]() Option[T, PT] {
	return func(h *EventHandler[T, PT]) {
		h.useWait = true
	}
}

// WithRetryOnce adds a single retry in case of version mismatch. Useful to
// let racy projections finish without an error.
func WithRetryOnce[T any, PT eventing.GenericEntity[T]]() Option[T, PT] {
	return func(h *EventHandler[T, PT]) {
		h.useRetryOnce = true
	}
}

// WithEntityLookup can be used to provide an alternative ID (from the aggregate ID)
// for fetching the projected entity. The lookup func can for example extract
// another field from the event or use a static ID for some singleton-like projections.
func WithEntityLookup[T any, PT eventing.GenericEntity[T]](f func(common.Event) string) Option[T, PT] {
	return func(h *EventHandler[T, PT]) {
		h.entityLookupFn = f
	}
}

// defaultEntitypLookupFn does a lookup by the aggregate ID of the event.
func defaultEntityLookupFn(event common.Event) string {
	return event.AggregateID()
}

// HandlerType implements the HandlerType method of the eventing.EventHandler interface.
func (h *EventHandler[T, PT]) HandlerType() common.EventHandlerType {
	return common.EventHandlerType("projector_" + h.projector.ProjectorType())
}

// HandleEvent implements the HandleEvent method of the eventing.EventHandler interface.
// It will try to find the correct version of the model, waiting for it the projector
// has the WithWait option set.
func (h *EventHandler[T, PT]) HandleEvent(ctx context.Context, event common.Event) error {
	if event == nil {
		return &Error{
			Err:       eventing.ErrMissingEvent,
			Projector: h.projector.ProjectorType().String(),
		}
	}

	// Used to retry once in case of a version mismatch.
	triedOnce := false
retryOnce:

	findCtx := ctx

	id := h.entityLookupFn(event)

	var entity any
	var err error
	// Get or create the model.
	entity, err = h.repo.Find(findCtx, id)
	if errors.Is(err, eventing.ErrEntityNotFound) {
		entity = new(T)
	} else if errors.Is(err, eventing.ErrIncorrectEntityVersion) {
		if h.useRetryOnce && !triedOnce {
			triedOnce = true

			time.Sleep(100 * time.Millisecond)

			goto retryOnce
		}

		return &Error{
			Err:       errors.Errorf("could not load entity with correct version: %w", err),
			Projector: h.projector.ProjectorType().String(),
			Event:     event,
			EntityID:  id,
		}
	} else if err != nil {
		return &Error{
			Err:       errors.Errorf("could not load entity: %w", err),
			Projector: h.projector.ProjectorType().String(),
			Event:     event,
			EntityID:  id,
		}
	}

	var newEntity any
	// Run the projection, which will possibly increment the version.
	newEntity, err = h.projector.Project(ctx, event, entity.(*T))
	if err != nil {
		return &Error{
			Err:       fmt.Errorf("could not project: %w", err),
			Projector: h.projector.ProjectorType().String(),
			Event:     event,
			EntityID:  id,
		}
	}

	// Update or remove the model.
	if newEntity.(*T) != nil {
		if newEntity.(PT).EntityID() != id {
			return &Error{
				Err:       fmt.Errorf("incorrect entity ID after projection"),
				Projector: h.projector.ProjectorType().String(),
				Event:     event,
				EntityID:  id,
			}
		}

		if err := h.repo.Save(ctx, newEntity.(*T)); err != nil {
			return &Error{
				Err:       fmt.Errorf("could not save: %w", err),
				Projector: h.projector.ProjectorType().String(),
				Event:     event,
				EntityID:  id,
			}
		}
	} else {
		if remover, ok := h.repo.(versionRemover); ok {
			if err := remover.RemoveVersion(ctx, id, event.Version()); err != nil {
				return &Error{
					Err:       fmt.Errorf("could not remove version: %w", err),
					Projector: h.projector.ProjectorType().String(),
					Event:     event,
					EntityID:  id,
				}
			}

			return nil
		}

		if err := h.repo.Remove(ctx, id); err != nil {
			return &Error{
				Err:       fmt.Errorf("could not remove: %w", err),
				Projector: h.projector.ProjectorType().String(),
				Event:     event,
				EntityID:  id,
			}
		}
	}

	return nil
}

type versionRemover interface {
	RemoveVersion(context.Context, string, uint64) error
}
