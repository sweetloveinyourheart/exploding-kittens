package lobby

import (
	"context"

	"github.com/gofrs/uuid"

	"github.com/sweetloveinyourheart/exploding-kittens/pkg/domain-eventing/common"
	"github.com/sweetloveinyourheart/exploding-kittens/pkg/timeutil"
)

type Projector struct {
	Projector *LobbyProjector
}

func NewProjector() *Projector {
	p := &Projector{}
	p.Projector = NewLobbyProjection(p)
	return p
}

var _ AllEventsProjector = (*Projector)(nil)
var _ AfterEntityHandler = (*Projector)(nil)

func (p *Projector) ProjectorType() common.ProjectorType {
	return common.ProjectorType(AggregateType.String())
}

func (p *Projector) AfterHandleEvent(ctx context.Context, event common.Event, data any, entity *Lobby) (*Lobby, error) {
	entity.UpdatedAt = timeutil.NowRoundedForGranularity()
	return entity, nil
}

func (p *Projector) Project(ctx context.Context, event common.Event, entity *Lobby) (*Lobby, error) {
	return p.Projector.Project(ctx, event, entity)
}

func (p *Projector) HandleLobbyCreated(ctx context.Context, event common.Event, data *LobbyCreated, entity *Lobby) (*Lobby, error) {
	entity.LobbyID = data.GetLobbyID()
	entity.LobbyCode = data.GetLobbyCode()
	entity.LobbyName = data.GetLobbyName()
	entity.HostUserID = data.GetHostUserID()
	entity.CreatedAt = timeutil.NowRoundedForGranularity()
	entity.Participants = []uuid.UUID{
		data.GetHostUserID(),
	}

	return entity, nil
}

func (p *Projector) HandleLobbyJoined(ctx context.Context, event common.Event, data *LobbyJoined, entity *Lobby) (*Lobby, error) {
	entity.Participants = append(entity.Participants, data.GetUserID())

	return entity, nil
}

func (p *Projector) HandleLobbyLeft(ctx context.Context, event common.Event, data *LobbyLeft, entity *Lobby) (*Lobby, error) {
	for i, participant := range entity.Participants {
		if participant == data.GetUserID() {
			entity.Participants = append(entity.Participants[:i], entity.Participants[i+1:]...)
			break
		}
	}

	return entity, nil
}
