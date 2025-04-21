package domains

import (
	eventing "github.com/sweetloveinyourheart/exploding-kittens/pkg/domain-eventing"
	"github.com/sweetloveinyourheart/exploding-kittens/pkg/domain-eventing/command_handler/bus"
)

// CommandBus is the command bus command processing.
var CommandBus *bus.CommandHandler

// GameEventBus is the event bus for consuming events in the game domain.
var GameEventBus eventing.EventBus

// DeskEventBus is the event bus for consuming events in the desk domain.
var DeskEventBus eventing.EventBus

// HandEventBus is the event bus for consuming events in the desk domain.
var HandEventBus eventing.EventBus
