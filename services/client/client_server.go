package client

import (
	"context"
	"fmt"

	"github.com/sweetloveinyourheart/exploding-kittens/pkg/config"
	"github.com/sweetloveinyourheart/exploding-kittens/pkg/domain-eventing/command_handler/bus"
	"github.com/sweetloveinyourheart/exploding-kittens/pkg/domains/desk"
	"github.com/sweetloveinyourheart/exploding-kittens/pkg/domains/game"
	"github.com/sweetloveinyourheart/exploding-kittens/pkg/domains/hand"
	"github.com/sweetloveinyourheart/exploding-kittens/pkg/domains/lobby"
	"github.com/sweetloveinyourheart/exploding-kittens/services/client/domains"
)

func InitializeRepos(ctx context.Context) error {
	appID := fmt.Sprintf("clientserver-%s", config.Instance().GetString("clientserver.id"))

	domains.CommandBus = bus.NewCommandHandler()
	err := lobby.AddNATSLobbyCommandHandlers(ctx, appID, domains.CommandBus)
	if err != nil {
		return err
	}

	err = game.AddNATSGameCommandHandlers(ctx, appID, domains.CommandBus)
	if err != nil {
		return err
	}

	domains.LobbyRepo, err = lobby.CreateNATSRepoLobbies(ctx, appID)
	if err != nil {
		return err
	}

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

	return nil
}
