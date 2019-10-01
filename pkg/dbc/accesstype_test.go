package dbc

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAccessType_Validate(t *testing.T) {
	for _, tt := range []AccessType{
		AccessTypeUnrestricted,
		AccessTypeRead,
		AccessTypeWrite,
		AccessTypeReadWrite,
	} {
		require.NoError(t, tt.Validate())
	}
}

func TestAccessType_Validate_Error(t *testing.T) {
	require.Error(t, AccessType("foo").Validate())
}
