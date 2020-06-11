package dbc

import (
	"testing"

	"gotest.tools/v3/assert"
)

func TestEnvironmentVariableType_Validate(t *testing.T) {
	for _, tt := range []EnvironmentVariableType{
		EnvironmentVariableTypeInteger,
		EnvironmentVariableTypeFloat,
		EnvironmentVariableTypeString,
	} {
		assert.NilError(t, tt.Validate())
	}
}

func TestEnvironmentVariableType_Validate_Error(t *testing.T) {
	assert.Error(t, EnvironmentVariableType(42).Validate(), "invalid environment variable type: 42")
}
