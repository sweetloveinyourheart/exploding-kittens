package eventing

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/cockroachdb/errors"
	"github.com/gofrs/uuid"

	"github.com/sweetloveinyourheart/exploding-kittens/pkg/domain-eventing/common"
)

var (
	// ErrAggregateNotFound is when no aggregate can be found.
	ErrAggregateNotFound = errors.New("aggregate not found")
	// ErrAggregateNotRegistered is when no aggregate factory was registered.
	ErrAggregateNotRegistered = errors.New("aggregate not registered")
	// ErrMissingAggregateType is when the aggregate type is missing.
	ErrMissingAggregateType = errors.New("missing aggregate type")
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

// AggregateStore is responsible for loading and saving aggregates.
type AggregateStore interface {
	// Load loads the most recent version of an aggregate with a type and id.
	Load(context.Context, common.AggregateType, string) (Aggregate, error)

	// Save saves the uncommitted events for an aggregate.
	Save(context.Context, Aggregate) ([]common.Event, error)
}

// SnapshotStrategy determines if a snapshot should be taken or not.
type SnapshotStrategy interface {
	ShouldTakeSnapshot(lastSnapshotVersion uint64, lastSnapshotTimestamp time.Time, event common.Event) bool
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

// AggregateStoreOperation is the operation done when an error happened.
type AggregateStoreOperation string

const (
	// AggregateStoreOpLoad - errors during loading of aggregates.
	AggregateStoreOpLoad = "load"
	// AggregateStoreOpSave - errors during saving of aggregates.
	AggregateStoreOpSave = "save"
)

// AggregateStoreError contains related info about errors in the store.
type AggregateStoreError struct {
	// Err is the error that happened when applying the event.
	Err error
	// Op is the operation for the error.
	Op AggregateStoreOperation
	// AggregateType of related operation.
	AggregateType common.AggregateType
	// AggregateID of related operation.
	AggregateID string
}

// Error implements the Error method of the error interface.
func (e *AggregateStoreError) Error() string {
	str := "aggregate store: "

	if e.Op != "" {
		str += string(e.Op) + ": "
	}

	if e.Err != nil {
		str += e.Err.Error()
	} else {
		str += "unknown error"
	}

	if e.AggregateID != "" && e.AggregateID != uuid.Nil.String() {
		at := "Aggregate"
		if e.AggregateType != "" {
			at = string(e.AggregateType)
		}

		str += fmt.Sprintf(", %s(%s)", at, e.AggregateID)
	}

	return str
}

// Unwrap implements the errors.Unwrap method.
func (e *AggregateStoreError) Unwrap() error {
	return e.Err
}

// Cause implements the github.com/cockroachdb/errors Unwrap method.
func (e *AggregateStoreError) Cause() error {
	return e.Unwrap()
}

// AggregateError is an error caused in the aggregate when handling a command.
type AggregateError struct {
	// Err is the error.
	Err error
}

// Error implements the Error method of the errors.Error interface.
func (e *AggregateError) Error() string {
	return "aggregate error: " + e.Err.Error()
}

// Unwrap implements the errors.Unwrap method.
func (e *AggregateError) Unwrap() error {
	return e.Err
}

// Cause implements the github.com/cockroachdb/errors Unwrap method.
func (e *AggregateError) Cause() error {
	return e.Unwrap()
}

// CreateAggregate creates an aggregate of a type with an ID using the factory
// registered with RegisterAggregate.
func CreateAggregate(aggregateType common.AggregateType, id string) (Aggregate, error) {
	aggregatesMu.RLock()
	defer aggregatesMu.RUnlock()

	if factory, ok := aggregates[aggregateType]; ok {
		agg := factory(id)
		a, ok := agg.(OnCreateHook)
		if ok {
			a.OnCreate(id)
		}
		return agg, nil
	}

	return nil, ErrAggregateNotRegistered
}

var aggregates = make(map[common.AggregateType]func(string) Aggregate)
var aggregatesMu sync.RWMutex
