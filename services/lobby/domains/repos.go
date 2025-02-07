package domains

import (
	eventing "github.com/sweetloveinyourheart/exploding-kittens/pkg/domain-eventing"
	"github.com/sweetloveinyourheart/exploding-kittens/pkg/domain-eventing/command_handler/bus"
	"github.com/sweetloveinyourheart/exploding-kittens/pkg/domains/lobby"
)

// CommandBus is the command bus command processing.
var CommandBus *bus.CommandHandler

// LobbyEventBus is the event bus for consuming events in the lobby domain.
var LobbyEventBus eventing.EventBus

// LobbyRepo is the repository for the Lobby aggregate.
var LobbyRepo eventing.ReadRepo[lobby.Lobby, *lobby.Lobby]
