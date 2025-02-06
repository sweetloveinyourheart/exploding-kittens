package lobby

import (
	"context"
	"fmt"

	"github.com/cockroachdb/errors"
	"github.com/gofrs/uuid"

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
				eventing.NewEventSubjectToken("studio_table_id", "Studio Table ID", uuid.Nil, 1),
				eventing.NewEventSubjectToken("aggregate_id", "Edge Table ID", uuid.Nil, 2),
			}
		}
		if a, ok := event.Data().(lobbyIDer); ok {
			return []common.EventSubjectToken{
				eventing.NewEventSubjectToken("aggregate_type", "Aggregate Type", AggregateType, 0),
				eventing.NewEventSubjectToken("studio_table_id", "Studio Table ID", a.GetLobbyID(), 1),
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
			CreatedAt:    timeutil.NowRoundedForGranularity(),
			UpdatedAt:    timeutil.NowRoundedForGranularity(),
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

	default:
		return errors.WithStack(fmt.Errorf("could not apply event: %s", event.EventType()))
	}
	return nil
}
