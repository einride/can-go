package cantext

import (
	"strings"
	"testing"
	"time"

	"github.com/blueinnovationsgroup/can-go"
	"github.com/blueinnovationsgroup/can-go/pkg/descriptor"
	"github.com/blueinnovationsgroup/can-go/pkg/generated"
	"gotest.tools/v3/assert"
)

func TestMarshal(t *testing.T) {
	for _, tt := range []struct {
		name            string
		msg             generated.Message
		expected        string
		expectedCompact string
	}{
		{
			name: "with enum",
			msg: &testMessage{
				frame:      can.Frame{ID: 100, Length: 1, Data: can.Data{2}},
				descriptor: newDriverHeartbeatDescriptor(),
			},
			expected: `
DriverHeartbeat
	Command: 2 (0x2) Reboot
`,
			expectedCompact: `{Command: Reboot}`,
		},
		{
			name: "with unit",
			msg: &testMessage{
				frame:      can.Frame{ID: 100, Length: 3, Data: can.Data{1, 0x7b}},
				descriptor: newMotorStatusDescriptor(),
			},
			expected: `
MotorStatus
	WheelError: true
	SpeedKph: 0.123km/h (0x7b)
`,
			expectedCompact: `{WheelError: true, SpeedKph: 0.123km/h}`,
		},
	} {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Run("standard", func(t *testing.T) {
				txt := Marshal(tt.msg)
				assert.Equal(t, strings.TrimSpace(tt.expected), string(txt))
			})
			t.Run("compact", func(t *testing.T) {
				txt := MarshalCompact(tt.msg)
				assert.Equal(t, strings.TrimSpace(tt.expectedCompact), string(txt))
			})
		})
	}
}

func TestAppendID(t *testing.T) {
	const expected = "ID: 100 (0x64)"
	actual := string(AppendID([]byte{}, newDriverHeartbeatDescriptor()))
	assert.Equal(t, expected, actual)
}

func TestAppendSender(t *testing.T) {
	const expected = "Sender: DRIVER"
	actual := string(AppendSender([]byte{}, newDriverHeartbeatDescriptor()))
	assert.Equal(t, expected, actual)
}

func TestAppendSendType(t *testing.T) {
	const expected = "SendType: Cyclic"
	actual := string(AppendSendType([]byte{}, newDriverHeartbeatDescriptor()))
	assert.Equal(t, expected, actual)
}

func TestAppendCycleTime(t *testing.T) {
	const expected = "CycleTime: 100ms"
	actual := string(AppendCycleTime([]byte{}, newDriverHeartbeatDescriptor()))
	assert.Equal(t, expected, actual)
}

func TestAppendDelayTime(t *testing.T) {
	const expected = "DelayTime: 2s"
	actual := string(AppendDelayTime([]byte{}, newDriverHeartbeatDescriptor()))
	assert.Equal(t, expected, actual)
}

func TestAppendFrame(t *testing.T) {
	const expected = "Frame: 042#123456"
	actual := string(AppendFrame([]byte{}, can.Frame{ID: 0x42, Length: 3, Data: can.Data{0x12, 0x34, 0x56}}))
	assert.Equal(t, expected, actual)
}

func newDriverHeartbeatDescriptor() *descriptor.Message {
	return &descriptor.Message{
		Name:        (string)("DriverHeartbeat"),
		ID:          (uint32)(100),
		IsExtended:  (bool)(false),
		Length:      (uint8)(1),
		Description: (string)("Sync message used to synchronize the controllers"),
		SendType:    descriptor.SendTypeCyclic,
		CycleTime:   100 * time.Millisecond,
		DelayTime:   2 * time.Second,
		Signals: []*descriptor.Signal{
			{
				Name:             (string)("Command"),
				Start:            (uint8)(0),
				Length:           (uint8)(8),
				IsBigEndian:      (bool)(false),
				IsSigned:         (bool)(false),
				IsMultiplexer:    (bool)(false),
				IsMultiplexed:    (bool)(false),
				MultiplexerValue: (uint)(0),
				Offset:           (float64)(0),
				Scale:            (float64)(1),
				Min:              (float64)(0),
				Max:              (float64)(0),
				Unit:             (string)(""),
				Description:      (string)(""),
				ValueDescriptions: []*descriptor.ValueDescription{
					{
						Value:       (int64)(0),
						Description: (string)("None"),
					},
					{
						Value:       (int64)(1),
						Description: (string)("Sync"),
					},
					{
						Value:       (int64)(2),
						Description: (string)("Reboot"),
					},
				},
				ReceiverNodes: []string{
					(string)("SENSOR"),
					(string)("MOTOR"),
				},
				DefaultValue: (int)(0),
			},
		},
		SenderNode: (string)("DRIVER"),
	}
}

func newMotorStatusDescriptor() *descriptor.Message {
	return &descriptor.Message{
		Name:        (string)("MotorStatus"),
		ID:          (uint32)(400),
		IsExtended:  (bool)(false),
		Length:      (uint8)(3),
		Description: (string)(""),
		Signals: []*descriptor.Signal{
			{
				Name:              (string)("WheelError"),
				Start:             (uint8)(0),
				Length:            (uint8)(1),
				IsBigEndian:       (bool)(false),
				IsSigned:          (bool)(false),
				IsMultiplexer:     (bool)(false),
				IsMultiplexed:     (bool)(false),
				MultiplexerValue:  (uint)(0),
				Offset:            (float64)(0),
				Scale:             (float64)(1),
				Min:               (float64)(0),
				Max:               (float64)(0),
				Unit:              (string)(""),
				Description:       (string)(""),
				ValueDescriptions: ([]*descriptor.ValueDescription)(nil),
				ReceiverNodes: []string{
					(string)("DRIVER"),
					(string)("IO"),
				},
				DefaultValue: (int)(0),
			},
			{
				Name:              (string)("SpeedKph"),
				Start:             (uint8)(8),
				Length:            (uint8)(16),
				IsBigEndian:       (bool)(false),
				IsSigned:          (bool)(false),
				IsMultiplexer:     (bool)(false),
				IsMultiplexed:     (bool)(false),
				MultiplexerValue:  (uint)(0),
				Offset:            (float64)(0),
				Scale:             (float64)(0.001),
				Min:               (float64)(0),
				Max:               (float64)(0),
				Unit:              (string)("km/h"),
				Description:       (string)(""),
				ValueDescriptions: ([]*descriptor.ValueDescription)(nil),
				ReceiverNodes: []string{
					(string)("DRIVER"),
					(string)("IO"),
				},
				DefaultValue: (int)(0),
			},
		},
		SenderNode: (string)("MOTOR"),
		CycleTime:  (time.Duration)(100000000),
		DelayTime:  (time.Duration)(0),
	}
}

type testMessage struct {
	frame      can.Frame
	descriptor *descriptor.Message
}

func (m *testMessage) Frame() can.Frame {
	return m.frame
}

func (m *testMessage) Descriptor() *descriptor.Message {
	return m.descriptor
}

func (m *testMessage) MarshalFrame() (can.Frame, error) {
	panic("should not be called")
}

func (testMessage) Reset() {
	panic("should not be called")
}

func (testMessage) String() string {
	panic("should not be called")
}

func (testMessage) UnmarshalFrame(can.Frame) error {
	panic("should not be called")
}
