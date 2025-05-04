package game_test

import (
	"context"
	"time"

	"github.com/avast/retry-go"
	"github.com/gofrs/uuid"
	"github.com/sweetloveinyourheart/exploding-kittens/pkg/constants/cards"
	"github.com/sweetloveinyourheart/exploding-kittens/pkg/domain-eventing/command_handler/bus"
	retrymw "github.com/sweetloveinyourheart/exploding-kittens/pkg/domain-eventing/middleware/command_hander"
	"github.com/sweetloveinyourheart/exploding-kittens/pkg/domain-eventing/repo/version"
	deskDomain "github.com/sweetloveinyourheart/exploding-kittens/pkg/domains/desk"
	gameDomain "github.com/sweetloveinyourheart/exploding-kittens/pkg/domains/game"
)

func (gs *GameSuite) TestGamePlayExecutor_HandleCardPlay_Skip() {
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

	findCtx, cancel := version.NewContextWithMinVersionWait(ctx, 1)
	defer cancel()

	gs.Eventually(func() bool {
		gameState, gameStateErr := gameRepo.Find(findCtx, gameID.String())
		return gameStateErr == nil && gameState.GetPlayerTurn() == player01
	}, 5*time.Second, 10*time.Millisecond)

	err = commandBus.HandleCommand(ctx, &gameDomain.PlayCard{
		GameID:   gameID,
		PlayerID: player01,
		CardIDs:  []uuid.UUID{uuid.Must(uuid.FromString(cardsCodeMap[cards.Skip].GetCardId()))},
	})
	gs.NoError(err)

	gs.Eventually(func() bool {
		gameState, gameStateErr := gameRepo.Find(findCtx, gameID.String())
		return gameStateErr == nil && gameState.GetPlayerTurn() == player02
	}, 5*time.Second, 10*time.Millisecond)
}

func (gs *GameSuite) TestGamePlayExecutor_HandleCardPlay_ShuffleDesk() {
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

	deskRepo, err := deskDomain.CreateNATSRepoDesk(ctx, "test")
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

	findCtx, cancel := version.NewContextWithMinVersionWait(ctx, 1)
	defer cancel()

	var previousDeskState *deskDomain.Desk
	gs.Eventually(func() bool {
		gameState, gameStateErr := gameRepo.Find(findCtx, gameID.String())
		deskState, deskStateErr := deskRepo.Find(findCtx, gameState.Desk.String())
		previousDeskState = deskState
		return gameStateErr == nil && deskStateErr == nil && gameState.Desk != uuid.Nil
	}, 5*time.Second, 10*time.Millisecond)

	err = commandBus.HandleCommand(ctx, &gameDomain.PlayCard{
		GameID:   gameID,
		PlayerID: player01,
		CardIDs:  []uuid.UUID{uuid.Must(uuid.FromString(cardsCodeMap[cards.Shuffle].GetCardId()))},
	})
	gs.NoError(err)

	findCtx, cancel = version.NewContextWithMinVersionWait(ctx, 2)
	defer cancel()

	var afterShuffleDeskState *deskDomain.Desk
	gs.Eventually(func() bool {
		gameState, gameStateErr := gameRepo.Find(findCtx, gameID.String())
		deskState, deskStateErr := deskRepo.Find(findCtx, gameState.Desk.String())
		afterShuffleDeskState = deskState
		return gameStateErr == nil && deskStateErr == nil && gameState.Desk != uuid.Nil && gameState.GetPlayerTurn() == player01
	}, 5*time.Second, 10*time.Millisecond)

	gs.NotEqual(previousDeskState.Cards, afterShuffleDeskState.Cards)
}
