package eventing

import "github.com/sweetloveinyourheart/exploding-kittens/pkg/domain-eventing/common"

type GenericEntity[T any] interface {
	common.Entity
	*T
}

// Versionable is an item that has a version number,
// used by version.ReadRepo and projector.EventHandler.
type Versionable interface {
	// AggregateVersion returns the version of the item.
	AggregateVersion() uint64
}
