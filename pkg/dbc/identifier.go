package dbc

import (
	"fmt"

	"go.einride.tech/can/internal/identifiers"
)

// Identifier represents a DBC identifier.
type Identifier string

// maxIdentifierLength is the length of the longest valid DBC identifier.
const maxIdentifierLength = 128

// Validate returns an error for invalid DBC identifiers.
func (id Identifier) Validate() error {
	if len(id) == 0 {
		return fmt.Errorf("zero-length")
	}
	if len(id) > maxIdentifierLength {
		return fmt.Errorf("length %v: exceeds max length: %v", len(id), maxIdentifierLength)
	}
	for i, r := range id {
		if i == 0 && r != '_' && !identifiers.IsAlphaChar(r) { // first char
			return fmt.Errorf("invalid first char: '%v'", r)
		} else if i > 0 && r != '_' && !identifiers.IsAlphaChar(r) && !identifiers.IsNumChar(r) {
			return fmt.Errorf("invalid char: '%v'", r)
		}
	}
	return nil
}
