package dbc

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSignalValueType_Validate(t *testing.T) {
	for _, tt := range []SignalValueType{
		SignalValueTypeInt,
		SignalValueTypeFloat32,
		SignalValueTypeFloat64,
	} {
		require.NoError(t, tt.Validate())
	}
}

func TestSignalValueType_Validate_Error(t *testing.T) {
	require.Error(t, SignalValueType(42).Validate())
}
