package eventing

import (
	"github.com/cockroachdb/errors"

	"github.com/gofrs/uuid"
)

var (
	// ErrMissingCommand is when there is no command to be handled.
	ErrMissingCommand = errors.New("missing command")
	// ErrMissingAggregateID is when a command is missing an aggregate ID.
	ErrMissingAggregateID = errors.New("missing aggregate ID")
)

// CheckCommand checks a command for errors.
func CheckCommand(cmd Command) error {
	if cmd == nil {
		return ErrMissingCommand
	}

	if cmd.AggregateID() == "" || cmd.AggregateID() == uuid.Nil.String() {
		return ErrMissingAggregateID
	}

	return cmd.Validate()
}
