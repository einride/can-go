package identifiers

import (
	"testing"

	"gotest.tools/v3/assert"
)

func TestIsAlphaChar(t *testing.T) {
	assert.Assert(t, IsAlphaChar('b'))
	assert.Assert(t, IsAlphaChar('C'))
	assert.Assert(t, !IsAlphaChar('Ö'))
	assert.Assert(t, !IsAlphaChar('_'))
}

func TestIsNumChar(t *testing.T) {
	assert.Assert(t, IsNumChar('0'))
	assert.Assert(t, IsNumChar('1'))
	assert.Assert(t, IsNumChar('2'))
	assert.Assert(t, IsNumChar('9'))
	assert.Assert(t, !IsNumChar('/'))
	assert.Assert(t, !IsNumChar('a'))
}
