package domains

import (
	"github.com/juju/pubsub/v2"

	eventing "github.com/sweetloveinyourheart/exploding-kittens/pkg/domain-eventing"
	"github.com/sweetloveinyourheart/exploding-kittens/pkg/domain-eventing/command_handler/bus"
	"github.com/sweetloveinyourheart/exploding-kittens/pkg/domains/lobby"
)

// LobbyRepo is the repository for the Lobby aggregate.
var LobbyRepo eventing.ReadRepo[lobby.Lobby, *lobby.Lobby]

var LobbySubscriber = pubsub.NewSimpleHub(&pubsub.SimpleHubConfig{})

var CommandBus *bus.CommandHandler
