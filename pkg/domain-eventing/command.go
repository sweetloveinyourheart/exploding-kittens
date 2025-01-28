package eventing

import "github.com/sweetloveinyourheart/exploding-kittens/pkg/domain-eventing/common"

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
