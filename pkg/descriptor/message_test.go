package descriptor

import (
	"testing"

	"github.com/stretchr/testify/require"
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
	require.True(t, ok)
	require.Equal(t, mux, actualMux)
}

func TestMessage_MultiplexerSignal_NotFound(t *testing.T) {
	m := &Message{
		Signals: []*Signal{
			{Name: "NotMux"},
			{Name: "AlsoNotMux"},
		},
	}
	actualMux, ok := m.MultiplexerSignal()
	require.False(t, ok)
	require.Nil(t, actualMux)
}
