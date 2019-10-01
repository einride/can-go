package identifiers

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIsAlphaChar(t *testing.T) {
	require.True(t, IsAlphaChar('b'))
	require.True(t, IsAlphaChar('C'))
	require.False(t, IsAlphaChar('Ö'))
	require.False(t, IsAlphaChar('_'))
}

func TestIsNumChar(t *testing.T) {
	require.True(t, IsNumChar('0'))
	require.True(t, IsNumChar('1'))
	require.True(t, IsNumChar('2'))
	require.True(t, IsNumChar('9'))
	require.False(t, IsNumChar('/'))
	require.False(t, IsNumChar('a'))
}
