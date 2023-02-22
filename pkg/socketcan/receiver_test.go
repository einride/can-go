package socketcan

import (
	"bytes"
	"io"
	"testing"

	"github.com/blueinnovationsgroup/can-go"
	"gotest.tools/v3/assert"
)

func TestReceiver_ReceiveFrames_Options(t *testing.T) {
	testReceive := func(opt ReceiverOption) {
		input := []byte{
			// id---------------> | dlc | padding-------> | data----------------------------------------> |
			0x01, 0x00, 0x00, 0x00, 0x02, 0x00, 0x00, 0x00, 0x12, 0x34, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		}
		expected := can.Frame{ID: 0x01, Length: 2, Data: can.Data{0x12, 0x34}}
		receiver := NewReceiver(io.NopCloser(bytes.NewReader(input)), opt)
		assert.Assert(t, receiver.Receive(), "expecting 1 CAN frames")
		assert.NilError(t, receiver.Err())
		assert.Assert(t, !receiver.HasErrorFrame())
		assert.DeepEqual(t, expected, receiver.Frame())
		assert.Assert(t, !receiver.Receive(), "expecting exactly 1 CAN frames")
		assert.NilError(t, receiver.Err())
	}

	// no options
	testReceive(func(*receiverOpts) {})

	// frame interceptor
	run := false
	intFunc := func(can.Frame) {
		run = true
	}
	testReceive(ReceiverFrameInterceptor(intFunc))
	assert.Assert(t, run)
}

func TestReceiver_ReceiveFrames(t *testing.T) {
	for _, tt := range []struct {
		msg            string
		input          []byte
		expectedFrames []can.Frame
	}{
		{
			msg:            "no data",
			input:          []byte{},
			expectedFrames: []can.Frame{},
		},
		{
			msg: "incomplete frame",
			input: []byte{
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
			},
			expectedFrames: []can.Frame{},
		},
		{
			msg: "whole single frame",
			input: []byte{
				// id---------------> | dlc | padding-------> | data----------------------------------------> |
				0x01, 0x00, 0x00, 0x00, 0x02, 0x00, 0x00, 0x00, 0x12, 0x34, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
			},
			expectedFrames: []can.Frame{
				{ID: 0x01, Length: 2, Data: can.Data{0x12, 0x34}},
			},
		},
		{
			msg: "one whole one incomplete",
			input: []byte{
				// id---------------> | dlc | padding-------> | data----------------------------------------> |
				0x01, 0x00, 0x00, 0x00, 0x02, 0x00, 0x00, 0x00, 0x12, 0x34, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x00,
			},
			expectedFrames: []can.Frame{
				{ID: 0x01, Length: 2, Data: can.Data{0x12, 0x34}},
			},
		},
		{
			msg: "two whole frames",
			input: []byte{
				// id---------------> | dlc | padding-------> | data----------------------------------------> |
				0x01, 0x00, 0x00, 0x00, 0x02, 0x00, 0x00, 0x00, 0x12, 0x34, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				// id---------------> | dlc | padding-------> | data----------------------------------------> |
				0x02, 0x00, 0x00, 0x00, 0x02, 0x00, 0x00, 0x00, 0x56, 0x78, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
			},
			expectedFrames: []can.Frame{
				{ID: 0x01, Length: 2, Data: can.Data{0x12, 0x34}},
				{ID: 0x02, Length: 2, Data: can.Data{0x56, 0x78}},
			},
		},
	} {
		tt := tt
		t.Run(tt.msg, func(t *testing.T) {
			receiver := NewReceiver(io.NopCloser(bytes.NewReader(tt.input)))
			for i, expected := range tt.expectedFrames {
				assert.Assert(t, receiver.Receive(), "expecting %d CAN frames", i+1)
				assert.NilError(t, receiver.Err())
				assert.Assert(t, !receiver.HasErrorFrame())
				assert.DeepEqual(t, expected, receiver.Frame())
			}
			assert.Assert(t, !receiver.Receive(), "expecting exactly %d CAN frames", len(tt.expectedFrames))
			assert.NilError(t, receiver.Err())
		})
	}
}

func TestReceiver_ReceiveErrorFrame(t *testing.T) {
	input := []byte{
		// frame
		// id---------------> | dlc | padding-------> | data----------------------------------------> |
		0x01, 0x00, 0x00, 0x00, 0x02, 0x00, 0x00, 0x00, 0x12, 0x34, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		// error frame
		// id---------------> | dlc | padding-------> | data----------------------------------------> |
		0x01, 0x00, 0x00, 0x20, 0x02, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		// frame
		// id---------------> | dlc | padding-------> | data----------------------------------------> |
		0x02, 0x00, 0x00, 0x00, 0x02, 0x00, 0x00, 0x00, 0x12, 0x34, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	}
	receiver := NewReceiver(io.NopCloser(bytes.NewReader(input)))
	// expect frame
	assert.Assert(t, receiver.Receive())
	assert.Assert(t, !receiver.HasErrorFrame())
	assert.Equal(t, can.Frame{ID: 0x01, Length: 2, Data: can.Data{0x12, 0x34}}, receiver.Frame())
	// expect error frame
	assert.Assert(t, receiver.Receive())
	assert.Assert(t, receiver.HasErrorFrame())
	assert.Equal(t, ErrorFrame{ErrorClass: ErrorClassTxTimeout}, receiver.ErrorFrame())
	// expect frame
	assert.Assert(t, receiver.Receive())
	assert.Assert(t, !receiver.HasErrorFrame())
	assert.Equal(t, can.Frame{ID: 0x02, Length: 2, Data: can.Data{0x12, 0x34}}, receiver.Frame())
	// expect end of stream
	assert.Assert(t, !receiver.Receive())
	assert.NilError(t, receiver.Err())
}
