package dbc

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestEnvironmentVariableType_Validate(t *testing.T) {
	for _, tt := range []EnvironmentVariableType{
		EnvironmentVariableTypeInteger,
		EnvironmentVariableTypeFloat,
		EnvironmentVariableTypeString,
	} {
		require.NoError(t, tt.Validate())
	}
}

func TestEnvironmentVariableType_Validate_Error(t *testing.T) {
	require.Error(t, EnvironmentVariableType(42).Validate())
}
