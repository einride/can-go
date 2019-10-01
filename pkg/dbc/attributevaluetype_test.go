package dbc

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAttributeValueType_Validate(t *testing.T) {
	for _, tt := range []AttributeValueType{
		AttributeValueTypeInt,
		AttributeValueTypeHex,
		AttributeValueTypeFloat,
		AttributeValueTypeString,
		AttributeValueTypeEnum,
	} {
		require.NoError(t, tt.Validate())
	}
}

func TestAttributeValueType_Validate_Error(t *testing.T) {
	require.Error(t, AttributeValueType("foo").Validate())
}
