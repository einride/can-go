package dbc

import (
	"testing"

	"gotest.tools/v3/assert"
)

func TestAccessType_Validate(t *testing.T) {
	for _, tt := range []AccessType{
		AccessTypeUnrestricted,
		AccessTypeRead,
		AccessTypeWrite,
		AccessTypeReadWrite,
	} {
		assert.NilError(t, tt.Validate())
	}
}

func TestAccessType_Validate_Error(t *testing.T) {
	assert.ErrorContains(t, AccessType("foo").Validate(), "invalid access type")
}
