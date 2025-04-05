package desk

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
	type deskIDer interface {
		GetDeskID() uuid.UUID
	}

	subjectFunc := func(ctx context.Context, event common.Event) string {
		if event == nil {
			return ""
		}

		if a, ok := event.Data().(deskIDer); ok {
			return fmt.Sprintf("%s.%s.%s", event.AggregateType(), a.GetDeskID(), event.AggregateID())
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
				eventing.NewEventSubjectToken("desk_id", "Desk ID", uuid.Nil, 1),
				eventing.NewEventSubjectToken("aggregate_id", "Aggregate ID", uuid.Nil, 2),
			}
		}
		if a, ok := event.Data().(deskIDer); ok {
			return []common.EventSubjectToken{
				eventing.NewEventSubjectToken("aggregate_type", "Aggregate Type", AggregateType, 0),
				eventing.NewEventSubjectToken("desk_id", "Desk ID", a.GetDeskID(), 1),
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

const AggregateType = common.AggregateType("desk")

type Aggregate struct {
	*aggregate.AggregateBase
	currentDeskID uuid.UUID
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
	case *CreateDesk:
		// An aggregate can only be created once.
		if a.currentDeskID == typed.DeskID {
			return ErrDeskAlreadyCreated
		}
	}

	return nil
}

func (a *Aggregate) createEvent(cmd eventing.Command) error {
	switch cmd := cmd.(type) {
	case *CreateDesk:
		a.AppendEvent(EventTypeDeskCreated, &DeskCreated{
			DeskID: cmd.DeskID,
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
