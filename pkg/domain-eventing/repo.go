package eventing

import (
	"context"

	"github.com/cockroachdb/errors"
	"github.com/gofrs/uuid"
)

var (
	// ErrEntityNotFound is when a entity could not be found.
	ErrEntityNotFound = errors.New("could not find entity")
	// ErrEntityHasNoVersion is when an entity has no version number.
	ErrEntityHasNoVersion = errors.New("entity has no version")
	// ErrIncorrectEntityVersion is when an entity has an incorrect version.
	ErrIncorrectEntityVersion = errors.New("incorrect entity version")
	// ErrModelIsMissing is when there has no model (model is nil, encountered while debugging).
	ErrModelIsMissing = errors.New("model is nil")
)

// RepoOperation is the operation done when an error happened.
type RepoOperation string

const (
	// RepoOpFind - errors during finding of an entity.
	RepoOpFind = "find"
	// RepoOpFindAll - errors during finding of all entities.
	RepoOpFindAll = "find all"
	// RepoOpFindQuery - errors during finding of entities by query.
	RepoOpFindQuery = "find query"
	// RepoOpSave - errors during saving of an entity.
	RepoOpSave = "save"
	// RepoOpRemove - errors during removing of an entity.
	RepoOpRemove = "remove"
	// RepoOpClear - errors during clearing of all entities.
	RepoOpClear = "clear"
	// RepoOpClose - errors during the closing of the repo
	RepoOpClose = "close"
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

// WriteRepo is a write repository for entities.
type WriteRepo[T any, PT GenericEntity[T]] interface {
	// Save saves a entity in the storage.
	Save(context.Context, *T) error

	// Remove removes a entity by ID from the storage.
	Remove(context.Context, string) error
}

// ReadWriteRepo is a combined read and write repo, mainly useful for testing.
type ReadWriteRepo[T any, PT GenericEntity[T]] interface {
	ReadRepo[T, PT]
	WriteRepo[T, PT]
}

// RepoError is an error in the read repository.
type RepoError struct {
	// Err is the error.
	Err error
	// Op is the operation for the error.
	Op RepoOperation
	// EntityID of related operation.
	EntityID string
	// EntityType is the type of entity.
	EntityType string
}

// Error implements the Error method of the errors.Error interface.
func (e *RepoError) Error() string {
	str := "repo: "

	if e.Op != "" {
		str += string(e.Op) + ": "
	}

	if e.Err != nil {
		str += e.Err.Error()
	} else {
		str += "unknown error"
	}

	if e.EntityID != "" && e.EntityID != uuid.Nil.String() {
		str += " " + e.EntityID + " "
	}

	if e.EntityType != "" {
		str += " (" + e.EntityType + ") "
	}

	return str
}

// Unwrap implements the errors.Unwrap method.
func (e *RepoError) Unwrap() error {
	return e.Err
}

// Cause implements the github.com/cockroachdb/errors Unwrap method.
func (e *RepoError) Cause() error {
	return e.Unwrap()
}
