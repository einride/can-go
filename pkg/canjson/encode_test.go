package canjson

import (
	"strings"
	"testing"

	examplecan "github.com/blueinnovationsgroup/can-go/testdata/gen/go/example"
	"gotest.tools/v3/assert"
)

func TestMarshal(t *testing.T) {
	driverHeartbeat := examplecan.NewDriverHeartbeat().SetCommand(examplecan.DriverHeartbeat_Command_Reboot)
	js, err := Marshal(driverHeartbeat)
	assert.NilError(t, err)
	expected := strings.TrimSpace(`
		{"Command":{"Raw":2,"Physical":2,"Description":"Reboot"}}
	`)
	assert.Equal(t, expected, string(js))
}
