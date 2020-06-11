package descriptor

import (
	"testing"

	"gotest.tools/v3/assert"
	is "gotest.tools/v3/assert/cmp"
)

func TestMessage_MultiplexerSignal(t *testing.T) {
	mux := &Signal{
		Name:          "Mux",
		IsMultiplexer: true,
	}
	m := &Message{
		Signals: []*Signal{
			{Name: "NotMux"},
			mux,
			{Name: "AlsoNotMux"},
		},
	}
	actualMux, ok := m.MultiplexerSignal()
	assert.Assert(t, ok)
	assert.DeepEqual(t, mux, actualMux)
}

func TestMessage_MultiplexerSignal_NotFound(t *testing.T) {
	m := &Message{
		Signals: []*Signal{
			{Name: "NotMux"},
			{Name: "AlsoNotMux"},
		},
	}
	actualMux, ok := m.MultiplexerSignal()
	assert.Assert(t, !ok)
	assert.Assert(t, is.Nil(actualMux))
}
