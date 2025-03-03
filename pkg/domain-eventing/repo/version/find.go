package version

import (
	"context"
	"fmt"
	"time"

	eventing "github.com/sweetloveinyourheart/exploding-kittens/pkg/domain-eventing"
	log "github.com/sweetloveinyourheart/exploding-kittens/pkg/logger"

	"github.com/cockroachdb/errors"
	"go.uber.org/zap"
	"k8s.io/apimachinery/pkg/util/wait"
)

type Findable[T any] func(ctx context.Context, id string) (*T, error)

// RepoFind searches for an item in the repository with the given id.
// If the context contains a min version set by WithMinVersion it will only
// return an item if its version is at least min version. If a timeout or
// deadline is set on the context it will repeatedly try to get the item until
// either the version matches or the deadline is reached.
func RepoFind[T any](ctx context.Context, find Findable[T], id string, minVersion uint64) (*T, error) {
	// If there is no min version set just return the item as normally.
	if minVersion < 1 {
		okMin := false
		minVersion, okMin = MinVersionFromContext(ctx)
		if !okMin || minVersion < 1 {
			return find(ctx, id)
		}
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
		entity, err := findMinVersion(ctx, find, id, minVersion, nil)
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
func findMinVersion[T any](ctx context.Context, find Findable[T], id string, minVersion uint64, matchVersion func(e *T) bool) (*T, error) {
	var entity any
	var err error
	entity, err = find(ctx, id)
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
