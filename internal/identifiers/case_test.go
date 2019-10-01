package identifiers

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIsCamelCase(t *testing.T) {
	require.True(t, IsCamelCase("SOC"))
	require.True(t, IsCamelCase("Camel"))
	require.True(t, IsCamelCase("CamelCase"))
	require.False(t, IsCamelCase("camelCase"))
	require.False(t, IsCamelCase("snake_case"))
	require.False(t, IsCamelCase("kebab-case"))
}
