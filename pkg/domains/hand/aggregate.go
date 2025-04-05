package hand

import (
	"context"
	"fmt"

	"github.com/gofrs/uuid"

	eventing "github.com/sweetloveinyourheart/exploding-kittens/pkg/domain-eventing"
	"github.com/sweetloveinyourheart/exploding-kittens/pkg/domain-eventing/aggregate"
	"github.com/sweetloveinyourheart/exploding-kittens/pkg/domain-eventing/common"
	"github.com/sweetloveinyourheart/exploding-kittens/pkg/timeutil"
)

var SubjectFactory common.EventSubject

func init() {
	eventing.RegisterAggregate[Aggregate, *Aggregate]()
	type handIDer interface {
		GetHandID() uuid.UUID
	}

	subjectFunc := func(ctx context.Context, event common.Event) string {
		if event == nil {
			return ""
		}

		if a, ok := event.Data().(handIDer); ok {
			return fmt.Sprintf("%s.%s.%s", event.AggregateType(), a.GetHandID(), event.AggregateID())
		}

		return ""
	}

	subjectRootFunc := func(ctx context.Context, event common.Event) string {
		if event == nil {
			return ""
		}
		return string(event.AggregateType())
	}

	subjectTokenFunc := func(ctx context.Context, event common.Event) []common.EventSubjectToken {
		if event == nil {
			return []common.EventSubjectToken{
				eventing.NewEventSubjectToken("aggregate_type", "Aggregate Type", AggregateType, 0),
				eventing.NewEventSubjectToken("hand_id", "Hand ID", uuid.Nil, 1),
				eventing.NewEventSubjectToken("aggregate_id", "Aggregate ID", uuid.Nil, 2),
			}
		}
		if a, ok := event.Data().(handIDer); ok {
			return []common.EventSubjectToken{
				eventing.NewEventSubjectToken("aggregate_type", "Aggregate Type", AggregateType, 0),
				eventing.NewEventSubjectToken("hand_id", "Hand ID", a.GetHandID(), 1),
				eventing.NewEventSubjectToken("aggregate_id", "Aggregate ID", event.AggregateID(), 2),
			}
		}
		return nil
	}
	SubjectFactory = eventing.NewEventSubjectFactory(
		subjectFunc,
		subjectRootFunc,
		2,
		subjectTokenFunc)
	registerEvents(subjectFunc, subjectRootFunc, 2, subjectTokenFunc)
}

const AggregateType = common.AggregateType("hand")

type Aggregate struct {
	*aggregate.AggregateBase
	currentHandID uuid.UUID
}

var _ eventing.Aggregate = (*Aggregate)(nil)
var _ eventing.OnCreateHook = (*Aggregate)(nil)

// TimeNow is a mockable version of time.Now.
var TimeNow = timeutil.NowRoundedForGranularity

func (a *Aggregate) OnCreate(id string) {
	a.AggregateBase = aggregate.NewAggregateBase(AggregateType, SubjectFactory, id)
}

func (a *Aggregate) validateCommand(cmd eventing.Command) error {
	switch typed := cmd.(type) {
	case *CreateHand:
		// An aggregate can only be created once.
		if a.currentHandID == typed.HandID {
			return ErrHandAlreadyCreated
		}
	}

	return nil
}

func (a *Aggregate) createEvent(cmd eventing.Command) error {
	switch cmd := cmd.(type) {
	case *CreateHand:
		a.AppendEvent(EventTypeHandCreated, &HandCreated{
			HandID: cmd.HandID,
			Cards:  cmd.Cards,
		}, TimeNow())

	default:
		return fmt.Errorf("could not handle command: %s", cmd.CommandType())
	}

	return nil
}

// HandleCommand implements the HandleCommand method of the
// eventing.CommandHandler interface.
func (a *Aggregate) HandleCommand(ctx context.Context, cmd eventing.Command) error {
	if err := a.validateCommand(cmd); err != nil {
		return err
	}

	if err := a.createEvent(cmd); err != nil {
		return err
	}

	return nil
}
