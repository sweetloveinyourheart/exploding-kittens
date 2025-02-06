package eventing

import (
	"context"
	"fmt"
	"strings"

	"github.com/cockroachdb/errors"
	"github.com/gofrs/uuid"

	"github.com/sweetloveinyourheart/exploding-kittens/pkg/domain-eventing/common"
	"github.com/sweetloveinyourheart/exploding-kittens/pkg/stringsutil"
)

// EventStore is an interface for an event sourcing event store.
type EventStore interface {
	// Save appends all events in the event stream to the store.
	Save(ctx context.Context, events []common.Event, originalVersion uint64) error

	// Load loads all events for the aggregate id from the store.
	Load(context.Context, string) ([]common.Event, error)

	// LoadFrom loads all events from version for the aggregate id from the store.
	LoadFrom(ctx context.Context, id string, version uint64) ([]common.Event, error)

	// Close closes the EventStore.
	Close() error
}

// SnapshotStore is an interface for snapshot store.
type SnapshotStore interface {
	LoadSnapshot(ctx context.Context, id string) (*Snapshot, error)
	SaveSnapshot(ctx context.Context, id string, snapshot Snapshot) error
}

var (
	// ErrMissingEvents missing events for save operation.
	ErrMissingEvents = errors.New("missing events")
	// ErrMismatchedEventAggregateIDs events in the same save operation is for different aggregate IDs.
	ErrMismatchedEventAggregateIDs = errors.New("mismatched event aggregate IDs")
	// ErrMismatchedEventAggregateTypes events in the same save operation is for different aggregate types.
	ErrMismatchedEventAggregateTypes = errors.New("mismatched event aggregate types")
	// ErrIncorrectEventVersion events in the same operation have non-serial versions or is not matching the original version.
	ErrIncorrectEventVersion = errors.New("incorrect event version")
	// ErrEventConflictFromOtherSave other events has been saved for this aggregate since the operation started.
	ErrEventConflictFromOtherSave = errors.New("event conflict from other save")
)

// EventStoreOperation is the operation done when an error happened.
type EventStoreOperation string

const (
	// Errors during loading of events.
	EventStoreOpLoad = "load"
	// Errors during saving of events.
	EventStoreOpSave = "save"
	// Errors during replacing of events.
	EventStoreOpReplace = "replace"
	// Errors during renaming of event types.
	EventStoreOpRename = "rename"
	// Errors during clearing of the event store.
	EventStoreOpClear = "clear"

	// Errors during loading of snapshot.
	EventStoreOpLoadSnapshot = "load_snapshot"
	// Errors during saving of snapshot.
	EventStoreOpSaveSnapshot = "save_snapshot"
)

// EventStoreError is an error in the event store.
type EventStoreError struct {
	// Err is the error.
	Err error
	// Op is the operation for the error.
	Op EventStoreOperation
	// AggregateType of related operation.
	AggregateType common.AggregateType
	// AggregateID of related operation.
	AggregateID string
	// AggregateVersion of related operation.
	AggregateVersion uint64
	// Events of the related operation.
	Events []common.Event
	// Subject of related operation.
	Subject string
}

// Error implements the Error method of the errors.Error interface.
func (e *EventStoreError) Error() string {
	str := "event store: "

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

		str += fmt.Sprintf(", %s(%s, v%d)", at, e.AggregateID, e.AggregateVersion)
	}

	if len(e.Events) > 0 {
		var es []string
		for _, ev := range e.Events {
			if ev != nil {
				es = append(es, ev.String())
			} else {
				es = append(es, "nil event")
			}
		}

		str += " [" + strings.Join(es, ", ") + "]"
	}

	if !stringsutil.IsBlank(e.Subject) {
		str += " subject(" + e.Subject + ")"
	}

	return str
}

// Unwrap implements the errors.Unwrap method.
func (e *EventStoreError) Unwrap() error {
	return e.Err
}

// Cause implements the github.com/cockroachdb/errors Unwrap method.
func (e *EventStoreError) Cause() error {
	return e.Unwrap()
}
