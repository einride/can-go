package canjson

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	examplecan "go.einride.tech/can/testdata/gen/go/example"
)

func TestMarshal(t *testing.T) {
	driverHeartbeat := examplecan.NewDriverHeartbeat().SetCommand(examplecan.DriverHeartbeat_Command_Reboot)
	js, err := Marshal(driverHeartbeat)
	require.NoError(t, err)
	expected := strings.TrimSpace(`
		{"Command":{"Raw":2,"Physical":2,"Description":"Reboot"}}
	`)
	require.Equal(t, expected, string(js))
}
