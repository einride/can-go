package generate

import (
	"os"
	"testing"
	"time"

	"github.com/blueinnovationsgroup/can-go/pkg/descriptor"
	examplecan "github.com/blueinnovationsgroup/can-go/testdata/gen/go/example"
	"gotest.tools/v3/assert"
)

func TestCompile_ExampleDBC(t *testing.T) {
	finish := runTestInDir(t, "../..")
	defer finish()
	const exampleDBCFile = "testdata/dbc/example/example.dbc"
	exampleDatabase := &descriptor.Database{
		SourceFile: exampleDBCFile,
		Version:    "",
		Nodes: []*descriptor.Node{
			{
				Name: "DBG",
			},
			{
				Name:        "DRIVER",
				Description: "The driver controller driving the car",
			},
			{
				Name: "IO",
			},
			{
				Name:        "MOTOR",
				Description: "The motor controller of the car",
			},
			{
				Name:        "SENSOR",
				Description: "The sensor controller of the car",
			},
		},

		Messages: []*descriptor.Message{
			{
				ID:         1,
				Name:       "EmptyMessage",
				SenderNode: "DBG",
			},

			{
				ID:          100,
				Name:        "DriverHeartbeat",
				Length:      1,
				SenderNode:  "DRIVER",
				Description: "Sync message used to synchronize the controllers",
				SendType:    descriptor.SendTypeCyclic,
				CycleTime:   time.Second,
				Signals: []*descriptor.Signal{
					{
						Name:          "Command",
						Start:         0,
						Length:        8,
						Scale:         1,
						ReceiverNodes: []string{"SENSOR", "MOTOR"},
						ValueDescriptions: []*descriptor.ValueDescription{
							{Value: 0, Description: "None"},
							{Value: 1, Description: "Sync"},
							{Value: 2, Description: "Reboot"},
							{Value: 3, Description: "Headlights On"},
						},
					},
				},
			},

			{
				ID:         101,
				Name:       "MotorCommand",
				Length:     1,
				SenderNode: "DRIVER",
				SendType:   descriptor.SendTypeCyclic,
				CycleTime:  100 * time.Millisecond,
				Signals: []*descriptor.Signal{
					{
						Name:          "Steer",
						Start:         0,
						Length:        4,
						IsSigned:      true,
						Scale:         1,
						Offset:        -5,
						Min:           -5,
						Max:           5,
						ReceiverNodes: []string{"MOTOR"},
					},
					{
						Name:          "Drive",
						Start:         4,
						Length:        4,
						Scale:         1,
						Max:           9,
						ReceiverNodes: []string{"MOTOR"},
					},
				},
			},

			{
				ID:         200,
				Name:       "SensorSonars",
				Length:     8,
				SenderNode: "SENSOR",
				SendType:   descriptor.SendTypeCyclic,
				CycleTime:  100 * time.Millisecond,
				Signals: []*descriptor.Signal{
					{
						Name:          "Mux",
						IsMultiplexer: true,
						Start:         0,
						Length:        4,
						Scale:         1,
						ReceiverNodes: []string{"DRIVER", "IO"},
					},
					{
						Name:          "ErrCount",
						Start:         4,
						Length:        12,
						Scale:         1,
						ReceiverNodes: []string{"DRIVER", "IO"},
					},
					{
						Name:             "Left",
						IsMultiplexed:    true,
						MultiplexerValue: 0,
						Start:            16,
						Length:           12,
						Scale:            0.1,
						ReceiverNodes:    []string{"DRIVER", "IO"},
					},
					{
						Name:             "NoFiltLeft",
						IsMultiplexed:    true,
						MultiplexerValue: 1,
						Start:            16,
						Length:           12,
						Scale:            0.1,
						ReceiverNodes:    []string{"DBG"},
					},
					{
						Name:             "Middle",
						IsMultiplexed:    true,
						MultiplexerValue: 0,
						Start:            28,
						Length:           12,
						Scale:            0.1,
						ReceiverNodes:    []string{"DRIVER", "IO"},
					},
					{
						Name:             "NoFiltMiddle",
						IsMultiplexed:    true,
						MultiplexerValue: 1,
						Start:            28,
						Length:           12,
						Scale:            0.1,
						ReceiverNodes:    []string{"DBG"},
					},
					{
						Name:             "Right",
						IsMultiplexed:    true,
						MultiplexerValue: 0,
						Start:            40,
						Length:           12,
						Scale:            0.1,
						ReceiverNodes:    []string{"DRIVER", "IO"},
					},
					{
						Name:             "NoFiltRight",
						IsMultiplexed:    true,
						MultiplexerValue: 1,
						Start:            40,
						Length:           12,
						Scale:            0.1,
						ReceiverNodes:    []string{"DBG"},
					},
					{
						Name:             "Rear",
						IsMultiplexed:    true,
						MultiplexerValue: 0,
						Start:            52,
						Length:           12,
						Scale:            0.1,
						ReceiverNodes:    []string{"DRIVER", "IO"},
					},
					{
						Name:             "NoFiltRear",
						IsMultiplexed:    true,
						MultiplexerValue: 1,
						Start:            52,
						Length:           12,
						Scale:            0.1,
						ReceiverNodes:    []string{"DBG"},
					},
				},
			},

			{
				ID:         400,
				Name:       "MotorStatus",
				Length:     3,
				SenderNode: "MOTOR",
				SendType:   descriptor.SendTypeCyclic,
				CycleTime:  100 * time.Millisecond,
				Signals: []*descriptor.Signal{
					{
						Name:          "WheelError",
						Start:         0,
						Length:        1,
						Scale:         1,
						ReceiverNodes: []string{"DRIVER", "IO"},
					},
					{
						Name:          "SpeedKph",
						Start:         8,
						Length:        16,
						Scale:         0.001,
						Unit:          "km/h",
						ReceiverNodes: []string{"DRIVER", "IO"},
					},
				},
			},

			{
				ID:         500,
				Name:       "IODebug",
				Length:     6,
				SenderNode: "IO",
				SendType:   descriptor.SendTypeEvent,
				Signals: []*descriptor.Signal{
					{
						Name:          "TestUnsigned",
						Start:         0,
						Length:        8,
						Scale:         1,
						ReceiverNodes: []string{"DBG"},
					},
					{
						Name:          "TestEnum",
						Start:         8,
						Length:        6,
						Scale:         1,
						ReceiverNodes: []string{"DBG"},
						DefaultValue:  int(examplecan.IODebug_TestEnum_Two),
						ValueDescriptions: []*descriptor.ValueDescription{
							{Value: 1, Description: "One"},
							{Value: 2, Description: "Two"},
						},
					},
					{
						Name:          "TestSigned",
						Start:         16,
						Length:        8,
						IsSigned:      true,
						Scale:         1,
						ReceiverNodes: []string{"DBG"},
					},
					{
						Name:          "TestFloat",
						Start:         24,
						Length:        8,
						Scale:         0.5,
						ReceiverNodes: []string{"DBG"},
					},
					{
						Name:          "TestBoolEnum",
						Start:         32,
						Length:        1,
						Scale:         1,
						ReceiverNodes: []string{"DBG"},
						ValueDescriptions: []*descriptor.ValueDescription{
							{Value: 0, Description: "Zero"},
							{Value: 1, Description: "One"},
						},
					},
					{
						Name:          "TestScaledEnum",
						Start:         40,
						Length:        2,
						Scale:         2,
						Min:           0,
						Max:           6,
						ReceiverNodes: []string{"DBG"},
						ValueDescriptions: []*descriptor.ValueDescription{
							{Value: 0, Description: "Zero"},
							{Value: 1, Description: "Two"},
							{Value: 2, Description: "Four"},
							{Value: 3, Description: "Six"},
						},
					},
				},
			},
		},
	}
	input, err := os.ReadFile(exampleDBCFile)
	assert.NilError(t, err)
	result, err := Compile(exampleDBCFile, input)
	if err != nil {
		t.Fatal(err)
	}
	if len(result.Warnings) > 0 {
		t.Fatal(result.Warnings)
	}
	assert.DeepEqual(t, exampleDatabase, result.Database)
}
