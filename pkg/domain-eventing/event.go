package eventing

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/cockroachdb/errors"
	"github.com/gofrs/uuid"
	"github.com/jinzhu/copier"
	"github.com/sweetloveinyourheart/exploding-kittens/pkg/domain-eventing/common"
)

// EventOption is an option to use when creating events.
type EventOption func(common.Event)

// ForAggregate adds aggregate data when creating an event.
func ForAggregate(aggregateType common.AggregateType, aggregateID string, version uint64) EventOption {
	return func(e common.Event) {
		if evt, ok := e.(*event); ok {
			evt.aggregateType = aggregateType
			evt.aggregateID = aggregateID
			evt.version = version
		}
	}
}

// WithSubject sets the subject of the event.
func WithSubject(subject common.EventSubject) EventOption {
	return func(e common.Event) {
		if evt, ok := e.(*event); ok {
			evt.subject = subject.(*eventSubject)
		}
	}
}

// NewEvent creates a new event with a type and data, setting its timestamp.
func NewEvent(eventType common.EventType, data any, timestamp time.Time, options ...EventOption) common.Event {
	e := &event{
		eventType: eventType,
		data:      data,
		timestamp: timestamp,
		subject:   defaultSubject,
	}

	for _, option := range options {
		if option == nil {
			continue
		}

		option(e)
	}

	return e
}

// event is an internal representation of an event, returned when the aggregate
// uses NewEvent to create a new event. The events loaded from the db is
// represented by each DBs internal event type, implementing Event.
type event struct {
	eventType        common.EventType
	data             any
	timestamp        time.Time
	aggregateType    common.AggregateType
	namespace        string
	aggregateID      string
	version          uint64
	previousSequence uint64
	sequence         uint64
	sequenced        bool
	unregistered     bool
	metadata         map[string]any
	subject          *eventSubject
}

// EventType implements the EventType method of the Event interface.
func (e *event) EventType() common.EventType {
	return e.eventType
}

// Data implements the Data method of the Event interface.
func (e *event) Data() any {
	return e.data
}

// Timestamp implements the Timestamp method of the Event interface.
func (e *event) Timestamp() time.Time {
	return e.timestamp
}

// AggregateType implements the AggregateType method of the Event interface.
func (e *event) AggregateType() common.AggregateType {
	return e.aggregateType
}

func (e *event) Namespace() string {
	return e.namespace
}

// AggregateID implements the AggregateID method of the Event interface.
func (e *event) AggregateID() string {
	return e.aggregateID
}

// Version implements the Version method of the Event interface.
func (e *event) Version() uint64 {
	return e.version
}

// Metadata implements the Metadata method of the Event interface.
func (e *event) Metadata() map[string]any {
	return e.metadata
}

func (e *event) PreviousSequence() uint64 {
	return e.previousSequence
}

func (e *event) Sequence() uint64 {
	return e.sequence
}

func (e *event) Sequenced() bool {
	return e.sequenced
}

func (e *event) Unregistered() bool {
	return e.unregistered
}

// String implements the String method of the Event interface.
func (e *event) String() string {
	str := string(e.eventType)

	if e.aggregateID != "" && e.aggregateID != uuid.Nil.String() && e.version != 0 {
		str += fmt.Sprintf("(%s, v%d)", e.aggregateID, e.version)
	}

	return str
}

// Subject implements the Subject method of the Event interface.
func (e *event) Subject(ctx context.Context) common.EventSubject {
	return e.subject.forEvent(ctx, e)
}

// Clone implements the Clone method of the Event interface.
func (e *event) Clone(ctx context.Context) (common.Event, error) {
	clone := *e
	if b, ok := e.data.([]byte); ok {
		clone.data = make([]byte, len(b))
		copy(clone.data.([]byte), b)
	} else {
		if e.data != nil {
			newData, subj, err := CreateEventData(ctx, e.eventType)
			if err != nil {
				return nil, errors.WithStack(fmt.Errorf("could not create event data: %w", err))
			}

			err = copier.Copy(newData, e.Data())
			if err != nil {
				return nil, fmt.Errorf("could not copy event data: %w", err)
			}
			clone.data = newData
			clone.subject = subj.(*eventSubject)
		} else {
			if subjFact, ok := ctx.Value(eventSubjectKey).(*eventSubject); ok {
				clone.subject = subjFact
			} else if e.subject != nil {
				clone.subject = e.subject
			} else {
				clone.subject = defaultSubject
			}
		}

	}
	return &clone, nil
}

// RegisterEventData registers an event data factory for a type. The factory is
// used to create concrete event data structs when loading from the database.
//
// An example would be:
//
//	RegisterEventData(MyEventType, func() Event { return &MyEventData{} })
func RegisterEventData[T any](eventType common.EventType, registrationOpts ...EventRegistrationOption) {
	if eventType == common.EventType("") {
		panic("eventing: attempt to register empty event type")
	}

	eventDataFactoriesMu.Lock()
	defer eventDataFactoriesMu.Unlock()

	if _, ok := eventDataFactories[eventType]; ok {
		panic(fmt.Sprintf("eventing: registering duplicate types for %q", eventType))
	}

	var opts eventRegistrationOptions
	for _, opt := range registrationOpts {
		opt(&opts)
	}

	eventDataFactories[eventType] = func(ctx context.Context) (any, *eventSubject) {
		var result any = new(T)
		sub := new(eventSubject)
		if opts.subjectFunc != nil {
			sub.subjectFunc = opts.subjectFunc
		} else {
			sub.subjectFunc = defaultSubjectFunc
		}
		if opts.subjectRootFunc != nil {
			sub.subjectRootFunc = opts.subjectRootFunc
		} else {
			sub.subjectRootFunc = defaultSubjectRootFunc
		}
		if opts.tokensFunc != nil {
			sub.subjectTokenPos = opts.subjectTokenPos
		} else {
			sub.subjectTokenPos = 1
		}
		if opts.tokensFunc != nil {
			sub.tokensFunc = opts.tokensFunc
		} else {
			sub.tokensFunc = defaultTokensFunc
		}

		return result, sub
	}
}

type EventRegistrationOption func(opts *eventRegistrationOptions)

type eventRegistrationOptions struct {
	subjectFunc     func(ctx context.Context, event common.Event) string
	subjectRootFunc func(ctx context.Context, event common.Event) string
	subjectTokenPos int
	tokensFunc      func(ctx context.Context, event common.Event) []common.EventSubjectToken
}

func WithRegisterSubjectFunc(subjectFunc func(ctx context.Context, event common.Event) string) EventRegistrationOption {
	return func(opts *eventRegistrationOptions) {
		opts.subjectFunc = subjectFunc
	}
}

func WithRegisterSubjectRootFunc(subjectRootFunc func(ctx context.Context, event common.Event) string) EventRegistrationOption {
	return func(opts *eventRegistrationOptions) {
		opts.subjectRootFunc = subjectRootFunc
	}
}

func WithRegisterSubjectTokenPosition(subjectTokenPos int) EventRegistrationOption {
	return func(opts *eventRegistrationOptions) {
		opts.subjectTokenPos = subjectTokenPos
	}
}

func WithRegisterTokensFunc(tokensFunc func(ctx context.Context, event common.Event) []common.EventSubjectToken) EventRegistrationOption {
	return func(opts *eventRegistrationOptions) {
		opts.tokensFunc = tokensFunc
	}
}

// ErrEventDataNotRegistered is when no event data factory was registered.
var ErrEventDataNotRegistered = errors.New("event data not registered")

// CreateEventData creates an event data of a type using the factory registered
// with RegisterEventData.
func CreateEventData(ctx context.Context, eventType common.EventType) (any, common.EventSubject, error) {
	eventDataFactoriesMu.RLock()
	defer eventDataFactoriesMu.RUnlock()

	if factory, ok := eventDataFactories[eventType]; ok {
		result, sub := factory(ctx)
		return result, sub, nil
	}

	return nil, nil, errors.WithStack(errors.Wrap(ErrEventDataNotRegistered, "event type "+string(eventType)))
}

var eventDataFactories = make(map[common.EventType]func(ctx context.Context) (any, *eventSubject))
var eventDataFactoriesMu sync.RWMutex
