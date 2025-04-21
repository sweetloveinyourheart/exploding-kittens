package client

import (
	"context"
	"fmt"

	"github.com/sweetloveinyourheart/exploding-kittens/pkg/config"
	eventing "github.com/sweetloveinyourheart/exploding-kittens/pkg/domain-eventing"
	"github.com/sweetloveinyourheart/exploding-kittens/pkg/domain-eventing/command_handler/bus"
	"github.com/sweetloveinyourheart/exploding-kittens/pkg/domain-eventing/common"
	consumerinvalidator "github.com/sweetloveinyourheart/exploding-kittens/pkg/domain-eventing/middleware/event_handler/consumer_invalidator"
	"github.com/sweetloveinyourheart/exploding-kittens/pkg/domains/lobby"
	"github.com/sweetloveinyourheart/exploding-kittens/services/client/domains"
	"github.com/sweetloveinyourheart/exploding-kittens/services/client/domains/match"
)

func InitializeRepos(ctx context.Context) error {
	appID := fmt.Sprintf("clientserver-%s", config.Instance().GetString("clientserver.id"))

	domains.CommandBus = bus.NewCommandHandler()
	err := lobby.AddNATSLobbyCommandHandlers(ctx, appID, domains.CommandBus)
	if err != nil {
		return err
	}

	allEvents := []common.EventType{}
	allEvents = append(allEvents, lobby.AllEventTypes...)
	mw := consumerinvalidator.NewMiddleware(eventing.MatchEvents(allEvents), func(ctx context.Context, event common.Event) {
		switch event.EventType() {
		case lobby.EventTypeLobbyCreated:
			if lobby, ok := event.Data().(match.LobbyIDer); ok {
				domains.LobbySubscriber.Publish(lobby.GetLobbyId().String(), struct{}{})
			}
		}

		if event.AggregateType() == lobby.AggregateType {
			if lobby, ok := event.Data().(match.LobbyIDer); ok {
				domains.LobbySubscriber.Publish(lobby.GetLobbyId().String(), struct{}{})
			}
		}
	})

	domains.LobbyRepo, err = lobby.CreateNATSRepoLobbies(ctx, appID, mw)
	if err != nil {
		return err
	}

	return nil
}
