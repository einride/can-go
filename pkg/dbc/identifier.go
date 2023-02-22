package dbc

import (
	"fmt"
	"unicode"

	"github.com/blueinnovationsgroup/can-go/internal/identifiers"
)

// Identifier represents a DBC identifier.
type Identifier string

// maxIdentifierLength is the length of the longest valid DBC identifier.
const maxIdentifierLength = 128

// Validate returns an error for invalid DBC identifiers.
func (id Identifier) Validate() (err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("invalid identifier '%s': %w", id, err)
		}
	}()
	if len(id) == 0 {
		return fmt.Errorf("zero-length")
	}
	if len(id) > maxIdentifierLength {
		return fmt.Errorf("length %v exceeds max length %v", len(id), maxIdentifierLength)
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

func (id Identifier) FirstCharUpper() (Identifier, bool) {
	firstChar := rune(id[0])
	if identifiers.IsAlphaChar(firstChar) && !identifiers.IsUpperAlphaChar(firstChar) {
		conv := []byte(id)
		conv[0] = byte(unicode.ToUpper(firstChar))
		return Identifier(conv), true
	} else {
		return id, false
	}
}
