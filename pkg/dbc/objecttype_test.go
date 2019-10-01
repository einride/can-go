package dbc

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestObjectType_Validate(t *testing.T) {
	for _, tt := range []ObjectType{
		ObjectTypeUnspecified,
		ObjectTypeNetworkNode,
		ObjectTypeMessage,
		ObjectTypeSignal,
		ObjectTypeEnvironmentVariable,
	} {
		require.NoError(t, tt.Validate())
	}
}

func TestObjectType_Validate_Error(t *testing.T) {
	require.Error(t, ObjectType("foo").Validate())
}
