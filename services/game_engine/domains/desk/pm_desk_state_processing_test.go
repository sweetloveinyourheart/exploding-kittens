package desk_test

import (
	"context"
	"runtime"
	"time"

	"github.com/avast/retry-go"
	"github.com/gofrs/uuid"

	"github.com/sweetloveinyourheart/exploding-kittens/pkg/domain-eventing/command_handler/bus"
	retrymw "github.com/sweetloveinyourheart/exploding-kittens/pkg/domain-eventing/middleware/command_hander"
	"github.com/sweetloveinyourheart/exploding-kittens/pkg/domains/desk"
	"github.com/sweetloveinyourheart/exploding-kittens/pkg/domains/game"
)

func (hs *DeskSuite) TestDeskStateProcessing_DeskCreated_Successfully() {
	hs.setupEnvironment()
	cardIDs, _, _ := hs.prepareCards()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	mw := retrymw.NewCommandHandlerMiddleware(retry.Attempts(4), retry.MaxDelay(1*time.Second))
	commandBus := bus.NewCommandHandler()

	err := desk.AddNATSDeskCommandHandlers(ctx, "test", commandBus, mw)
	hs.NoError(err)

	deskRepo, err := desk.CreateNATSRepoDesk(ctx, "test")
	hs.NoError(err)

	cards := make([]uuid.UUID, 0, len(cardIDs))
	for _, card := range cardIDs {
		cards = append(cards, uuid.Must(uuid.FromString(card.GetCardId())))
	}

	deskID := uuid.Must(uuid.NewV7())
	err = commandBus.HandleCommand(ctx, &desk.CreateDesk{
		DeskID:  deskID,
		CardIDs: cards,
	})
	hs.NoError(err)

	hs.Eventually(func() bool {
		runtime.Gosched()

		deskState, deskStateErr := deskRepo.Find(ctx, deskID.String())
		return deskStateErr == nil && len(deskState.GetCardIDs()) > 0
	}, 5*time.Second, 10*time.Millisecond)
}

func (hs *DeskSuite) TestDeskStateProcessing_DeskShuffled_Successfully() {
	hs.setupEnvironment()
	cardIDs, _, _ := hs.prepareCards()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	mw := retrymw.NewCommandHandlerMiddleware(retry.Attempts(4), retry.MaxDelay(1*time.Second))
	commandBus := bus.NewCommandHandler()

	err := desk.AddNATSDeskCommandHandlers(ctx, "test", commandBus, mw)
	hs.NoError(err)

	deskRepo, err := desk.CreateNATSRepoDesk(ctx, "test")
	hs.NoError(err)

	cards := make([]uuid.UUID, 0, len(cardIDs))
	for _, card := range cardIDs {
		cards = append(cards, uuid.Must(uuid.FromString(card.GetCardId())))
	}

	deskID := uuid.Must(uuid.NewV7())
	err = commandBus.HandleCommand(ctx, &desk.CreateDesk{
		DeskID:  deskID,
		CardIDs: cards,
	})
	hs.NoError(err)

	hs.Eventually(func() bool {
		runtime.Gosched()

		deskState, deskStateErr := deskRepo.Find(ctx, deskID.String())
		return deskStateErr == nil && len(deskState.GetCardIDs()) > 0
	}, 5*time.Second, 10*time.Millisecond)

	err = commandBus.HandleCommand(ctx, &desk.ShuffleDesk{
		DeskID: deskID,
	})
	hs.NoError(err)

	hs.Eventually(func() bool {
		runtime.Gosched()

		deskState, deskStateErr := deskRepo.Find(ctx, deskID.String())
		return deskStateErr == nil && len(deskState.GetCardIDs()) > 0 && isDeskShuffled(deskState.GetCardIDs(), cards)
	}, 5*time.Second, 10*time.Millisecond)
}

func (hs *DeskSuite) TestDeskStateProcessing_DeskDrawn_Successfully() {
	hs.setupEnvironment()

	cardIDs, _, _ := hs.prepareCards()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	mw := retrymw.NewCommandHandlerMiddleware(retry.Attempts(4), retry.MaxDelay(1*time.Second))
	commandBus := bus.NewCommandHandler()

	err := game.AddNATSGameCommandHandlers(ctx, "test", commandBus, mw)
	hs.NoError(err)

	err = desk.AddNATSDeskCommandHandlers(ctx, "test", commandBus, mw)
	hs.NoError(err)

	gameRepo, err := game.CreateNATSRepoGames(ctx, "test")
	hs.NoError(err)

	deskRepo, err := desk.CreateNATSRepoDesk(ctx, "test")
	hs.NoError(err)

	cards := make([]uuid.UUID, 0, len(cardIDs))
	for _, card := range cardIDs {
		cards = append(cards, uuid.Must(uuid.FromString(card.GetCardId())))
	}

	gameID := uuid.Must(uuid.NewV7())
	player01 := uuid.Must(uuid.NewV7())
	player02 := uuid.Must(uuid.NewV7())
	playerIDs := []uuid.UUID{
		player01,
		player02,
	}

	err = commandBus.HandleCommand(ctx, &game.CreateGame{
		GameID:    gameID,
		PlayerIDs: playerIDs,
	})
	hs.NoError(err)

	var deskID uuid.UUID
	var deskCards int
	hs.Eventually(func() bool {
		runtime.Gosched()

		gameState, gameStateErr := gameRepo.Find(ctx, gameID.String())
		deskID = gameState.Desk

		deskState, deskStateErr := deskRepo.Find(ctx, deskID.String())
		deskCards = len(deskState.GetCardIDs())

		return gameStateErr == nil && gameState.GetPlayerTurn() == player01 &&
			deskStateErr == nil && deskCards > 0
	}, 5*time.Second, 10*time.Millisecond)

	err = commandBus.HandleCommand(ctx, &desk.DrawCards{
		DeskID:   deskID,
		GameID:   gameID,
		PlayerID: player01,
		Count:    1,
	})
	hs.NoError(err)

	hs.Eventually(func() bool {
		runtime.Gosched()

		deskState, deskStateErr := deskRepo.Find(ctx, deskID.String())
		return deskStateErr == nil && len(deskState.GetCardIDs()) == deskCards-1
	}, 5*time.Second, 10*time.Millisecond)
}

func isDeskShuffled(original, shuffled []uuid.UUID) bool {
	if len(original) != len(shuffled) {
		return false
	}
	sameOrder := true
	for i := range original {
		if original[i] != shuffled[i] {
			sameOrder = false
			break
		}
	}
	return !sameOrder
}
