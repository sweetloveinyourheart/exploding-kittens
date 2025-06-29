package tracing

import (
	"context"
	"fmt"

	"github.com/cockroachdb/errors"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"

	eventing "github.com/sweetloveinyourheart/exploding-kittens/pkg/domain-eventing"
)

// Repo is a ReadWriteRepo that adds tracing.
type Repo[T any, PT eventing.GenericEntity[T]] struct {
	eventing.ReadWriteRepo[T, PT]
	tracer     trace.Tracer
	entityType string
}

// NewRepo creates a new Repo.
func NewRepo[T any, PT eventing.GenericEntity[T]](repo eventing.ReadWriteRepo[T, PT]) *Repo[T, PT] {
	return &Repo[T, PT]{
		ReadWriteRepo: repo,
		tracer:        otel.Tracer(TracerName),
		entityType:    fmt.Sprintf("%T", new(T)),
	}
}

// InnerRepo implements the InnerRepo method of the eventing.ReadRepo interface.
func (r *Repo[T, PT]) InnerRepo(ctx context.Context) eventing.ReadRepo[T, PT] {
	return r.ReadWriteRepo
}

// AdaptFrom tries to convert an eventing.ReadRepo into a Repo by recursively looking at
// inner repos. Returns nil if none was found.
func AdaptFrom[T any, PT eventing.GenericEntity[T]](ctx context.Context, repo eventing.ReadRepo[T, PT]) *Repo[T, PT] {
	if repo == nil {
		return nil
	}

	if r, ok := repo.(*Repo[T, PT]); ok {
		return r
	}

	return AdaptFrom(ctx, repo.InnerRepo(ctx))
}

// Find implements the Find method of the eventing.ReadModel interface.
func (r *Repo[T, PT]) Find(ctx context.Context, id string) (*T, error) {
	opName := "Repo.Find"

	opts := []trace.SpanStartOption{
		trace.WithAttributes(
			AggregateID(id),
			EntityType(r.entityType)),
		trace.WithSpanKind(trace.SpanKindInternal),
	}

	ctx, span := r.tracer.Start(ctx, opName, opts...)

	entity, err := r.ReadWriteRepo.Find(ctx, id)

	if err != nil && !errors.Is(err, eventing.ErrEntityNotFound) {
		span.RecordError(err)
	}

	span.End()

	return entity, err
}

// FindAll implements the FindAll method of the eventing.ReadRepo interface.
func (r *Repo[T, PT]) FindAll(ctx context.Context) ([]*T, error) {
	opName := "Repo.FindAll"

	opts := []trace.SpanStartOption{
		trace.WithAttributes(EntityType(r.entityType)),
		trace.WithSpanKind(trace.SpanKindInternal),
	}

	ctx, span := r.tracer.Start(ctx, opName, opts...)

	entities, err := r.ReadWriteRepo.FindAll(ctx)
	if err != nil {
		span.RecordError(err)
	}

	span.End()

	return entities, err
}

// Save implements the Save method of the eventing.WriteRepo interface.
func (r *Repo[T, PT]) Save(ctx context.Context, entity *T) error {
	opName := "Repo.Save"

	opts := []trace.SpanStartOption{
		trace.WithAttributes(
			AggregateID(PT(entity).EntityID()),
			EntityType(r.entityType),
		),
		trace.WithSpanKind(trace.SpanKindInternal),
	}

	ctx, span := r.tracer.Start(ctx, opName, opts...)

	err := r.ReadWriteRepo.Save(ctx, entity)
	if err != nil {
		span.RecordError(err)
	}

	span.End()

	return err
}

// RemoveVersion removes an entity by ID and version.
func (r *Repo[T, PT]) RemoveVersion(ctx context.Context, id string, version uint64) error {
	opName := "Repo.Remove"

	opts := []trace.SpanStartOption{
		trace.WithAttributes(AggregateID(id)),
		trace.WithSpanKind(trace.SpanKindInternal),
	}

	ctx, span := r.tracer.Start(ctx, opName, opts...)

	if remover, ok := r.ReadWriteRepo.(versionRemover); ok {
		err := remover.RemoveVersion(ctx, id, version)
		if err != nil {
			span.RecordError(err)
		}

		span.End()

		return err
	}

	err := r.ReadWriteRepo.Remove(ctx, id)
	if err != nil {
		span.RecordError(err)
	}
	span.End()

	return err
}

// Remove implements the Remove method of the eventing.WriteRepo interface.
func (r *Repo[T, PT]) Remove(ctx context.Context, id string) error {
	opName := "Repo.Remove"

	opts := []trace.SpanStartOption{
		trace.WithAttributes(AggregateID(id)),
		trace.WithSpanKind(trace.SpanKindInternal),
	}

	ctx, span := r.tracer.Start(ctx, opName, opts...)

	err := r.ReadWriteRepo.Remove(ctx, id)
	if err != nil {
		span.RecordError(err)
	}

	span.End()

	return err
}

type versionRemover interface {
	RemoveVersion(context.Context, string, uint64) error
}
