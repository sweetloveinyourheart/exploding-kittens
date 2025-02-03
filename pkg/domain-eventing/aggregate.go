package eventing

import (
	"fmt"
	"sync"

	"github.com/gofrs/uuid"
	"github.com/sweetloveinyourheart/exploding-kittens/pkg/domain-eventing/common"
)

// Aggregate is an interface representing a versioned data entity created from
// events. It receives commands and generates events that are stored.
//
// The aggregate is created/loaded and saved by the Repository inside the
// Dispatcher. A domain specific aggregate can either implement the full interface,
// or more commonly embed *AggregateBase to take care of the common methods.
type Aggregate interface {
	// Entity provides the ID of the aggregate.
	common.Entity

	// AggregateType returns the type name of the aggregate.
	// AggregateType() string
	AggregateType() common.AggregateType

	// CommandHandler is used to handle commands.
	SimpleCommandHandler
}

type GenericAggregate[T any] interface {
	Aggregate
	*T
}

type OnCreateHook interface {
	OnCreate(id string)
}

// RegisterAggregate registers an aggregate factory for a type. The factory is
// used to create concrete aggregate types when loading from the database.
//
// An example would be:
//
//	RegisterAggregate(func(id UUID) Aggregate { return &MyAggregate{id} })
func RegisterAggregate[T any, PT GenericAggregate[T]]() {
	// Check that the created aggregate matches the registered type.
	var aggregate any = new(T)

	a, ok := aggregate.(OnCreateHook)
	if ok {
		a.OnCreate(uuid.Nil.String())
	}

	aggregateType := aggregate.(PT).AggregateType()
	if aggregateType == common.AggregateType("") {
		panic("eventing: attempt to register empty aggregate type")
	}

	aggregatesMu.Lock()
	defer aggregatesMu.Unlock()

	if _, ok := aggregates[aggregateType]; ok {
		panic(fmt.Sprintf("eventing: registering duplicate types for %q", aggregateType))
	}

	aggregates[aggregateType] = func(id string) Aggregate {
		return PT(new(T))
	}
}

var aggregates = make(map[common.AggregateType]func(string) Aggregate)
var aggregatesMu sync.RWMutex
