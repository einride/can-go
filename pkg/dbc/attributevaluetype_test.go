package dbc

import (
	"testing"

	"gotest.tools/v3/assert"
)

func TestAttributeValueType_Validate(t *testing.T) {
	for _, tt := range []AttributeValueType{
		AttributeValueTypeInt,
		AttributeValueTypeHex,
		AttributeValueTypeFloat,
		AttributeValueTypeString,
		AttributeValueTypeEnum,
	} {
		assert.NilError(t, tt.Validate())
	}
}

func TestAttributeValueType_Validate_Error(t *testing.T) {
	assert.ErrorContains(t, AttributeValueType("foo").Validate(), "invalid attribute value type")
}
