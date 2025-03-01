package aggregate

import (
	"time"

	eventing "github.com/sweetloveinyourheart/exploding-kittens/pkg/domain-eventing"
	"github.com/sweetloveinyourheart/exploding-kittens/pkg/domain-eventing/common"
)

type AggregateBase struct {
	id          string
	subjFactory common.EventSubject
	t           common.AggregateType
	v           uint64
	s           uint64
	events      []common.Event
}

// NewAggregateBase creates an aggregate.
func NewAggregateBase(t common.AggregateType, subjFactory common.EventSubject, id string) *AggregateBase {
	return &AggregateBase{
		id:          id,
		t:           t,
		subjFactory: subjFactory,
	}
}

// EntityID implements the EntityID method of the eventing.Entity and eventing.Aggregate interface.
func (a *AggregateBase) EntityID() string {
	return a.id
}

// AggregateType implements the AggregateType method of the eventing.Aggregate interface.
func (a *AggregateBase) AggregateType() common.AggregateType {
	return a.t
}

// AggregateVersion implements the AggregateVersion method of the Aggregate interface.
func (a *AggregateBase) AggregateVersion() uint64 {
	return a.v
}

// SetAggregateVersion implements the SetAggregateVersion method of the Aggregate interface.
func (a *AggregateBase) SetAggregateVersion(v uint64) {
	a.v = v
}

// AggregateSequence implements the AggregateSequence method of the Aggregate interface.
func (a *AggregateBase) AggregateSequence() uint64 {
	return a.s
}

// SetAggregateSequence implements the SetAggregateSequence method of the Aggregate interface.
func (a *AggregateBase) SetAggregateSequence(v uint64) {
	a.s = v
}

// UncommittedEvents implements the UncommittedEvents method of the eventing.EventSource
// interface.
func (a *AggregateBase) UncommittedEvents() []common.Event {
	return a.events
}

// ClearUncommittedEvents implements the ClearUncommittedEvents method of the eventing.EventSource
// interface.
func (a *AggregateBase) ClearUncommittedEvents() {
	a.events = nil
}

func (a *AggregateBase) SubjectFactory() common.EventSubject {
	return a.subjFactory
}

// AppendEvent appends an event for later retrieval by Events().
func (a *AggregateBase) AppendEvent(t common.EventType, data any, timestamp time.Time, options ...eventing.EventOption) common.Event {
	options = append(options, eventing.ForAggregate(
		a.AggregateType(),
		a.EntityID(),
		a.AggregateVersion()+uint64(len(a.events))+1),
		eventing.WithSubject(a.SubjectFactory()),
	)
	e := eventing.NewEvent(t, data, timestamp, options...)
	a.events = append(a.events, e)

	return e
}
