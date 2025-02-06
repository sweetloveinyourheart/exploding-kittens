package eventing

import (
	"context"

	"github.com/cockroachdb/errors"
)

var (
	// ErrEntityNotFound is when a entity could not be found.
	ErrEntityNotFound = errors.New("could not find entity")
	// ErrEntityHasNoVersion is when an entity has no version number.
	ErrEntityHasNoVersion = errors.New("entity has no version")
	// ErrIncorrectEntityVersion is when an entity has an incorrect version.
	ErrIncorrectEntityVersion = errors.New("incorrect entity version")
	// ErrModelIsMissing is when a playerWagersEntry has no model (model is nil, encountered while debugging).
	ErrModelIsMissing = errors.New("playerWagersEntry model is nil")
)

// ReadRepo is a read repository for entities.
type ReadRepo[T any, PT GenericEntity[T]] interface {
	// InnerRepo returns the inner read repository, if there is one.
	// Useful for iterating a wrapped set of repositories to get a specific one.
	InnerRepo(context.Context) ReadRepo[T, PT]

	// Find returns an entity for an ID.
	Find(context.Context, string) (*T, error)

	// FindAll returns all entities in the repository.
	FindAll(context.Context) ([]*T, error)

	// Close closes the ReadRepo.
	Close() error
}
