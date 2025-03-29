package lobby

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
	"github.com/sweetloveinyourheart/exploding-kittens/pkg/domains/lobby"
	log "github.com/sweetloveinyourheart/exploding-kittens/pkg/logger"
	"github.com/sweetloveinyourheart/exploding-kittens/services/lobby/domains"
	lobbyDomain "github.com/sweetloveinyourheart/exploding-kittens/services/lobby/domains/lobby"
)

func InitializeRepos(ctx context.Context) error {
	appID := fmt.Sprintf("lobbyserver-%s", config.Instance().GetString("lobbyserver.id"))

	if err := InitializeCoreRepos(appID, ctx); err != nil {
		log.Global().ErrorContext(ctx, "failed to initialize core repos", zap.Error(err))
		return err
	}

	_, err := lobbyDomain.NewLobbyInteractionProcessor(ctx)
	if err != nil {
		return err
	}

	return nil
}

func InitializeCoreRepos(appID string, ctx context.Context) error {
	var err error

	domains.LobbyRepo, err = lobby.CreateNATSRepoLobbies(ctx, appID)
	if err != nil {
		return err
	}

	domains.CommandBus = bus.NewCommandHandler()
	err = lobby.AddNATSLobbyCommandHandlers(ctx, appID, domains.CommandBus)
	if err != nil {
		return err
	}

	connPool, err := do.InvokeNamed[*pool.ConnPool](nil, string(constants.ConnectionPool))
	if err != nil {
		return err
	}

	if err := initLobbyEventBus(appID, connPool); err != nil {
		return err
	}

	return nil
}

func initLobbyEventBus(appID string, connPool *pool.ConnPool) error {
	neb, err := nats.NewEventBus(connPool, fmt.Sprintf("%s-bus", appID), nats.WithStreamName(constants.LobbyStream))
	if err != nil {
		return err
	}

	domains.LobbyEventBus = neb

	return nil
}
