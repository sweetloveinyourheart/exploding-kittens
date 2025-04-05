package domains

import (
	eventing "github.com/sweetloveinyourheart/exploding-kittens/pkg/domain-eventing"
	"github.com/sweetloveinyourheart/exploding-kittens/pkg/domain-eventing/command_handler/bus"
	"github.com/sweetloveinyourheart/exploding-kittens/pkg/domains/desk"
	"github.com/sweetloveinyourheart/exploding-kittens/pkg/domains/game"
	"github.com/sweetloveinyourheart/exploding-kittens/pkg/domains/hand"
)

// CommandBus is the command bus command processing.
var CommandBus *bus.CommandHandler

// GameEventBus is the event bus for consuming events in the game domain.
var GameEventBus eventing.EventBus

// DeskEventBus is the event bus for consuming events in the desk domain.
var DeskEventBus eventing.EventBus

// HandEventBus is the event bus for consuming events in the desk domain.
var HandEventBus eventing.EventBus

// GameRepo is the repository for the game aggregate.
var GameRepo eventing.ReadRepo[game.Game, *game.Game]

// DeskRepo is the repository for the desk aggregate.
var DeskRepo eventing.ReadRepo[desk.Desk, *desk.Desk]

// HandRepo is the repository for the hand aggregate.
var HandRepo eventing.ReadRepo[hand.Hand, *hand.Hand]
