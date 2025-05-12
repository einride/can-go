package dbc

import (
	"fmt"
	"strings"
	"testing"

	"gotest.tools/v3/assert"
)

func TestIdentifier_Validate(t *testing.T) {
	for _, tt := range []Identifier{
		"_",
		"_foo",
		"foo",
		"foo32",
		"_43",
		Identifier(strings.Repeat("a", maxIdentifierLength)),
	} {
		t.Run(fmt.Sprintf("%v", tt), func(t *testing.T) {
			assert.NilError(t, tt.Validate())
		})
	}
}

func TestIdentifier_Validate_Error(t *testing.T) {
	for _, tt := range []Identifier{
		"42",
		"",
		"42foo",
		"☃",
		"foo☃",
		"foo bar",
		Identifier(strings.Repeat("a", maxIdentifierLength+1)),
	} {
		t.Run(fmt.Sprintf("%v", tt), func(t *testing.T) {
			assert.ErrorContains(t, tt.Validate(), "invalid identifier")
		})
	}
}
