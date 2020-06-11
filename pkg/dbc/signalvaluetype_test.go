package dbc

import (
	"testing"

	"gotest.tools/v3/assert"
)

func TestSignalValueType_Validate(t *testing.T) {
	for _, tt := range []SignalValueType{
		SignalValueTypeInt,
		SignalValueTypeFloat32,
		SignalValueTypeFloat64,
	} {
		assert.NilError(t, tt.Validate())
	}
}

func TestSignalValueType_Validate_Error(t *testing.T) {
	assert.Error(t, SignalValueType(42).Validate(), "invalid signal value type: 42")
}
