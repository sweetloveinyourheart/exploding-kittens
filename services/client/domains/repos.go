package domains

import (
	"github.com/juju/pubsub/v2"

	eventing "github.com/sweetloveinyourheart/exploding-kittens/pkg/domain-eventing"
	"github.com/sweetloveinyourheart/exploding-kittens/pkg/domain-eventing/command_handler/bus"
	"github.com/sweetloveinyourheart/exploding-kittens/pkg/domains/desk"
	"github.com/sweetloveinyourheart/exploding-kittens/pkg/domains/game"
	"github.com/sweetloveinyourheart/exploding-kittens/pkg/domains/hand"
	"github.com/sweetloveinyourheart/exploding-kittens/pkg/domains/lobby"
)

// LobbyRepo is the repository for the Lobby aggregate.
var LobbyRepo eventing.ReadRepo[lobby.Lobby, *lobby.Lobby]

// GameRepo is the repository for the Game aggregate.
var GameRepo eventing.ReadRepo[game.Game, *game.Game]

// DeskRepo is the repository for the Desk aggregate.
var DeskRepo eventing.ReadRepo[desk.Desk, *desk.Desk]

// HandRepo is the repository for the Hand aggregate.
var HandRepo eventing.ReadRepo[hand.Hand, *hand.Hand]

var LobbySubscriber = pubsub.NewSimpleHub(&pubsub.SimpleHubConfig{})

var GameSubscriber = pubsub.NewSimpleHub(&pubsub.SimpleHubConfig{})

var CommandBus *bus.CommandHandler
