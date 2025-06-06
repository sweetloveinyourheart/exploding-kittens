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

	"github.com/sweetloveinyourheart/exploding-kittens/pkg/domain-eventing/repo/version"
	deskDomain "github.com/sweetloveinyourheart/exploding-kittens/pkg/domains/desk"
	gameDomain "github.com/sweetloveinyourheart/exploding-kittens/pkg/domains/game"
	handDomain "github.com/sweetloveinyourheart/exploding-kittens/pkg/domains/hand"

	dataProviderProto "github.com/sweetloveinyourheart/exploding-kittens/proto/code/dataprovider/go"
)

func (gs *GameSuite) TestGameStateProcessing_HandleGameCreated() {
	gs.setupEnvironment()
	cards, _, _ := gs.prepareCards()

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

	handRepo, err := handDomain.CreateNATSRepoHand(ctx, "test")
	gs.NoError(err)

	gameID := uuid.Must(uuid.NewV7())
	playerIDs := []uuid.UUID{
		uuid.Must(uuid.NewV7()),
		uuid.Must(uuid.NewV7()),
		uuid.Must(uuid.NewV7()),
	}

	err = commandBus.HandleCommand(ctx, &gameDomain.CreateGame{
		GameID:    gameID,
		PlayerIDs: playerIDs,
	})
	gs.NoError(err)

	findCtx, cancel := version.NewContextWithMinVersionWait(ctx, 1)
	defer cancel()

	gs.Eventually(func() bool {
		runtime.Gosched()

		gameState, gameStateErr := gameRepo.Find(findCtx, gameID.String())
		gs.NoError(gameStateErr)

		deskState, deskStateErr := deskRepo.Find(findCtx, gameState.DeskID.String())
		gs.NoError(deskStateErr)

		validDesk := isValidDesk(cards, deskState.GetCardIDs(), len(gameState.PlayerHands))

		return gameStateErr == nil && deskStateErr == nil && validDesk && gameState.DeskID != uuid.Nil
	}, 5*time.Second, 10*time.Millisecond)

	gs.Eventually(func() bool {
		runtime.Gosched()

		for _, playerID := range playerIDs {
			playerHandID := handDomain.NewPlayerHandID(gameID, playerID)
			handState, err := handRepo.Find(findCtx, playerHandID.String())
			gs.NoError(err)

			isValid := err == nil && len(handState.GetCardIDs()) == 8
			if !isValid {
				return false
			}
		}

		return true
	}, 5*time.Second, 10*time.Millisecond)
}

func isValidDesk(original []*dataProviderProto.Card, deskCards []uuid.UUID, playerNum int) bool {
	cardMap := make(map[uuid.UUID]*dataProviderProto.Card, len(original))
	defuseQuantity := 0

	for _, card := range original {
		cardMap[uuid.FromStringOrNil(card.CardId)] = card
		if card.Code == cards.Defuse {
			defuseQuantity = int(card.Quantity)
		}
	}

	explodingKittenCards := 0
	defuseCards := 0

	for _, cardID := range deskCards {
		if card, ok := cardMap[cardID]; ok {
			switch card.Code {
			case cards.ExplodingKitten:
				explodingKittenCards++
			case cards.Defuse:
				defuseCards++
			}
		}
	}

	// Ensure that the number of Exploding Kitten and Defuse cards is valid
	return explodingKittenCards == playerNum-1 && defuseCards == defuseQuantity-playerNum
}

func (gs *GameSuite) TestGameStateProcessing_HandlePlayerEliminated() {
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
	player03 := uuid.Must(uuid.NewV7())
	playerIDs := []uuid.UUID{
		player01,
		player02,
		player03,
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

	err = commandBus.HandleCommand(ctx, &gameDomain.EliminatePlayer{
		GameID:   gameID,
		PlayerID: player01,
	})
	gs.NoError(err)

	gs.Eventually(func() bool {
		runtime.Gosched()

		gameState, gameStateErr := gameRepo.Find(ctx, gameID.String())
		return gameStateErr == nil &&
			gameState.GetPlayerTurn() == player02 &&
			gameState.GetGamePhase() == gameDomain.GAME_PHASE_TURN_START &&
			gameState.Players[0].Active == false
	}, 5*time.Second, 10*time.Millisecond)
}

func (gs *GameSuite) TestGameStateProcessing_HandleGameFinished() {
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

	err = commandBus.HandleCommand(ctx, &gameDomain.EliminatePlayer{
		GameID:   gameID,
		PlayerID: player01,
	})
	gs.NoError(err)

	gs.Eventually(func() bool {
		runtime.Gosched()

		gameState, gameStateErr := gameRepo.Find(ctx, gameID.String())
		return gameStateErr == nil &&
			gameState.GetGamePhase() == gameDomain.GAME_PHASE_GAME_FINISH &&
			gameState.WinnerID == player02
	}, 5*time.Second, 10*time.Millisecond)
}
