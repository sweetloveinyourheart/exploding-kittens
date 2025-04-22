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
	handDomain "github.com/sweetloveinyourheart/exploding-kittens/pkg/domains/hand"

	dataProviderProto "github.com/sweetloveinyourheart/exploding-kittens/proto/code/dataprovider/go"
)

func (gs *GameSuite) TestGameStateProcessing_HandleGameCreated() {
	gs.setupEnvironment()
	cards := gs.prepareCards()

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
		gameState, gameStateErr := gameRepo.Find(findCtx, gameID.String())
		deskState, deskStateErr := deskRepo.Find(findCtx, gameState.Desk.String())

		validDesk := isValidDesk(cards, deskState.Cards, len(gameState.PlayerHands))

		return gameStateErr == nil && deskStateErr == nil && validDesk && gameState.Desk != uuid.Nil
	}, 5*time.Second, 10*time.Millisecond)

	gs.Eventually(func() bool {
		for _, playerID := range playerIDs {
			playerHandID := handDomain.NewPlayerHandID(gameID, playerID)
			handState, err := handRepo.Find(findCtx, playerHandID.String())

			isValid := err == nil && len(handState.GetCards()) == 8
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
		if card.Name == cards.Defuse {
			defuseQuantity = int(card.Quantity)
		}
	}

	explodingKittenCards := 0
	defuseCards := 0

	for _, cardID := range deskCards {
		if card, ok := cardMap[cardID]; ok {
			switch card.Name {
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
