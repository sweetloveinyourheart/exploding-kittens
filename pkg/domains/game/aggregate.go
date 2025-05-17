package game

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
	type gameIDer interface {
		GetGameID() uuid.UUID
	}

	subjectFunc := func(ctx context.Context, event common.Event) string {
		if event == nil {
			return ""
		}

		if a, ok := event.Data().(gameIDer); ok {
			return fmt.Sprintf("%s.%s.%s", event.AggregateType(), a.GetGameID(), event.AggregateID())
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
				eventing.NewEventSubjectToken("game_id", "Game ID", uuid.Nil, 1),
				eventing.NewEventSubjectToken("aggregate_id", "Aggregate ID", uuid.Nil, 2),
			}
		}
		if a, ok := event.Data().(gameIDer); ok {
			return []common.EventSubjectToken{
				eventing.NewEventSubjectToken("aggregate_type", "Aggregate Type", AggregateType, 0),
				eventing.NewEventSubjectToken("game_id", "Game ID", a.GetGameID(), 1),
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

const AggregateType = common.AggregateType("game")

type Aggregate struct {
	*aggregate.AggregateBase

	actived       bool
	currentGameID uuid.UUID
	playerTurn    uuid.UUID
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
	case *CreateGame:
		// An aggregate can only be created once.
		if a.currentGameID == typed.GameID {
			return ErrGameAlreadyCreated
		}

		if len(typed.PlayerIDs) < 2 {
			return ErrNotEnoughUsersToPlay
		}

		if len(typed.PlayerIDs) > 5 {
			return ErrTooManyUsersToPlay
		}

	case *InitializeGame:
		if a.currentGameID != typed.GameID {
			return ErrGameNotFound
		}

		if a.actived {
			return ErrGameAlreadyInitialized
		}

	case *StartTurn:
		if a.currentGameID != typed.GameID {
			return ErrGameNotFound
		}

		if a.playerTurn == typed.PlayerID {
			return ErrPlayerIsAlreadyInTheirTurn
		}

	case *FinishTurn:
		if a.currentGameID != typed.GameID {
			return ErrGameNotFound
		}

		if a.playerTurn != typed.PlayerID {
			return ErrPlayerNotInTheirTurn
		}

	case *ReverseTurn:
		if a.currentGameID != typed.GameID {
			return ErrGameNotFound
		}

		if a.playerTurn != typed.PlayerID {
			return ErrPlayerNotInTheirTurn
		}

	case *PlayCards:
		if a.currentGameID != typed.GameID {
			return ErrGameNotFound
		}

		if a.playerTurn != typed.PlayerID {
			return ErrPlayerNotInTheirTurn
		}

	case *CreateAction:
		if a.currentGameID != typed.GameID {
			return ErrGameNotFound
		}

	case *ExecuteAction:
		if a.currentGameID != typed.GameID {
			return ErrGameNotFound
		}

	case *SelectAffectedPlayer:
		if a.currentGameID != typed.GameID {
			return ErrGameNotFound
		}

	}
	return nil
}

func (a *Aggregate) createEvent(cmd eventing.Command) error {
	switch cmd := cmd.(type) {
	case *CreateGame:
		a.AppendEvent(EventTypeGameCreated, &GameCreated{
			GameID:    cmd.GameID,
			PlayerIDs: cmd.PlayerIDs,
		}, TimeNow())

	case *InitializeGame:
		a.AppendEvent(EventTypeGameInitialized, &GameInitialized{
			GameID:      cmd.GameID,
			Desk:        cmd.Desk,
			PlayerHands: cmd.PlayerHands,
		}, TimeNow())

	case *StartTurn:
		a.AppendEvent(EventTypeTurnStarted, &TurnStarted{
			GameID:   cmd.GameID,
			PlayerID: cmd.PlayerID,
		}, TimeNow())

	case *FinishTurn:
		a.AppendEvent(EventTypeTurnFinished, &TurnFinished{
			GameID:   cmd.GameID,
			PlayerID: cmd.PlayerID,
		}, TimeNow())

	case *ReverseTurn:
		a.AppendEvent(EventTypeTurnReversed, &TurnReversed{
			GameID:   cmd.GameID,
			PlayerID: cmd.PlayerID,
		}, TimeNow())

	case *PlayCards:
		a.AppendEvent(EventTypeCardsPlayed, &CardsPlayed{
			GameID:   cmd.GameID,
			PlayerID: cmd.PlayerID,
			CardIDs:  cmd.CardIDs,
		}, TimeNow())

	case *CreateAction:
		a.AppendEvent(EventTypeActionCreated, &ActionCreated{
			GameID: cmd.GameID,
			Effect: cmd.Effect,
		}, TimeNow())

	case *SelectAffectedPlayer:
		a.AppendEvent(EventTypeAffectedPlayerSelected, &AffectedPlayerSelected{
			GameID:   cmd.GameID,
			PlayerID: cmd.PlayerID,
		}, TimeNow())

	case *ExecuteAction:
		a.AppendEvent(EventTypeActionExecuted, &ActionExecuted{
			GameID: cmd.GameID,
			Effect: cmd.Effect,
			Args:   cmd.Args,
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
	case EventTypeGameCreated:
		data, ok := event.Data().(*GameCreated)
		if !ok {
			return fmt.Errorf("could not apply event: %s", event.EventType())
		}

		a.currentGameID = data.GetGameID()

	case EventTypeGameInitialized:
		a.actived = true

	case EventTypeTurnStarted:
		data, ok := event.Data().(*TurnStarted)
		if !ok {
			return fmt.Errorf("could not apply event: %s", event.EventType())
		}
		a.playerTurn = data.GetPlayerID()

	case EventTypeTurnFinished:
		a.playerTurn = uuid.Nil

	case EventTypeTurnReversed:
		a.playerTurn = uuid.Nil

	case EventTypeCardsPlayed:
	case EventTypeActionCreated:
	case EventTypeAffectedPlayerSelected:
	case EventTypeActionExecuted:
	}

	return nil
}
