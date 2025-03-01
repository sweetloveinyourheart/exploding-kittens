package common

import "github.com/sweetloveinyourheart/exploding-kittens/pkg/stringsutil"

// CommandType is the type of a command, used as its unique identifier.
type CommandType string

// String returns the string representation of a command type.
func (ct CommandType) String() string {
	return string(ct)
}

// CommandFieldError is returned by Dispatch when a field is incorrect.
type CommandFieldError struct {
	Field   string
	Details string
}

// Error implements the Error method of the error interface.
func (c *CommandFieldError) Error() string {
	if stringsutil.IsBlank(c.Details) {
		return "missing field: " + c.Field
	} else {
		return "invalid field: " + c.Field + ", " + c.Details
	}
}
