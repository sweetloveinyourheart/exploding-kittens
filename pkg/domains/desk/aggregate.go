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
	cardCount     int
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

	case *ShuffleDesk:
		// An aggregate can only be shuffled once.
		if a.currentDeskID != typed.DeskID {
			return ErrDeskNotAvailable
		}

	case *DiscardCards:
		// An aggregate can only be discarded once.
		if a.currentDeskID != typed.DeskID {
			return ErrDeskNotAvailable
		}

		if len(typed.CardIDs) == 0 {
			return ErrNoCardsToDiscard
		}

	case *PeekCards:
		// An aggregate can only be peeked once.
		if a.currentDeskID != typed.DeskID {
			return ErrDeskNotAvailable
		}

		if typed.Count <= 0 {
			return ErrInvalidPeekCount
		}

	case *DrawCard:
		// An aggregate can only be drawn from once.
		if a.currentDeskID != typed.DeskID {
			return ErrDeskNotAvailable
		}

		if a.cardCount <= 0 {
			return ErrNoCardToDraw
		}

	case *InsertCard:
		// An aggregate can only be inserted into once.
		if a.currentDeskID != typed.DeskID {
			return ErrDeskNotAvailable
		}

		if typed.Index < 0 || typed.Index > a.cardCount {
			return ErrInvalidInsertIndex
		}

	}

	return nil
}

func (a *Aggregate) createEvent(cmd eventing.Command) error {
	switch cmd := cmd.(type) {
	case *CreateDesk:
		a.AppendEvent(EventTypeDeskCreated, &DeskCreated{
			DeskID:  cmd.DeskID,
			CardIDs: cmd.CardIDs,
		}, TimeNow())

	case *ShuffleDesk:
		a.AppendEvent(EventTypeDeskShuffled, &DeskShuffled{
			DeskID: cmd.DeskID,
		}, TimeNow())

	case *DiscardCards:
		a.AppendEvent(EventTypeCardsDiscarded, &CardsDiscarded{
			DeskID:  cmd.DeskID,
			CardIDs: cmd.CardIDs,
		}, TimeNow())

	case *PeekCards:
		a.AppendEvent(EventTypeCardsPeeked, &CardsPeeked{
			DeskID: cmd.DeskID,
			Count:  cmd.Count,
		}, TimeNow())

	case *DrawCard:
		a.AppendEvent(EventTypeCardDrawn, &CardDrawn{
			DeskID:        cmd.DeskID,
			GameID:        cmd.GameID,
			PlayerID:      cmd.PlayerID,
			CanFinishTurn: cmd.CanFinishTurn,
		}, TimeNow())

	case *InsertCard:
		a.AppendEvent(EventTypeCardInserted, &CardInserted{
			DeskID: cmd.DeskID,
			CardID: cmd.CardID,
			Index:  cmd.Index,
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

// ApplyEvent implements the ApplyEvent method of the
// eventing.Aggregate interface.
func (a *Aggregate) ApplyEvent(ctx context.Context, event common.Event) error {
	switch event.EventType() {
	case EventTypeDeskCreated:
		data, ok := event.Data().(*DeskCreated)
		if !ok {
			return fmt.Errorf("could not apply event: %s", event.EventType())
		}
		a.currentDeskID = data.GetDeskID()
		a.cardCount = len(data.GetCardIDs())

	case EventTypeDeskShuffled:
	case EventTypeCardsDiscarded:
	case EventTypeCardsPeeked:
	case EventTypeCardDrawn:
		a.cardCount--
	case EventTypeCardInserted:
		a.cardCount++
	}

	return nil
}
