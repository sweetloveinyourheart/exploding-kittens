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
	entity.CreatedAt = timeutil.NowRoundedForGranularity()

	return entity, nil
}

func (p *Projector) HandleGameInitialized(ctx context.Context, event common.Event, data *GameInitialized, entity *Game) (*Game, error) {
	entity.GameID = data.GetGameID()
	entity.GamePhase = GAME_PHASE_INITIALIZING
	entity.DeskID = data.GetDeskID()
	entity.PlayerHands = data.GetPlayerHands()

	return entity, nil
}

func (p *Projector) HandleGameStarted(ctx context.Context, event common.Event, data *GameStarted, entity *Game) (*Game, error) {
	entity.GameID = data.GetGameID()
	return entity, nil
}

func (p *Projector) HandleTurnStarted(ctx context.Context, event common.Event, data *TurnStarted, entity *Game) (*Game, error) {
	entity.GameID = data.GetGameID()
	entity.GamePhase = GAME_PHASE_TURN_START
	entity.PlayerTurn = data.GetPlayerID()

	return entity, nil
}

func (p *Projector) HandleTurnFinished(ctx context.Context, event common.Event, data *TurnFinished, entity *Game) (*Game, error) {
	entity.GameID = data.GetGameID()
	entity.GamePhase = GAME_PHASE_TURN_FINISH
	entity.PlayerTurn = uuid.Nil

	return entity, nil
}

func (p *Projector) HandleTurnReversed(ctx context.Context, event common.Event, data *TurnReversed, entity *Game) (*Game, error) {
	entity.GameID = data.GetGameID()
	entity.GamePhase = GAME_PHASE_TURN_FINISH
	entity.PlayerTurn = data.GetPlayerID()

	return entity, nil
}

func (p *Projector) HandleGameFinished(ctx context.Context, event common.Event, data *GameFinished, entity *Game) (*Game, error) {
	entity.GameID = data.GetGameID()
	entity.GamePhase = GAME_PHASE_GAME_FINISH
	entity.PlayerTurn = uuid.Nil
	entity.WinnerID = data.GetWinnerID()

	return entity, nil
}

func (p *Projector) HandleCardsPlayed(ctx context.Context, event common.Event, data *CardsPlayed, entity *Game) (*Game, error) {
	entity.GameID = data.GetGameID()
	entity.GamePhase = GAME_PHASE_ACTION_PHASE
	return entity, nil
}

func (p *Projector) HandleActionCreated(ctx context.Context, event common.Event, data *ActionCreated, entity *Game) (*Game, error) {
	entity.GameID = data.GetGameID()
	entity.GamePhase = GAME_PHASE_ACTION_PHASE
	entity.ExecutingAction = data.GetEffect()

	return entity, nil
}

func (p *Projector) HandleAffectedPlayerSelected(ctx context.Context, event common.Event, data *AffectedPlayerSelected, entity *Game) (*Game, error) {
	entity.GameID = data.GetGameID()
	entity.GamePhase = GAME_PHASE_ACTION_PHASE
	entity.AffectedPlayer = data.GetPlayerID()

	return entity, nil
}

func (p *Projector) HandleActionExecuted(ctx context.Context, event common.Event, data *ActionExecuted, entity *Game) (*Game, error) {
	entity.GameID = data.GetGameID()
	entity.GamePhase = GAME_PHASE_TURN_START
	entity.ExecutingAction = ""
	entity.AffectedPlayer = uuid.Nil

	return entity, nil
}

func (p *Projector) HandleCardDrawn(ctx context.Context, event common.Event, data *CardDrawn, entity *Game) (*Game, error) {
	entity.GameID = data.GetGameID()

	return entity, nil
}

func (p *Projector) HandleExplodingDrawn(ctx context.Context, event common.Event, data *ExplodingDrawn, entity *Game) (*Game, error) {
	entity.GameID = data.GetGameID()
	entity.GamePhase = GAME_PHASE_EXPLODING_DRAWN

	return entity, nil
}

func (p *Projector) HandleExplodingDefused(ctx context.Context, event common.Event, data *ExplodingDefused, entity *Game) (*Game, error) {
	entity.GameID = data.GetGameID()
	entity.GamePhase = GAME_PHASE_EXPLODING_DEFUSED

	return entity, nil
}

func (p *Projector) HandlePlayerEliminated(ctx context.Context, event common.Event, data *PlayerEliminated, entity *Game) (*Game, error) {
	entity.GameID = data.GetGameID()
	entity.GamePhase = GAME_PHASE_PLAYER_ELIMINATED
	entity.PlayerTurn = uuid.Nil

	players := entity.Players
	for i, player := range players {
		if player.GetPlayerID() == data.GetPlayerID() && player.Active {
			players[i].Active = false
			break
		}
	}
	entity.Players = players

	return entity, nil
}

func (p *Projector) HandleKittenPlanted(ctx context.Context, event common.Event, data *KittenPlanted, entity *Game) (*Game, error) {
	entity.GameID = data.GetGameID()

	return entity, nil
}
