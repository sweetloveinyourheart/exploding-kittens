package common

import (
	"context"
	"time"
)

// Event is a domain event describing a change that has happened to an aggregate.
//
// An event struct and type name should:
//  1. Be in past tense (CustomerMoved)
//  2. Contain the intent (CustomerMoved vs CustomerAddressCorrected).
//
// The event should contain all the data needed when applying/handling it.
type Event = Of[any]

type PT[T any] interface {
	*T
}

// Of is a domain event describing a change that has happened to an aggregate.
//
// An event struct and type name should:
//  1. Be in past tense (CustomerMoved)
//  2. Contain the intent (CustomerMoved vs CustomerAddressCorrected).
//
// The event should contain all the data needed when applying/handling it.
type Of[Data any] interface {
	// EventType returns the type of the event.
	EventType() EventType
	// Data is the data attached to the event.
	Data() Data
	// Timestamp of when the event was created.
	Timestamp() time.Time

	// AggregateType is the type of the aggregate that the event can be
	// applied to.
	AggregateType() AggregateType
	// AggregateID is the ID of the aggregate that the event belongs to.
	AggregateID() string
	// Version is the version of the aggregate after the event has been applied.
	Version() uint64

	// Metadata is app-specific metadata such as request ID, originating user etc.
	Metadata() map[string]any

	// String returns a representation of the event.
	String() string

	// Subject returns the subject of the event.
	Subject(context.Context) EventSubject

	// Clone returns a copy of the event.
	Clone(ctx context.Context) (Event, error)
}

type EventSubject interface {
	Subject() string
	SubjectRoot() string
	SubjectTokenPosition() int
	Tokens() []EventSubjectToken
}

type EventSubjectToken interface {
	Key() string
	Value() any
	Position() int
	Description() string
	Type() string
}

// EventType is the type of an event, used as its unique identifier.
type EventType string

// String returns the string representation of an event type.
func (et EventType) String() string {
	return string(et)
}

// EventHandlerType is the type of event handler, used as its unique identifier.
type EventHandlerType string

// String returns the string representation of an event handler type.
func (ht EventHandlerType) String() string {
	return string(ht)
}

type UnregisteredEvent interface {
	Unregistered() bool
}
