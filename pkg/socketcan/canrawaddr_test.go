package socketcan

import (
	"testing"

	"gotest.tools/v3/assert"
)

func TestCanRawAddr_Network(t *testing.T) {
	addr := &canRawAddr{device: "can0"}
	assert.Equal(t, "can0", addr.String())
}

func TestCanRawAddr_String(t *testing.T) {
	addr := &canRawAddr{device: "can0"}
	assert.Equal(t, "can", addr.Network())
}
