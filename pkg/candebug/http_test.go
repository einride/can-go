package candebug

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/blueinnovationsgroup/can-go"
	"github.com/blueinnovationsgroup/can-go/pkg/descriptor"
	"github.com/blueinnovationsgroup/can-go/pkg/generated"
	"gotest.tools/v3/assert"
)

func TestServeMessagesHTTP_Single(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ServeMessagesHTTP(w, r, []generated.Message{
			&testMessage{
				frame:      can.Frame{ID: 100, Length: 1},
				descriptor: newDriverHeartbeatDescriptor(),
			},
		})
	}))
	c := http.DefaultClient
	ctx := context.Background()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, ts.URL, nil)
	assert.NilError(t, err)
	res, err := c.Do(req)
	assert.NilError(t, err)
	response, err := io.ReadAll(res.Body)
	assert.NilError(t, err)
	assert.NilError(t, res.Body.Close())
	const expected = `
DriverHeartbeat
===============
ID: 100 (0x64)
Sender: DRIVER
SendType: Cyclic
CycleTime: 100ms
DelayTime: 2s
===============
Command: 0 (0x0) None
`
	assert.Equal(t, strings.TrimSpace(expected), string(response))
}

func TestServeMessagesHTTP_Multi(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ServeMessagesHTTP(w, r, []generated.Message{
			&testMessage{
				frame:      can.Frame{ID: 100, Length: 1},
				descriptor: newDriverHeartbeatDescriptor(),
			},
			&testMessage{
				frame:      can.Frame{ID: 100, Length: 1, Data: can.Data{0x01}},
				descriptor: newDriverHeartbeatDescriptor(),
			},
		})
	}))
	c := http.DefaultClient
	ctx := context.Background()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, ts.URL, nil)
	assert.NilError(t, err)
	res, err := c.Do(req)
	assert.NilError(t, err)
	response, err := io.ReadAll(res.Body)
	assert.NilError(t, err)
	assert.NilError(t, res.Body.Close())
	const expected = `
DriverHeartbeat
===============
ID: 100 (0x64)
Sender: DRIVER
SendType: Cyclic
CycleTime: 100ms
DelayTime: 2s
===============
Command: 0 (0x0) None


DriverHeartbeat
===============
ID: 100 (0x64)
Sender: DRIVER
SendType: Cyclic
CycleTime: 100ms
DelayTime: 2s
===============
Command: 1 (0x1) Sync
`
	assert.Equal(t, strings.TrimSpace(expected), string(response))
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

func newDriverHeartbeatDescriptor() *descriptor.Message {
	return &descriptor.Message{
		Name:        "DriverHeartbeat",
		SenderNode:  "DRIVER",
		ID:          100,
		Length:      1,
		Description: "Sync message used to synchronize the controllers",
		SendType:    descriptor.SendTypeCyclic,
		CycleTime:   100 * time.Millisecond,
		DelayTime:   2 * time.Second,
		Signals: []*descriptor.Signal{
			{
				Name:   "Command",
				Start:  0,
				Length: 8,
				Scale:  1,
				ValueDescriptions: []*descriptor.ValueDescription{
					{Value: 0, Description: "None"},
					{Value: 1, Description: "Sync"},
					{Value: 2, Description: "Reboot"},
				},
				ReceiverNodes: []string{"SENSOR", "MOTOR"},
			},
		},
	}
}
