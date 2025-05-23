package game_test

import (
	"context"
	"runtime"
	"time"

	"github.com/avast/retry-go"
	"github.com/gofrs/uuid"

	"github.com/sweetloveinyourheart/exploding-kittens/pkg/constants/cards"
	"github.com/sweetloveinyourheart/exploding-kittens/pkg/domain-eventing/command_handler/bus"
	retrymw "github.com/sweetloveinyourheart/exploding-kittens/pkg/domain-eventing/middleware/command_hander"
	gameDomain "github.com/sweetloveinyourheart/exploding-kittens/pkg/domains/game"
)

func (gs *GameSuite) TestGameExplodingProcessManager_HandleExplodingCardDrawn() {
	gs.setupEnvironment()
	_, _, cardsCodeMap := gs.prepareCards()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	mw := retrymw.NewCommandHandlerMiddleware(retry.Attempts(4), retry.MaxDelay(1*time.Second))
	commandBus := bus.NewCommandHandler()

	err := gameDomain.AddNATSGameCommandHandlers(ctx, "test", commandBus, mw)
	gs.NoError(err)

	gameRepo, err := gameDomain.CreateNATSRepoGames(ctx, "test")
	gs.NoError(err)

	gameID := uuid.Must(uuid.NewV7())
	player01 := uuid.Must(uuid.NewV7())
	player02 := uuid.Must(uuid.NewV7())
	playerIDs := []uuid.UUID{
		player01,
		player02,
	}

	err = commandBus.HandleCommand(ctx, &gameDomain.CreateGame{
		GameID:    gameID,
		PlayerIDs: playerIDs,
	})
	gs.NoError(err)

	gs.Eventually(func() bool {
		runtime.Gosched()

		gameState, gameStateErr := gameRepo.Find(ctx, gameID.String())
		return gameStateErr == nil && gameState.GetPlayerTurn() == player01
	}, 5*time.Second, 10*time.Millisecond)

	err = commandBus.HandleCommand(ctx, &gameDomain.DrawExplodingKitten{
		GameID:   gameID,
		PlayerID: player01,
		CardID:   uuid.Must(uuid.FromString(cardsCodeMap[cards.ExplodingKitten].GetCardId())),
	})
	gs.NoError(err)

	gs.Eventually(func() bool {
		runtime.Gosched()

		gameState, gameStateErr := gameRepo.Find(ctx, gameID.String())
		return gameStateErr == nil && gameState.GetPlayerTurn() == player01 && gameState.GetGamePhase() == gameDomain.GAME_PHASE_EXPLODING_DRAWN
	}, 5*time.Second, 10*time.Millisecond)
}

func (gs *GameSuite) TestGameExplodingProcessManager_HandleDefuseExplodingCard() {
	gs.setupEnvironment()
	_, _, cardsCodeMap := gs.prepareCards()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	mw := retrymw.NewCommandHandlerMiddleware(retry.Attempts(4), retry.MaxDelay(1*time.Second))
	commandBus := bus.NewCommandHandler()

	err := gameDomain.AddNATSGameCommandHandlers(ctx, "test", commandBus, mw)
	gs.NoError(err)

	gameRepo, err := gameDomain.CreateNATSRepoGames(ctx, "test")
	gs.NoError(err)

	gameID := uuid.Must(uuid.NewV7())
	player01 := uuid.Must(uuid.NewV7())
	player02 := uuid.Must(uuid.NewV7())
	playerIDs := []uuid.UUID{
		player01,
		player02,
	}

	err = commandBus.HandleCommand(ctx, &gameDomain.CreateGame{
		GameID:    gameID,
		PlayerIDs: playerIDs,
	})
	gs.NoError(err)

	gs.Eventually(func() bool {
		runtime.Gosched()

		gameState, gameStateErr := gameRepo.Find(ctx, gameID.String())
		return gameStateErr == nil && gameState.GetPlayerTurn() == player01
	}, 5*time.Second, 10*time.Millisecond)

	err = commandBus.HandleCommand(ctx, &gameDomain.DrawExplodingKitten{
		GameID:   gameID,
		PlayerID: player01,
		CardID:   uuid.Must(uuid.FromString(cardsCodeMap[cards.ExplodingKitten].GetCardId())),
	})
	gs.NoError(err)

	gs.Eventually(func() bool {
		runtime.Gosched()

		gameState, gameStateErr := gameRepo.Find(ctx, gameID.String())

		return gameStateErr == nil &&
			gameState.GetPlayerTurn() == player01 &&
			gameState.GetGamePhase() == gameDomain.GAME_PHASE_EXPLODING_DRAWN

	}, 5*time.Second, 10*time.Millisecond)

	err = commandBus.HandleCommand(ctx, &gameDomain.DefuseExplodingKitten{
		GameID:   gameID,
		PlayerID: player01,
		CardID:   uuid.Must(uuid.FromString(cardsCodeMap[cards.Defuse].GetCardId())),
	})
	gs.NoError(err)

	gs.Eventually(func() bool {
		runtime.Gosched()

		gameState, gameStateErr := gameRepo.Find(ctx, gameID.String())
		return gameStateErr == nil && gameState.GetPlayerTurn() == player01 && gameState.GetGamePhase() == gameDomain.GAME_PHASE_EXPLODING_DEFUSED
	}, 5*time.Second, 10*time.Millisecond)

	err = commandBus.HandleCommand(ctx, &gameDomain.PlantTheKitten{
		GameID:   gameID,
		PlayerID: player01,
		Index:    1,
	})
	gs.NoError(err)

	gs.Eventually(func() bool {
		runtime.Gosched()

		gameState, gameStateErr := gameRepo.Find(ctx, gameID.String())
		return gameStateErr == nil &&
			gameState.GetPlayerTurn() == player02 &&
			gameState.GetGamePhase() == gameDomain.GAME_PHASE_TURN_START
	}, 5*time.Second, 10*time.Millisecond)
}
