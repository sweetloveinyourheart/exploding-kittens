package eventing

import (
	"context"
	"fmt"
	"reflect"

	"github.com/sweetloveinyourheart/exploding-kittens/pkg/domain-eventing/common"
)

var defaultSubject = &eventSubject{
	ctx:             context.Background(),
	subjectFunc:     defaultSubjectFunc,
	subjectRootFunc: defaultSubjectRootFunc,
	subjectTokenPos: 1,
	tokensFunc:      defaultTokensFunc,
}

type SubjectFunc func(ctx context.Context, event common.Event) string
type TokensFunc func(ctx context.Context, event common.Event) []common.EventSubjectToken

type eventSubject struct {
	ctx             context.Context
	event           common.Event
	subjectFunc     SubjectFunc
	subjectRootFunc SubjectFunc
	subjectTokenPos int
	tokensFunc      TokensFunc
}

var _ common.EventSubject = (*eventSubject)(nil)

func NewEventSubjectFactory(subjFunc, subjRootFunc SubjectFunc, subjTokenPos int, tokensFunc TokensFunc) common.EventSubject {
	return &eventSubject{
		subjectFunc:     subjFunc,
		subjectRootFunc: subjRootFunc,
		subjectTokenPos: subjTokenPos,
		tokensFunc:      tokensFunc,
	}
}

func (e *eventSubject) Subject() string {
	return e.subjectFunc(e.ctx, e.event)
}

func (e *eventSubject) SubjectRoot() string {
	return e.subjectRootFunc(e.ctx, e.event)
}

func (e *eventSubject) SubjectTokenPosition() int {
	return e.subjectTokenPos
}

func (e *eventSubject) Tokens() []common.EventSubjectToken {
	return e.tokensFunc(e.ctx, e.event)
}

func (e *eventSubject) forEvent(ctx context.Context, event common.Event) *eventSubject {
	return &eventSubject{
		ctx:             ctx,
		event:           event,
		subjectFunc:     e.subjectFunc,
		subjectRootFunc: e.subjectRootFunc,
		subjectTokenPos: e.subjectTokenPos,
		tokensFunc:      e.tokensFunc,
	}
}

type eventSubjectToken struct {
	key         string
	description string
	value       any
	position    int
}

var _ common.EventSubjectToken = (*eventSubjectToken)(nil)

func NewEventSubjectToken(key, description string, value any, position int) *eventSubjectToken {
	return &eventSubjectToken{
		key:         key,
		description: description,
		value:       value,
		position:    position,
	}
}

func (e *eventSubjectToken) Key() string {
	return e.key
}

func (e *eventSubjectToken) Value() any {
	return e.value
}

func (e *eventSubjectToken) Position() int {
	return e.position
}

func (e *eventSubjectToken) Description() string {
	return e.description
}

func (e *eventSubjectToken) Type() string {
	return reflect.TypeOf(e.value).String()
}

func defaultSubjectFunc(ctx context.Context, event common.Event) string {
	return fmt.Sprintf("%s.%s", event.AggregateType(), event.AggregateID())
}

func defaultSubjectRootFunc(ctx context.Context, event common.Event) string {
	return event.AggregateType().String()
}

func defaultTokensFunc(ctx context.Context, event common.Event) []common.EventSubjectToken {
	return []common.EventSubjectToken{
		NewEventSubjectToken("aggregate_type", "Aggregate Type", event.AggregateType(), 0),
		NewEventSubjectToken("aggregate_id", "Aggregate ID", event.AggregateID(), 1),
	}
}
