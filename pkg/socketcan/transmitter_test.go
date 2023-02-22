package socketcan

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"testing"
	"time"

	"github.com/blueinnovationsgroup/can-go"
	"golang.org/x/sync/errgroup"
	"gotest.tools/v3/assert"
)

func TestTransmitter_TransmitMessage(t *testing.T) {
	testTransmit := func(opt TransmitterOption) {
		w, r := net.Pipe()
		f := can.Frame{
			ID:     0x12,
			Length: 8,
			Data:   can.Data{0x12, 0x34, 0x56, 0x78, 0x9a, 0xbc, 0xde, 0xf0},
		}
		msg := &testMessage{frame: f}
		expected := []byte{
			// id---------------> | dlc | padding-------> | data----------------------------------------> |
			0x12, 0x00, 0x00, 0x00, 0x08, 0x00, 0x00, 0x00, 0x12, 0x34, 0x56, 0x78, 0x9a, 0xbc, 0xde, 0xf0,
		}
		// write
		var g errgroup.Group
		g.Go(func() error {
			tr := NewTransmitter(w, opt)
			ctx, done := context.WithTimeout(context.Background(), time.Second)
			defer done()
			if err := tr.TransmitMessage(ctx, msg); err != nil {
				return err
			}
			return w.Close()
		})
		// read
		actual := make([]byte, len(expected))
		_, err := io.ReadFull(r, actual)
		assert.NilError(t, err)
		assert.NilError(t, r.Close())
		// assert
		assert.DeepEqual(t, expected, actual)
		assert.NilError(t, g.Wait())
	}

	// No opts
	testTransmit(func(*transmitterOpts) {})

	// Frame Interceptor
	run := false
	intFunc := func(fr can.Frame) {
		run = true
	}
	testTransmit(TransmitterFrameInterceptor(intFunc))
	assert.Assert(t, run)
}

func TestTransmitter_TransmitMessage_Error(t *testing.T) {
	cause := fmt.Errorf("boom")
	msg := &testMessage{err: cause}
	tr := NewTransmitter(nil)
	ctx, done := context.WithTimeout(context.Background(), time.Second)
	defer done()
	err := tr.TransmitMessage(ctx, msg)
	assert.Error(t, err, "transmit message: boom")
	assert.Equal(t, cause, errors.Unwrap(err))
}

func TestTransmitter_TransmitFrame_Error(t *testing.T) {
	t.Run("set deadline", func(t *testing.T) {
		cause := fmt.Errorf("boom")
		w := &errCon{deadlineErr: cause}
		tr := NewTransmitter(w)
		ctx, done := context.WithTimeout(context.Background(), time.Second)
		defer done()
		err := tr.TransmitFrame(ctx, can.Frame{})
		assert.ErrorContains(t, err, "boom")
		assert.Equal(t, cause, errors.Unwrap(err))
	})
	t.Run("write", func(t *testing.T) {
		cause := fmt.Errorf("boom")
		w := &errCon{writeErr: cause}
		tr := NewTransmitter(w)
		ctx, done := context.WithTimeout(context.Background(), time.Second)
		defer done()
		err := tr.TransmitFrame(ctx, can.Frame{})
		assert.ErrorContains(t, err, "boom")
		assert.Equal(t, cause, errors.Unwrap(err))
	})
}

type testMessage struct {
	frame can.Frame
	err   error
}

func (t *testMessage) MarshalFrame() (can.Frame, error) {
	return t.frame, t.err
}

func (t *testMessage) UnmarshalFrame(can.Frame) error {
	panic("should not be called")
}

type errCon struct {
	deadlineErr error
	writeErr    error
}

func (e *errCon) Write(b []byte) (n int, err error) {
	return 0, e.writeErr
}

func (e *errCon) SetWriteDeadline(t time.Time) error {
	return e.deadlineErr
}

func (e *errCon) Read(b []byte) (n int, err error) {
	panic("should not be called")
}

func (e *errCon) Close() error {
	panic("should not be called")
}

func (e *errCon) LocalAddr() net.Addr {
	panic("should not be called")
}

func (e *errCon) RemoteAddr() net.Addr {
	panic("should not be called")
}

func (e *errCon) SetDeadline(t time.Time) error {
	panic("should not be called")
}

func (e *errCon) SetReadDeadline(t time.Time) error {
	panic("should not be called")
}
