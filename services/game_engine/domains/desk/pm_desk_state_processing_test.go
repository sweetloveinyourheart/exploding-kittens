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
)

func (hs *DeskSuite) TestDeskStateProcessing_DeskCreatedSuccessfully() {
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

func (hs *DeskSuite) TestDeskStateProcessing_DeskShuffledSuccessfully() {
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
