package dbc

import (
	"testing"

	"gotest.tools/v3/assert"
)

func TestObjectType_Validate(t *testing.T) {
	for _, tt := range []ObjectType{
		ObjectTypeUnspecified,
		ObjectTypeNetworkNode,
		ObjectTypeMessage,
		ObjectTypeSignal,
		ObjectTypeEnvironmentVariable,
	} {
		assert.NilError(t, tt.Validate())
	}
}

func TestObjectType_Validate_Error(t *testing.T) {
	assert.ErrorContains(t, ObjectType("foo").Validate(), "invalid object type")
}
