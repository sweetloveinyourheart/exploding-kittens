package lobby

import (
	"context"
	"fmt"

	"github.com/cockroachdb/errors"
	"github.com/gofrs/uuid"

	"slices"

	eventing "github.com/sweetloveinyourheart/exploding-kittens/pkg/domain-eventing"
	"github.com/sweetloveinyourheart/exploding-kittens/pkg/domain-eventing/aggregate"
	"github.com/sweetloveinyourheart/exploding-kittens/pkg/domain-eventing/common"
	"github.com/sweetloveinyourheart/exploding-kittens/pkg/timeutil"
)

var SubjectFactory common.EventSubject

func init() {
	eventing.RegisterAggregate[Aggregate, *Aggregate]()
	type lobbyIDer interface {
		GetLobbyID() uuid.UUID
	}

	subjectFunc := func(ctx context.Context, event common.Event) string {
		if event == nil {
			return ""
		}

		if a, ok := event.Data().(lobbyIDer); ok {
			return fmt.Sprintf("%s.%s.%s", event.AggregateType(), a.GetLobbyID(), event.AggregateID())
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
				eventing.NewEventSubjectToken("lobby_id", "Lobby ID", uuid.Nil, 1),
				eventing.NewEventSubjectToken("aggregate_id", "Aggregate ID", uuid.Nil, 2),
			}
		}
		if a, ok := event.Data().(lobbyIDer); ok {
			return []common.EventSubjectToken{
				eventing.NewEventSubjectToken("aggregate_type", "Aggregate Type", AggregateType, 0),
				eventing.NewEventSubjectToken("lobby_id", "Lobby ID", a.GetLobbyID(), 1),
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

const AggregateType = common.AggregateType("lobby")

type Aggregate struct {
	*aggregate.AggregateBase

	currentLobbyID uuid.UUID
	currentGameID  uuid.UUID
	actived        bool
	hostID         uuid.UUID
	playerIDs      []uuid.UUID
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
	case *CreateLobby:
		// An aggregate can only be created once.
		if a.currentLobbyID == typed.LobbyID {
			return ErrLobbyAlreadyCreated
		}

		if !a.currentLobbyID.IsNil() {
			return ErrLobbyInWaitingMode
		}
	case *JoinLobby:
		if !a.actived {
			return ErrLobbyNotAvailable
		}

		err := a.validatePlayerExists(typed.UserID)
		if err != nil {
			return err
		}
	case *LeaveLobby:
		if !a.actived {
			return ErrLobbyNotAvailable
		}

		err := a.validatePlayerNotExists(typed.UserID)
		if err != nil {
			return err
		}

	case *StartGame:
		if !a.actived {
			return ErrLobbyNotAvailable
		}

		if a.hostID != typed.HostUserID {
			return ErrHostUserNotRecognized
		}

		if a.currentGameID != uuid.Nil {
			return ErrGameIsAlreadyStarted
		}

		if len(a.playerIDs) < 2 {
			return ErrGameIsNotEnoughPlayer
		}

	default:
		// All other events require the aggregate to be created.
	}

	return nil
}

func (a *Aggregate) createEvent(cmd eventing.Command) error {
	switch cmd := cmd.(type) {
	case *CreateLobby:
		a.AppendEvent(EventTypeLobbyCreated, &LobbyCreated{
			LobbyID:      cmd.LobbyID,
			LobbyCode:    cmd.LobbyCode,
			LobbyName:    cmd.LobbyName,
			HostUserID:   cmd.HostUserID,
			Participants: []uuid.UUID{},
		}, TimeNow())
	case *JoinLobby:
		a.AppendEvent(EventTypeLobbyJoined, &LobbyJoined{
			LobbyID: cmd.LobbyID,
			UserID:  cmd.UserID,
		}, TimeNow())
	case *LeaveLobby:
		a.AppendEvent(EventTypeLobbyLeft, &LobbyLeft{
			LobbyID: cmd.LobbyID,
			UserID:  cmd.UserID,
		}, TimeNow())
	case *StartGame:
		a.AppendEvent(EventTypeGameStarted, &GameStarted{
			LobbyID:    cmd.LobbyID,
			HostUserID: cmd.HostUserID,
			GameID:     cmd.GameID,
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
	case EventTypeLobbyCreated:
		data, ok := event.Data().(*LobbyCreated)
		if !ok {
			return fmt.Errorf("could not apply event: %s", event.EventType())
		}

		a.currentLobbyID = data.LobbyID
		a.actived = true
		a.playerIDs = append(a.playerIDs, data.GetHostUserID())
		a.hostID = data.GetHostUserID()

	case EventTypeLobbyJoined:
		data, ok := event.Data().(*LobbyJoined)
		if !ok {
			return fmt.Errorf("could not apply event: %s", event.EventType())
		}
		a.playerIDs = append(a.playerIDs, data.GetUserID())

	case EventTypeLobbyLeft:
		data, ok := event.Data().(*LobbyLeft)
		if !ok {
			return fmt.Errorf("could not apply event: %s", event.EventType())
		}

		for i, id := range a.playerIDs {
			if id == data.GetUserID() {
				a.playerIDs = slices.Delete(a.playerIDs, i, i+1)
				break
			}
		}
		if len(a.playerIDs) == 0 {
			a.actived = false
		}

	case EventTypeGameStarted:
		data, ok := event.Data().(*GameStarted)
		if !ok {
			return fmt.Errorf("could not apply event: %s", event.EventType())
		}

		a.currentGameID = data.GetGameID()

	default:
		return errors.WithStack(fmt.Errorf("could not apply event: %s", event.EventType()))
	}
	return nil
}

func (a *Aggregate) validatePlayerExists(playerID uuid.UUID) error {
	if slices.Contains(a.playerIDs, playerID) {
		return errors.New("player is already exist")
	}
	return nil
}

func (a *Aggregate) validatePlayerNotExists(playerID uuid.UUID) error {
	if slices.Contains(a.playerIDs, playerID) {
		return nil
	}
	return errors.New("player does not exist")
}
