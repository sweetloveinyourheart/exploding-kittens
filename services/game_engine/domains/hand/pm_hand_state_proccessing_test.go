package hand_test

import (
	"context"
	"runtime"
	"time"

	"github.com/avast/retry-go"
	"github.com/gofrs/uuid"
	"github.com/sweetloveinyourheart/exploding-kittens/pkg/constants/cards"
	"github.com/sweetloveinyourheart/exploding-kittens/pkg/domain-eventing/command_handler/bus"
	retrymw "github.com/sweetloveinyourheart/exploding-kittens/pkg/domain-eventing/middleware/command_hander"
	"github.com/sweetloveinyourheart/exploding-kittens/pkg/domains/hand"
)

func (hs *HandSuite) TestHandStateProcessing_HandCreatedSuccessfully() {
	hs.setupEnvironment()
	_, _, cardsCodeMap := hs.prepareCards()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	mw := retrymw.NewCommandHandlerMiddleware(retry.Attempts(4), retry.MaxDelay(1*time.Second))
	commandBus := bus.NewCommandHandler()

	err := hand.AddNATSHandCommandHandlers(ctx, "test", commandBus, mw)
	hs.NoError(err)

	handRepo, err := hand.CreateNATSRepoHand(ctx, "test")
	hs.NoError(err)

	handID := uuid.Must(uuid.NewV7())
	err = commandBus.HandleCommand(ctx, &hand.CreateHand{
		HandID: handID,
		CardIDs: []uuid.UUID{
			uuid.Must(uuid.FromString(cardsCodeMap[cards.HairyPotatoCat].GetCardId())),
			uuid.Must(uuid.FromString(cardsCodeMap[cards.Attack].GetCardId())),
			uuid.Must(uuid.FromString(cardsCodeMap[cards.HairyPotatoCat].GetCardId())),
			uuid.Must(uuid.FromString(cardsCodeMap[cards.Nope].GetCardId())),
		},
	})
	hs.NoError(err)

	hs.Eventually(func() bool {
		runtime.Gosched()

		handState, handStateErr := handRepo.Find(ctx, handID.String())
		return handStateErr == nil && len(handState.GetCardIDs()) == 4
	}, 5*time.Second, 10*time.Millisecond)
}

func (hs *HandSuite) TestHandStateProcessing_CardsGivenCorrectly() {
	hs.setupEnvironment()
	_, _, cardsCodeMap := hs.prepareCards()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	mw := retrymw.NewCommandHandlerMiddleware(retry.Attempts(4), retry.MaxDelay(1*time.Second))
	commandBus := bus.NewCommandHandler()

	err := hand.AddNATSHandCommandHandlers(ctx, "test", commandBus, mw)
	hs.NoError(err)

	handRepo, err := hand.CreateNATSRepoHand(ctx, "test")
	hs.NoError(err)

	handID_01 := uuid.Must(uuid.NewV7())
	err = commandBus.HandleCommand(ctx, &hand.CreateHand{
		HandID: handID_01,
		CardIDs: []uuid.UUID{
			uuid.Must(uuid.FromString(cardsCodeMap[cards.HairyPotatoCat].GetCardId())),
			uuid.Must(uuid.FromString(cardsCodeMap[cards.Attack].GetCardId())),
			uuid.Must(uuid.FromString(cardsCodeMap[cards.HairyPotatoCat].GetCardId())),
			uuid.Must(uuid.FromString(cardsCodeMap[cards.Nope].GetCardId())),
		},
	})
	hs.NoError(err)

	handID_02 := uuid.Must(uuid.NewV7())
	err = commandBus.HandleCommand(ctx, &hand.CreateHand{
		HandID: handID_02,
		CardIDs: []uuid.UUID{
			uuid.Must(uuid.FromString(cardsCodeMap[cards.HairyPotatoCat].GetCardId())),
			uuid.Must(uuid.FromString(cardsCodeMap[cards.Attack].GetCardId())),
			uuid.Must(uuid.FromString(cardsCodeMap[cards.HairyPotatoCat].GetCardId())),
			uuid.Must(uuid.FromString(cardsCodeMap[cards.Nope].GetCardId())),
		},
	})
	hs.NoError(err)

	hs.Eventually(func() bool {
		runtime.Gosched()

		handState_01, handStateErr_01 := handRepo.Find(ctx, handID_01.String())
		handState_02, handStateErr_02 := handRepo.Find(ctx, handID_02.String())

		return handStateErr_01 == nil && len(handState_01.GetCardIDs()) == 4 &&
			handStateErr_02 == nil && len(handState_02.GetCardIDs()) == 4
	}, 5*time.Second, 10*time.Millisecond)

	err = commandBus.HandleCommand(ctx, &hand.GiveCards{
		HandID:   handID_01,
		ToHandID: handID_02,
		CardIDs: []uuid.UUID{
			uuid.Must(uuid.FromString(cardsCodeMap[cards.Attack].GetCardId())),
		},
	})
	hs.NoError(err)

	hs.Eventually(func() bool {
		runtime.Gosched()

		handState_01, handStateErr_01 := handRepo.Find(ctx, handID_01.String())
		handState_02, handStateErr_02 := handRepo.Find(ctx, handID_02.String())

		return handStateErr_01 == nil && len(handState_01.GetCardIDs()) == 3 &&
			handStateErr_02 == nil && len(handState_02.GetCardIDs()) == 5
	}, 5*time.Second, 10*time.Millisecond)
}

func (hs *HandSuite) TestHandStateProcessing_PlayCard_Single_Success() {
	hs.setupEnvironment()
	_, _, cardsCodeMap := hs.prepareCards()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	mw := retrymw.NewCommandHandlerMiddleware(retry.Attempts(4), retry.MaxDelay(1*time.Second))
	commandBus := bus.NewCommandHandler()

	err := hand.AddNATSHandCommandHandlers(ctx, "test", commandBus, mw)
	hs.NoError(err)

	handRepo, err := hand.CreateNATSRepoHand(ctx, "test")
	hs.NoError(err)

	handID := uuid.Must(uuid.NewV7())
	err = commandBus.HandleCommand(ctx, &hand.CreateHand{
		HandID: handID,
		CardIDs: []uuid.UUID{
			uuid.Must(uuid.FromString(cardsCodeMap[cards.HairyPotatoCat].GetCardId())),
			uuid.Must(uuid.FromString(cardsCodeMap[cards.Attack].GetCardId())),
			uuid.Must(uuid.FromString(cardsCodeMap[cards.Attack].GetCardId())),
			uuid.Must(uuid.FromString(cardsCodeMap[cards.HairyPotatoCat].GetCardId())),
			uuid.Must(uuid.FromString(cardsCodeMap[cards.Nope].GetCardId())),
		},
	})
	hs.NoError(err)

	hs.Eventually(func() bool {
		runtime.Gosched()

		handState, handStateErr := handRepo.Find(ctx, handID.String())
		return handStateErr == nil && len(handState.GetCardIDs()) == 5
	}, 5*time.Second, 10*time.Millisecond)

	err = commandBus.HandleCommand(ctx, &hand.PlayCards{
		HandID: handID,
		CardIDs: []uuid.UUID{
			uuid.Must(uuid.FromString(cardsCodeMap[cards.Attack].GetCardId())),
		},
	})
	hs.NoError(err)

	hs.Eventually(func() bool {
		runtime.Gosched()

		handState, handStateErr := handRepo.Find(ctx, handID.String())
		return handStateErr == nil && len(handState.GetCardIDs()) == 4
	}, 5*time.Second, 10*time.Millisecond)
}

func (hs *HandSuite) TestHandStateProcessing_PlayCard_Twice_Success() {
	hs.setupEnvironment()
	_, _, cardsCodeMap := hs.prepareCards()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	mw := retrymw.NewCommandHandlerMiddleware(retry.Attempts(4), retry.MaxDelay(1*time.Second))
	commandBus := bus.NewCommandHandler()

	err := hand.AddNATSHandCommandHandlers(ctx, "test", commandBus, mw)
	hs.NoError(err)

	handRepo, err := hand.CreateNATSRepoHand(ctx, "test")
	hs.NoError(err)

	handID := uuid.Must(uuid.NewV7())
	err = commandBus.HandleCommand(ctx, &hand.CreateHand{
		HandID: handID,
		CardIDs: []uuid.UUID{
			uuid.Must(uuid.FromString(cardsCodeMap[cards.HairyPotatoCat].GetCardId())),
			uuid.Must(uuid.FromString(cardsCodeMap[cards.Attack].GetCardId())),
			uuid.Must(uuid.FromString(cardsCodeMap[cards.Attack].GetCardId())),
			uuid.Must(uuid.FromString(cardsCodeMap[cards.HairyPotatoCat].GetCardId())),
			uuid.Must(uuid.FromString(cardsCodeMap[cards.Nope].GetCardId())),
		},
	})
	hs.NoError(err)

	hs.Eventually(func() bool {
		runtime.Gosched()

		handState, handStateErr := handRepo.Find(ctx, handID.String())
		return handStateErr == nil && len(handState.GetCardIDs()) == 5
	}, 5*time.Second, 10*time.Millisecond)

	err = commandBus.HandleCommand(ctx, &hand.PlayCards{
		HandID: handID,
		CardIDs: []uuid.UUID{
			uuid.Must(uuid.FromString(cardsCodeMap[cards.Attack].GetCardId())),
		},
	})
	hs.NoError(err)

	hs.Eventually(func() bool {
		runtime.Gosched()

		handState, handStateErr := handRepo.Find(ctx, handID.String())
		return handStateErr == nil && len(handState.GetCardIDs()) == 4
	}, 5*time.Second, 10*time.Millisecond)

	err = commandBus.HandleCommand(ctx, &hand.PlayCards{
		HandID: handID,
		CardIDs: []uuid.UUID{
			uuid.Must(uuid.FromString(cardsCodeMap[cards.Attack].GetCardId())),
		},
	})
	hs.NoError(err)

	hs.Eventually(func() bool {
		runtime.Gosched()

		handState, handStateErr := handRepo.Find(ctx, handID.String())
		return handStateErr == nil && len(handState.GetCardIDs()) == 3
	}, 5*time.Second, 10*time.Millisecond)
}
