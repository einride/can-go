package socketcan

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCanRawAddr_Network(t *testing.T) {
	addr := &canRawAddr{device: "can0"}
	require.Equal(t, "can0", addr.String())
}

func TestCanRawAddr_String(t *testing.T) {
	addr := &canRawAddr{device: "can0"}
	require.Equal(t, "can", addr.Network())
}
