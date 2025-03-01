package eventing

import "github.com/sweetloveinyourheart/exploding-kittens/pkg/domain-eventing/common"

type GenericEntity[T any] interface {
	common.Entity
	*T
}
