package eventing

import (
	"fmt"
	"sync"

	"github.com/sweetloveinyourheart/exploding-kittens/pkg/domain-eventing/common"
)

// Command is a domain command that is sent to a Dispatcher.
//
// A command name should 1) be in present tense and 2) contain the intent
// (MoveCustomer vs CorrectCustomerAddress).
//
// The command should contain all the data needed when handling it as fields.
// These fields can take an optional "eventing" tag, which adds properties. For now
// only "optional" is a valid tag: `eventing:"optional"`.
type Command interface {
	// AggregateID returns the ID of the aggregate that the command should be
	// handled by.
	AggregateID() string

	// AggregateType returns the type of the aggregate that the command can be
	// handled by.
	AggregateType() common.AggregateType

	// CommandType returns the type of the command.
	CommandType() common.CommandType

	Validate() error
}

type GenericCommand[T any] interface {
	Command
	*T
}

// RegisterCommand registers an command factory for a type. The factory is
// used to create concrete command types.
//
// An example would be:
//
//	RegisterCommand(func() Command { return &MyCommand{} })
func RegisterCommand[T any, PT GenericCommand[T]]() {
	var cmd PT = new(T)
	if cmd == nil {
		panic("eventing: created command is nil")
	}

	commandType := cmd.CommandType()
	if commandType == common.CommandType("") {
		panic("eventing: attempt to register empty command type")
	}

	commandsMu.Lock()
	defer commandsMu.Unlock()

	if _, ok := commands[commandType]; ok {
		panic(fmt.Sprintf("eventing: registering duplicate types for %q", commandType))
	}

	commands[commandType] = func() Command { return PT(new(T)) }
}

var commands = make(map[common.CommandType]func() Command)
var commandsMu sync.RWMutex
