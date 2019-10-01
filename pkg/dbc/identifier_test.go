package dbc

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
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
		tt := tt
		t.Run(fmt.Sprintf("%v", tt), func(t *testing.T) {
			require.NoError(t, tt.Validate())
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
		tt := tt
		t.Run(fmt.Sprintf("%v", tt), func(t *testing.T) {
			require.Error(t, tt.Validate())
		})
	}
}
