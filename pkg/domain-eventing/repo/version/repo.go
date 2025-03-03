package version

import (
	"context"
	"fmt"
	"time"

	"github.com/cockroachdb/errors"
	"go.uber.org/zap"
	"k8s.io/apimachinery/pkg/util/wait"

	eventing "github.com/sweetloveinyourheart/exploding-kittens/pkg/domain-eventing"
	log "github.com/sweetloveinyourheart/exploding-kittens/pkg/logger"
)

// Repo is a middleware that adds version checking to a read repository.
type Repo[T any, PT eventing.GenericEntity[T]] struct {
	eventing.ReadWriteRepo[T, PT]
}

// NewRepo creates a new Repo.
func NewRepo[T any, PT eventing.GenericEntity[T]](repo eventing.ReadWriteRepo[T, PT]) *Repo[T, PT] {
	return &Repo[T, PT]{
		ReadWriteRepo: repo,
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
// If the context contains a min version set by WithMinVersion it will only
// return an item if its version is at least min version. If a timeout or
// deadline is set on the context it will repeatedly try to get the item until
// either the version matches or the deadline is reached.
func (r *Repo[T, PT]) Find(ctx context.Context, id string) (*T, error) {
	// If there is no min version set just return the item as normally.
	minVersion, okMin := MinVersionFromContext(ctx)
	matchVersion, okMatch := MatchVersionFromContext[T, PT](ctx)
	if (!okMin || minVersion < 1) && (!okMatch || matchVersion == nil) {
		return r.ReadWriteRepo.Find(ctx, id)
	}

	// Try to get the correct version, retry with exponentially longer intervals
	// until the deadline expires. If there is no deadline just try once.
	delay := &wait.Backoff{
		Cap: 5 * time.Second,
	}
	// Skip the first duration, which is always 0.
	//_ = delay.Step()
	_, hasDeadline := ctx.Deadline()

	for {
		entity, err := r.findMinVersion(ctx, id, minVersion, matchVersion)
		if errors.Is(err, eventing.ErrIncorrectEntityVersion) || errors.Is(err, eventing.ErrEntityNotFound) {
			// Try again for incorrect version or if the entity was not found.
		} else if err != nil {
			// Return any real error.
			return nil, err
		} else {
			// Return the entity.
			return entity, nil
		}

		// If there is no deadline, return whatever we have at this point.
		if !hasDeadline {
			return entity, err
		}

		// Wait for the next try or cancellation.
		select {
		case <-time.After(delay.Step()):
		case <-ctx.Done():
			log.Global().WarnContext(ctx, "findMinVersion timed out", zap.Error(ctx.Err()), zap.Any("inner_error", err))
			if err == nil {
				err = ctx.Err()
			}
			return nil, errors.WithStack(errors.Wrap(err, fmt.Sprintf("could not find versioned entity with id %s of type %T", id, entity)))
		}
	}
}

// findMinVersion finds an item if it has a version and it is at least minVersion.
func (r *Repo[T, PT]) findMinVersion(ctx context.Context, id string, minVersion uint64, matchVersion func(e *T) bool) (*T, error) {
	var entity any
	var err error
	entity, err = r.ReadWriteRepo.Find(ctx, id)
	if err != nil {
		return nil, err
	}

	versionable, ok := entity.(eventing.Versionable)
	if !ok {
		return nil, errors.WithStack(&eventing.RepoError{
			Err:        eventing.ErrEntityHasNoVersion,
			EntityID:   id,
			Op:         eventing.RepoOpFind,
			EntityType: fmt.Sprintf("%T", *new(T)),
		})
	}

	if matchVersion != nil && !matchVersion(entity.(*T)) {
		return nil, errors.WithStack(&eventing.RepoError{
			Err:        errors.Wrap(eventing.ErrIncorrectEntityVersion, fmt.Sprintf("entity does not match version: %d", versionable.AggregateVersion())),
			EntityID:   id,
			Op:         eventing.RepoOpFind,
			EntityType: fmt.Sprintf("%T", *new(T)),
		})
	}

	if versionable.AggregateVersion() < minVersion {
		return nil, errors.WithStack(&eventing.RepoError{
			Err:        errors.Wrap(eventing.ErrIncorrectEntityVersion, fmt.Sprintf("entity version is too low: %d want %d", versionable.AggregateVersion(), minVersion)),
			EntityID:   id,
			Op:         eventing.RepoOpFind,
			EntityType: fmt.Sprintf("%T", *new(T)),
		})
	}

	return entity.(*T), nil
}
