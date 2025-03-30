package core

import "github.com/google/uuid"

// Command is a domain command that is sent to a Dispatcher.
//
// A command name should 1) be in present tense and 2) contain the intent
// (MoveCustomer vs CorrectCustomerAddress).
//
// The command should contain all the data needed when handling it as fields.
// These fields can take an optional "eh" tag, which adds properties. For now
// only "optional" is a valid tag: `eh:"optional"`.
type Command interface {
	// AggregateID returns the ID of the aggregate that the command should be
	// handled by.
	AggregateID() uuid.UUID

	// AggregateType returns the type of the aggregate that the command can be
	// handled by.
	AggregateType() AggregateType

	// CommandType returns the type of the command.
	CommandType() CommandType
}

// CommandType is the type of a command, used as its unique identifier.
type CommandType string

// String returns the string representation of a command type.
func (ct CommandType) String() string {
	return string(ct)
}

// CommandIDer provides a unique command ID to be used for request tracking etc.
type CommandIDer interface {
	// CommandID returns the ID of the command instance being handled.
	CommandID() uuid.UUID
}
