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
	players := make([]Player, 0)
	for _, playerID := range data.GetPlayerIDs() {
		players = append(players, Player{PlayerID: playerID, Active: true})
	}

	entity.GameID = data.GetGameID()
	entity.Players = players
	entity.DiscardPile = make([]uuid.UUID, 0)
	entity.CreatedAt = timeutil.NowRoundedForGranularity()

	return entity, nil
}

func (p *Projector) HandleGameInitialized(ctx context.Context, event common.Event, data *GameInitialized, entity *Game) (*Game, error) {
	entity.GamePhase = GAME_PHASE_INITIALIZING
	entity.Desk = data.GetDesk()
	entity.PlayerHands = data.GetPlayerHands()

	return entity, nil
}

func (p *Projector) HandleTurnStarted(ctx context.Context, event common.Event, data *TurnStarted, entity *Game) (*Game, error) {
	entity.GamePhase = GAME_PHASE_TURN_START
	entity.PlayerTurn = data.GetPlayerID()

	return entity, nil
}
func (p *Projector) HandleTurnFinished(ctx context.Context, event common.Event, data *TurnFinished, entity *Game) (*Game, error) {
	entity.GamePhase = GAME_PHASE_TURN_FINISH
	entity.PlayerTurn = uuid.Nil

	return entity, nil
}

func (p *Projector) HandleCardPlayed(ctx context.Context, event common.Event, data *CardPlayed, entity *Game) (*Game, error) {
	entity.GamePhase = GAME_PHASE_ACTION_PHASE
	entity.DiscardPile = append(entity.DiscardPile, data.GetCardIDs()...)

	return entity, nil
}
