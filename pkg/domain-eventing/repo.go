package eventing

import "context"

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
