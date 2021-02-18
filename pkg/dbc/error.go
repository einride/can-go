package dbc

import (
	"errors"
	"fmt"
	"text/scanner"
)

// Error represents an error in a DBC file.
type Error interface {
	error

	// Position of the error in the DBC file.
	Position() scanner.Position

	// Reason for the error.
	Reason() string
}

// validationError is an error resulting from an invalid DBC definition.
type validationError struct {
	pos    scanner.Position
	reason string
	cause  error
}

func (e *validationError) Unwrap() error {
	return e.cause
}

var _ Error = &validationError{}

func (e *validationError) Error() string {
	return fmt.Sprintf("%v: %s (validate)", e.Position(), e.reason)
}

// Reason returns the reason for the error.
func (e *validationError) Reason() string {
	return e.reason
}

// Position returns the position of the validation error in the DBC file.
//
// When the validation error results from an invalid nested definition, the position of the nested definition is
// returned.
func (e *validationError) Position() scanner.Position {
	var errValidation *validationError
	if errors.As(e.cause, &errValidation) {
		return errValidation.Position()
	}
	return e.pos
}

// parseError is an error resulting from a failure to parse a DBC file.
type parseError struct {
	pos    scanner.Position
	reason string
}

var _ Error = &parseError{}

func (e *parseError) Error() string {
	return fmt.Sprintf("%v: %s (parse)", e.pos, e.reason)
}

// Reason returns the reason for the error.
func (e *parseError) Reason() string {
	return e.reason
}

// Position returns the position of the parse error in the DBC file.
func (e *parseError) Position() scanner.Position {
	return e.pos
}
