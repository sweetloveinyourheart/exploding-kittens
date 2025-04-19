package game

import (
	"context"

	"github.com/gofrs/uuid"
	"github.com/sweetloveinyourheart/exploding-kittens/pkg/domain-eventing/common"
	"github.com/sweetloveinyourheart/exploding-kittens/pkg/timeutil"
)

type Projector struct {
	Projector *GameProjector
}

func NewProjector() *Projector {
	p := &Projector{}
	p.Projector = NewGameProjection(p)
	return p
}

var _ AllEventsProjector = (*Projector)(nil)
var _ AfterEntityHandler = (*Projector)(nil)

func (p *Projector) ProjectorType() common.ProjectorType {
	return common.ProjectorType(AggregateType.String())
}

func (p *Projector) AfterHandleEvent(ctx context.Context, event common.Event, data any, entity *Game) (*Game, error) {
	entity.UpdatedAt = timeutil.NowRoundedForGranularity()
	return entity, nil
}

func (p *Projector) Project(ctx context.Context, event common.Event, entity *Game) (*Game, error) {
	return p.Projector.Project(ctx, event, entity)
}

func (p *Projector) HandleGameCreated(ctx context.Context, event common.Event, data *GameCreated, entity *Game) (*Game, error) {
	entity.GameID = data.GetGameID()

	entity.PlayerHands = make(map[uuid.UUID]uuid.UUID)
	for _, playerID := range data.GetPlayerIDs() {
		entity.PlayerHands[playerID] = uuid.Nil
	}

	return entity, nil
}
