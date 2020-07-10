package identifiers

import (
	"testing"

	"gotest.tools/v3/assert"
)

func TestIsCamelCase(t *testing.T) {
	assert.Assert(t, IsCamelCase("SOC"))
	assert.Assert(t, IsCamelCase("Camel"))
	assert.Assert(t, IsCamelCase("CamelCase"))
	assert.Assert(t, IsCamelCase("111CamelCaseNr"))
	assert.Assert(t, !IsCamelCase("camelCase"))
	assert.Assert(t, !IsCamelCase("snake_case"))
	assert.Assert(t, !IsCamelCase("kebab-case"))
	assert.Assert(t, !IsCamelCase("111camelCaseNr"))
}
