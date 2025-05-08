package game_test

import (
	"context"
	"fmt"
	"time"

	"github.com/avast/retry-go"
	"github.com/gofrs/uuid"
	"github.com/nats-io/nats.go"
	"github.com/samber/do"

	"github.com/sweetloveinyourheart/exploding-kittens/pkg/constants"
	card_effects "github.com/sweetloveinyourheart/exploding-kittens/pkg/constants/card-effects"
	"github.com/sweetloveinyourheart/exploding-kittens/pkg/constants/cards"
	codecJson "github.com/sweetloveinyourheart/exploding-kittens/pkg/domain-eventing/codec/json"
	"github.com/sweetloveinyourheart/exploding-kittens/pkg/domain-eventing/command_handler/bus"
	natsPkg "github.com/sweetloveinyourheart/exploding-kittens/pkg/domain-eventing/event_bus/nats"
	retrymw "github.com/sweetloveinyourheart/exploding-kittens/pkg/domain-eventing/middleware/command_hander"
	deskDomain "github.com/sweetloveinyourheart/exploding-kittens/pkg/domains/desk"
	gameDomain "github.com/sweetloveinyourheart/exploding-kittens/pkg/domains/game"
	handDomain "github.com/sweetloveinyourheart/exploding-kittens/pkg/domains/hand"
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

	gs.Eventually(func() bool {
		gameState, gameStateErr := gameRepo.Find(ctx, gameID.String())
		return gameStateErr == nil && gameState.GetPlayerTurn() == player01
	}, 5*time.Second, 10*time.Millisecond)

	err = commandBus.HandleCommand(ctx, &gameDomain.PlayCard{
		GameID:   gameID,
		PlayerID: player01,
		CardIDs:  []uuid.UUID{uuid.Must(uuid.FromString(cardsCodeMap[cards.Skip].GetCardId()))},
	})
	gs.NoError(err)

	gs.Eventually(func() bool {
		gameState, gameStateErr := gameRepo.Find(ctx, gameID.String())
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

	var previousDeskState *deskDomain.Desk
	gs.Eventually(func() bool {
		gameState, gameStateErr := gameRepo.Find(ctx, gameID.String())
		deskState, deskStateErr := deskRepo.Find(ctx, gameState.Desk.String())
		previousDeskState = deskState
		return gameStateErr == nil && deskStateErr == nil && gameState.Desk != uuid.Nil
	}, 5*time.Second, 10*time.Millisecond)

	err = commandBus.HandleCommand(ctx, &gameDomain.PlayCard{
		GameID:   gameID,
		PlayerID: player01,
		CardIDs:  []uuid.UUID{uuid.Must(uuid.FromString(cardsCodeMap[cards.Shuffle].GetCardId()))},
	})
	gs.NoError(err)

	var afterShuffleDeskState *deskDomain.Desk
	gs.Eventually(func() bool {
		gameState, gameStateErr := gameRepo.Find(ctx, gameID.String())
		deskState, deskStateErr := deskRepo.Find(ctx, gameState.Desk.String())
		afterShuffleDeskState = deskState
		return gameStateErr == nil && deskStateErr == nil && gameState.Desk != uuid.Nil && gameState.GetPlayerTurn() == player01
	}, 5*time.Second, 10*time.Millisecond)

	gs.Equal(len(previousDeskState.CardIDs), len(afterShuffleDeskState.CardIDs))
	gs.NotEqual(previousDeskState.CardIDs, afterShuffleDeskState.CardIDs)
}

func (gs *GameSuite) TestGamePlayExecutor_HandleCardPlay_Favor() {
	gs.setupEnvironment()
	_, _, cardsCodeMap := gs.prepareCards()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	mw := retrymw.NewCommandHandlerMiddleware(retry.Attempts(4), retry.MaxDelay(1*time.Second))
	commandBus := bus.NewCommandHandler()

	err := gameDomain.AddNATSGameCommandHandlers(ctx, "test", commandBus, mw)
	gs.NoError(err)

	err = handDomain.AddNATSHandCommandHandlers(ctx, "test", commandBus, mw)
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
		gameState, gameStateErr := gameRepo.Find(ctx, gameID.String())
		return gameStateErr == nil && gameState.GetPlayerTurn() == player01
	}, 5*time.Second, 10*time.Millisecond)

	err = commandBus.HandleCommand(ctx, &gameDomain.PlayCard{
		GameID:   gameID,
		PlayerID: player01,
		CardIDs:  []uuid.UUID{uuid.Must(uuid.FromString(cardsCodeMap[cards.Favor].GetCardId()))},
	})
	gs.NoError(err)

	gs.Eventually(func() bool {
		gameState, gameStateErr := gameRepo.Find(ctx, gameID.String())

		isGameStateValid := gameStateErr == nil &&
			gameState.GetPlayerTurn() == player01 &&
			gameState.GetExecutingAction() == card_effects.StealCard &&
			gameState.GetGamePhase() == gameDomain.GAME_PHASE_ACTION_PHASE

		return isGameStateValid
	}, 5*time.Second, 10*time.Millisecond)

	err = commandBus.HandleCommand(ctx, &gameDomain.ExecuteAction{
		GameID:         gameID,
		PlayerID:       player01,
		Effect:         card_effects.StealCard,
		TargetPlayerID: player02,
		TargetCardID:   uuid.Must(uuid.FromString(cardsCodeMap[cards.Favor].GetCardId())),
	})
	gs.NoError(err)

	bus, err := do.InvokeNamed[nats.JetStreamContext](nil, string(constants.Bus))
	gs.NoError(err)
	_ = bus

	gs.Eventually(func() bool {
		events, err := natsPkg.LoadJetStreamCtx(ctx, bus, constants.HandStream, fmt.Sprintf("%s.>", constants.HandStream), &codecJson.EventCodec{})
		if err != nil {
			return false
		}

		for _, event := range events {
			if event.EventType() == handDomain.EventTypeCardStolen {
				return true
			}
		}

		return false
	}, 5*time.Second, 100*time.Millisecond)
}

func (gs *GameSuite) TestGamePlayExecutor_HandleCardPlay_Combo2_HairyPotatoCat() {
	gs.setupEnvironment()
	_, _, cardsCodeMap := gs.prepareCards()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	mw := retrymw.NewCommandHandlerMiddleware(retry.Attempts(4), retry.MaxDelay(1*time.Second))
	commandBus := bus.NewCommandHandler()

	err := gameDomain.AddNATSGameCommandHandlers(ctx, "test", commandBus, mw)
	gs.NoError(err)

	err = handDomain.AddNATSHandCommandHandlers(ctx, "test", commandBus, mw)
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
		gameState, gameStateErr := gameRepo.Find(ctx, gameID.String())
		return gameStateErr == nil && gameState.GetPlayerTurn() == player01
	}, 5*time.Second, 10*time.Millisecond)

	err = commandBus.HandleCommand(ctx, &gameDomain.PlayCard{
		GameID:   gameID,
		PlayerID: player01,
		CardIDs: []uuid.UUID{
			uuid.Must(uuid.FromString(cardsCodeMap[cards.HairyPotatoCat].GetCardId())),
			uuid.Must(uuid.FromString(cardsCodeMap[cards.HairyPotatoCat].GetCardId())),
		},
	})
	gs.NoError(err)

	gs.Eventually(func() bool {
		gameState, gameStateErr := gameRepo.Find(ctx, gameID.String())

		isGameStateValid := gameStateErr == nil &&
			gameState.GetPlayerTurn() == player01 &&
			gameState.GetExecutingAction() == card_effects.StealRandomCard &&
			gameState.GetGamePhase() == gameDomain.GAME_PHASE_ACTION_PHASE

		return isGameStateValid
	}, 5*time.Second, 10*time.Millisecond)

	err = commandBus.HandleCommand(ctx, &gameDomain.ExecuteAction{
		GameID:         gameID,
		PlayerID:       player01,
		Effect:         card_effects.StealRandomCard,
		TargetPlayerID: player02,
		TargetCardID:   uuid.Must(uuid.FromString(cardsCodeMap[cards.Defuse].GetCardId())),
	})
	gs.NoError(err)

	bus, err := do.InvokeNamed[nats.JetStreamContext](nil, string(constants.Bus))
	gs.NoError(err)
	_ = bus

	gs.Eventually(func() bool {
		events, err := natsPkg.LoadJetStreamCtx(ctx, bus, constants.HandStream, fmt.Sprintf("%s.>", constants.HandStream), &codecJson.EventCodec{})
		if err != nil {
			return false
		}

		for _, event := range events {
			if event.EventType() == handDomain.EventTypeCardStolen {
				return true
			}
		}

		return false
	}, 5*time.Second, 100*time.Millisecond)
}

func (gs *GameSuite) TestGamePlayExecutor_HandleCardPlay_Combo3_BreadCat() {
	gs.setupEnvironment()
	_, _, cardsCodeMap := gs.prepareCards()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	mw := retrymw.NewCommandHandlerMiddleware(retry.Attempts(4), retry.MaxDelay(1*time.Second))
	commandBus := bus.NewCommandHandler()

	err := gameDomain.AddNATSGameCommandHandlers(ctx, "test", commandBus, mw)
	gs.NoError(err)

	err = handDomain.AddNATSHandCommandHandlers(ctx, "test", commandBus, mw)
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
		gameState, gameStateErr := gameRepo.Find(ctx, gameID.String())
		return gameStateErr == nil && gameState.GetPlayerTurn() == player01
	}, 5*time.Second, 10*time.Millisecond)

	err = commandBus.HandleCommand(ctx, &gameDomain.PlayCard{
		GameID:   gameID,
		PlayerID: player01,
		CardIDs: []uuid.UUID{
			uuid.Must(uuid.FromString(cardsCodeMap[cards.BeardCat].GetCardId())),
			uuid.Must(uuid.FromString(cardsCodeMap[cards.BeardCat].GetCardId())),
			uuid.Must(uuid.FromString(cardsCodeMap[cards.BeardCat].GetCardId())),
		},
	})
	gs.NoError(err)

	gs.Eventually(func() bool {
		gameState, gameStateErr := gameRepo.Find(ctx, gameID.String())

		isGameStateValid := gameStateErr == nil &&
			gameState.GetPlayerTurn() == player01 &&
			gameState.GetExecutingAction() == card_effects.StealNamedCard &&
			gameState.GetGamePhase() == gameDomain.GAME_PHASE_ACTION_PHASE

		return isGameStateValid
	}, 5*time.Second, 10*time.Millisecond)

	err = commandBus.HandleCommand(ctx, &gameDomain.ExecuteAction{
		GameID:         gameID,
		PlayerID:       player01,
		Effect:         card_effects.StealNamedCard,
		TargetPlayerID: player02,
		TargetCardID:   uuid.Must(uuid.FromString(cardsCodeMap[cards.Defuse].GetCardId())),
	})
	gs.NoError(err)

	bus, err := do.InvokeNamed[nats.JetStreamContext](nil, string(constants.Bus))
	gs.NoError(err)
	_ = bus

	gs.Eventually(func() bool {
		events, err := natsPkg.LoadJetStreamCtx(ctx, bus, constants.HandStream, fmt.Sprintf("%s.>", constants.HandStream), &codecJson.EventCodec{})
		if err != nil {
			return false
		}

		for _, event := range events {
			if event.EventType() == handDomain.EventTypeCardStolen {
				return true
			}
		}

		return false
	}, 5*time.Second, 100*time.Millisecond)
}

func (gs *GameSuite) TestGamePlayExecutor_HandleCardPlay_SeeTheFuture() {
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
		gameState, gameStateErr := gameRepo.Find(ctx, gameID.String())
		return gameStateErr == nil && gameState.GetPlayerTurn() == player01
	}, 5*time.Second, 10*time.Millisecond)

	err = commandBus.HandleCommand(ctx, &gameDomain.PlayCard{
		GameID:   gameID,
		PlayerID: player01,
		CardIDs:  []uuid.UUID{uuid.Must(uuid.FromString(cardsCodeMap[cards.SeeTheFuture].GetCardId()))},
	})
	gs.NoError(err)

	gs.Eventually(func() bool {
		gameState, gameStateErr := gameRepo.Find(ctx, gameID.String())

		isGameStateValid := gameStateErr == nil &&
			gameState.GetPlayerTurn() == player01 &&
			gameState.GetExecutingAction() == card_effects.PeekCards &&
			gameState.GetGamePhase() == gameDomain.GAME_PHASE_ACTION_PHASE

		return isGameStateValid
	}, 5*time.Second, 10*time.Millisecond)
}

func (gs *GameSuite) TestGamePlayExecutor_HandleCardPlay_Nope() {
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
		gameState, gameStateErr := gameRepo.Find(ctx, gameID.String())
		return gameStateErr == nil && gameState.GetPlayerTurn() == player01
	}, 5*time.Second, 10*time.Millisecond)

	err = commandBus.HandleCommand(ctx, &gameDomain.PlayCard{
		GameID:   gameID,
		PlayerID: player01,
		CardIDs:  []uuid.UUID{uuid.Must(uuid.FromString(cardsCodeMap[cards.Attack].GetCardId()))},
	})
	gs.NoError(err)

	gs.Eventually(func() bool {
		gameState, gameStateErr := gameRepo.Find(ctx, gameID.String())

		isGameStateValid := gameStateErr == nil &&
			gameState.GetPlayerTurn() == player02 &&
			gameState.GetExecutingAction() == "" &&
			gameState.GetGamePhase() == gameDomain.GAME_PHASE_TURN_START

		return isGameStateValid
	}, 5*time.Second, 10*time.Millisecond)

	err = commandBus.HandleCommand(ctx, &gameDomain.PlayCard{
		GameID:   gameID,
		PlayerID: player02,
		CardIDs:  []uuid.UUID{uuid.Must(uuid.FromString(cardsCodeMap[cards.Nope].GetCardId()))},
	})
	gs.NoError(err)

	gs.Eventually(func() bool {
		gameState, gameStateErr := gameRepo.Find(ctx, gameID.String())

		isGameStateValid := gameStateErr == nil &&
			gameState.GetPlayerTurn() == player01 &&
			gameState.GetExecutingAction() == "" &&
			gameState.GetGamePhase() == gameDomain.GAME_PHASE_TURN_START

		return isGameStateValid
	}, 5*time.Second, 10*time.Millisecond)
}
