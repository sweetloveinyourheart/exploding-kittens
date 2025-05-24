package gameengine

import (
	"context"
	"fmt"

	pool "github.com/octu0/nats-pool"
	"github.com/samber/do"
	"go.uber.org/zap"

	"github.com/sweetloveinyourheart/exploding-kittens/pkg/config"
	"github.com/sweetloveinyourheart/exploding-kittens/pkg/constants"
	"github.com/sweetloveinyourheart/exploding-kittens/pkg/domain-eventing/command_handler/bus"
	"github.com/sweetloveinyourheart/exploding-kittens/pkg/domain-eventing/event_bus/nats"
	"github.com/sweetloveinyourheart/exploding-kittens/pkg/domains/desk"
	"github.com/sweetloveinyourheart/exploding-kittens/pkg/domains/game"
	"github.com/sweetloveinyourheart/exploding-kittens/pkg/domains/hand"
	log "github.com/sweetloveinyourheart/exploding-kittens/pkg/logger"
	"github.com/sweetloveinyourheart/exploding-kittens/services/game_engine/domains"
	deskDomain "github.com/sweetloveinyourheart/exploding-kittens/services/game_engine/domains/desk"
	gameDomain "github.com/sweetloveinyourheart/exploding-kittens/services/game_engine/domains/game"
	handDomain "github.com/sweetloveinyourheart/exploding-kittens/services/game_engine/domains/hand"
)

func InitializeRepos(ctx context.Context) error {
	appID := fmt.Sprintf("gameengineserver-%s", config.Instance().GetString("gameengineserver.id"))

	if err := InitializeCoreRepos(appID, ctx); err != nil {
		log.Global().ErrorContext(ctx, "failed to initialize core repos", zap.Error(err))
		return err
	}

	_, err := gameDomain.NewGameInteractionProcessor(ctx)
	if err != nil {
		return err
	}

	_, err = handDomain.NewHandStateProcessor(ctx)
	if err != nil {
		return err
	}

	_, err = deskDomain.NewDeskStateProcessor(ctx)
	if err != nil {
		return err
	}

	_, err = gameDomain.NewGamePlayExecutor(ctx)
	if err != nil {
		return err
	}

	_, err = gameDomain.NewGameExplodingProcessManager(ctx)
	if err != nil {
		return err
	}

	return nil
}

func InitializeCoreRepos(appID string, ctx context.Context) error {
	var err error

	domains.CommandBus = bus.NewCommandHandler()

	err = game.AddNATSGameCommandHandlers(ctx, appID, domains.CommandBus)
	if err != nil {
		return err
	}

	err = desk.AddNATSDeskCommandHandlers(ctx, appID, domains.CommandBus)
	if err != nil {
		return err
	}

	err = hand.AddNATSHandCommandHandlers(ctx, appID, domains.CommandBus)
	if err != nil {
		return err
	}

	connPool, err := do.InvokeNamed[*pool.ConnPool](nil, string(constants.ConnectionPool))
	if err != nil {
		return err
	}

	if err := initGameEventBus(appID, connPool); err != nil {
		return err
	}

	return nil
}

func initGameEventBus(appID string, connPool *pool.ConnPool) error {
	neb, err := nats.NewEventBus(connPool, fmt.Sprintf("%s-bus", appID), nats.WithStreamName(constants.GameStream))
	if err != nil {
		return err
	}

	domains.GameEventBus = neb

	return nil
}
