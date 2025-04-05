package gameengine

import (
	"context"
	"embed"
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
	gameDomain "github.com/sweetloveinyourheart/exploding-kittens/services/game_engine/domains/game"
)

//go:embed migrations/*.sql
var FS embed.FS

func InitializeRepos(ctx context.Context) error {
	appID := fmt.Sprintf("gameserver-%s", config.Instance().GetString("gameserver.id"))

	if err := InitializeCoreRepos(appID, ctx); err != nil {
		log.Global().ErrorContext(ctx, "failed to initialize core repos", zap.Error(err))
		return err
	}

	_, err := gameDomain.NewGameInteractionProcessor(ctx)
	if err != nil {
		return err
	}

	return nil
}

func InitializeCoreRepos(appID string, ctx context.Context) error {
	var err error

	domains.GameRepo, err = game.CreateNATSRepoGames(ctx, appID)
	if err != nil {
		return err
	}

	domains.DeskRepo, err = desk.CreateNATSRepoDesk(ctx, appID)
	if err != nil {
		return err
	}

	domains.HandRepo, err = hand.CreateNATSRepoHand(ctx, appID)
	if err != nil {
		return err
	}

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
