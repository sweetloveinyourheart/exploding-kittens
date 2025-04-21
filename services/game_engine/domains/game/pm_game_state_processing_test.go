package game_test

import (
	"context"
	"time"

	"github.com/avast/retry-go"
	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/mock"

	"github.com/sweetloveinyourheart/exploding-kittens/pkg/constants/cards"
	"github.com/sweetloveinyourheart/exploding-kittens/pkg/domain-eventing/command_handler/bus"
	retrymw "github.com/sweetloveinyourheart/exploding-kittens/pkg/domain-eventing/middleware/command_hander"

	"github.com/sweetloveinyourheart/exploding-kittens/pkg/domain-eventing/repo/version"
	deskDomain "github.com/sweetloveinyourheart/exploding-kittens/pkg/domains/desk"
	gameDomain "github.com/sweetloveinyourheart/exploding-kittens/pkg/domains/game"
	"github.com/sweetloveinyourheart/exploding-kittens/services/game_engine/models"
	"github.com/sweetloveinyourheart/exploding-kittens/services/game_engine/repos"
)

func (gs *GameSuite) TestGameStateProcessing_HandleGameCreated() {
	gs.setupEnvironment()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	cards := []repos.CardDetail{
		{Card: models.Card{CardID: uuid.FromStringOrNil("123e4567-e89b-12d3-a456-426655440001"), Name: cards.ExplodingKitten, Quantity: 4}},
		{Card: models.Card{CardID: uuid.FromStringOrNil("123e4567-e89b-12d3-a456-426655440002"), Name: cards.Defuse, Quantity: 4}},
		{Card: models.Card{CardID: uuid.FromStringOrNil("123e4567-e89b-12d3-a456-426655440003"), Name: cards.Attack, Quantity: 4}},
		{Card: models.Card{CardID: uuid.FromStringOrNil("123e4567-e89b-12d3-a456-426655440004"), Name: cards.Nope, Quantity: 5}},
		{Card: models.Card{CardID: uuid.FromStringOrNil("123e4567-e89b-12d3-a456-426655440005"), Name: cards.SeeTheFuture, Quantity: 5}},
		{Card: models.Card{CardID: uuid.FromStringOrNil("123e4567-e89b-12d3-a456-426655440006"), Name: cards.Shuffle, Quantity: 4}},
		{Card: models.Card{CardID: uuid.FromStringOrNil("123e4567-e89b-12d3-a456-426655440007"), Name: cards.Skip, Quantity: 4}},
		{Card: models.Card{CardID: uuid.FromStringOrNil("123e4567-e89b-12d3-a456-426655440008"), Name: cards.Favor, Quantity: 4}},
		{Card: models.Card{CardID: uuid.FromStringOrNil("123e4567-e89b-12d3-a456-426655440009"), Name: cards.BeardCat, Quantity: 4}},
		{Card: models.Card{CardID: uuid.FromStringOrNil("123e4567-e89b-12d3-a456-426655440010"), Name: cards.Catermelon, Quantity: 4}},
		{Card: models.Card{CardID: uuid.FromStringOrNil("123e4567-e89b-12d3-a456-426655440011"), Name: cards.HairyPotatoCat, Quantity: 4}},
		{Card: models.Card{CardID: uuid.FromStringOrNil("123e4567-e89b-12d3-a456-426655440012"), Name: cards.RainbowRalphingCat, Quantity: 4}},
	}

	gs.mockCardRepository.On("GetCards", mock.Anything).Return(cards, nil)
	gs.NoError(nil, nil)

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

	gs.Eventually(func() bool {
		gameState, err := gameRepo.Find(findCtx, gameID.String())
		deskState, err := deskRepo.Find(findCtx, gameState.Desk.String())

		shuffled := isCardsShuffled(cards, deskState.Cards)

		return err == nil && shuffled && gameState.Desk != uuid.Nil
	}, 5*time.Second, 10*time.Millisecond)
}

func isCardsShuffled(original []repos.CardDetail, shuffled []uuid.UUID) bool {
	for i := range original {
		if original[i].Card.CardID != shuffled[i] {
			return true // cards are same, but order differs = shuffled
		}
	}

	return false // same order, not shuffled
}
