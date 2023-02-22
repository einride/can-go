//go:build linux && go1.12

package socketcan

import (
	"context"
	"fmt"
	"net"
	"runtime"
	"testing"
	"time"

	"github.com/blueinnovationsgroup/can-go"
	"golang.org/x/sync/errgroup"
	"gotest.tools/v3/assert"
	is "gotest.tools/v3/assert/cmp"
)

func TestDial_CANRaw(t *testing.T) {
	requireVCAN0(t)
	conn, err := Dial("can", "vcan0")
	assert.NilError(t, err)
	assert.NilError(t, conn.Close())
}

func TestDialContext_CANRaw(t *testing.T) {
	requireVCAN0(t)
	ctx, done := context.WithTimeout(context.Background(), time.Second)
	defer done()
	conn, err := DialContext(ctx, "can", "vcan0")
	assert.NilError(t, err)
	assert.NilError(t, conn.Close())
}

func TestConn_DialFail(t *testing.T) {
	t.Run("bad file name", func(t *testing.T) {
		_, err := Dial("can", "badFileName#")
		assert.ErrorContains(t, err, "dial")
	})
	t.Run("timeout", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		cancel()
		_, err := DialContext(ctx, "can", "vcan0")
		assert.ErrorContains(t, err, "context canceled")
	})
}

func TestConn_Addr(t *testing.T) {
	requireVCAN0(t)
	conn, err := Dial("can", "vcan0")
	assert.NilError(t, err)
	assert.Assert(t, is.Nil(conn.LocalAddr())) // SocketCAN connections don't have a local connection
	assert.Equal(t, "can", conn.RemoteAddr().Network())
	assert.Equal(t, "vcan0", conn.RemoteAddr().String())
}

func TestConn_SetDeadline(t *testing.T) {
	requireVCAN0(t)
	// Given that a vcan device exists and that I can open a connection to it
	receiver, err := Dial("can", "vcan0")
	assert.NilError(t, err)
	// When I set the can
	timeout := 20 * time.Millisecond
	assert.NilError(t, receiver.SetDeadline(time.Now().Add(timeout)))
	// Then I expect a read without a corresponding write to time out
	data := make([]byte, lengthOfFrame)
	n, err := receiver.Read(data)
	assert.Equal(t, 0, n)
	assert.Assert(t, is.ErrorContains(err, ""))
	// When I clear the timeouts
	assert.NilError(t, receiver.SetDeadline(time.Time{}))
	// Then I don't expect the read to timeout anymore
	errChan := make(chan error, 1)
	go func() {
		_, err = receiver.Read(data)
		errChan <- err
	}()
	select {
	case <-errChan:
		t.Fatal("unexpected read result")
	case <-time.After(timeout):
		assert.NilError(t, receiver.Close())
		assert.Assert(t, is.ErrorContains(<-errChan, ""))
	}
}

func TestConn_ReadWrite(t *testing.T) {
	requireVCAN0(t)
	// given a reader and writer
	reader, err := Dial("can", "vcan0")
	assert.NilError(t, err)
	writer, err := Dial("can", "vcan0")
	assert.NilError(t, err)
	// when the reader reads
	var g errgroup.Group
	var readFrame can.Frame
	g.Go(func() error {
		rec := NewReceiver(reader)
		if !rec.Receive() {
			return fmt.Errorf("receive")
		}
		readFrame = rec.Frame()
		return reader.Close()
	})
	// and the writer writes
	writeFrame := can.Frame{ID: 32}
	tr := NewTransmitter(writer)
	ctx, done := context.WithTimeout(context.Background(), time.Second)
	defer done()
	assert.NilError(t, tr.TransmitFrame(ctx, writeFrame))
	assert.NilError(t, writer.Close())
	// then the written and read frames should be identical
	assert.NilError(t, g.Wait())
	assert.DeepEqual(t, writeFrame, readFrame)
}

func TestConn_WriteOnClosedFails(t *testing.T) {
	requireVCAN0(t)
	conn, err := Dial("can", "vcan0")
	assert.NilError(t, err)
	tr := NewTransmitter(conn)
	ctx, done := context.WithTimeout(context.Background(), time.Second)
	defer done()
	assert.NilError(t, tr.TransmitFrame(ctx, can.Frame{}))
	// When I close the connection and then write to it
	assert.NilError(t, conn.Close())
	// Then it should fail
	assert.Assert(t, is.ErrorContains(tr.TransmitFrame(ctx, can.Frame{}), ""), "WriteFrame on a closed Conn should fail")
}

func TestConn_ReadOnClose(t *testing.T) {
	requireVCAN0(t)
	t.Run("close then read", func(t *testing.T) {
		conn, err := Dial("can", "vcan0")
		assert.NilError(t, err)
		// When I close the connection and then read from it
		assert.NilError(t, conn.Close())
		rec := NewReceiver(conn)
		assert.Assert(t, !rec.Receive())
		assert.Assert(t, is.ErrorContains(rec.Err(), ""))
	})
	t.Run("read then close", func(t *testing.T) {
		conn, err := Dial("can", "vcan0")
		assert.NilError(t, err)
		// And when I read from a connection
		var g errgroup.Group
		var receiveErr error
		g.Go(func() error {
			rec := NewReceiver(conn)
			if rec.Receive() {
				return fmt.Errorf("receive")
			}
			receiveErr = rec.Err()
			return nil
		})
		runtime.Gosched()
		// And then close it
		assert.NilError(t, conn.Close())
		// Then the read operation should fail
		assert.NilError(t, g.Wait())
		assert.Assert(t, is.ErrorContains(receiveErr, ""))
	})
}

func requireVCAN0(t *testing.T) {
	t.Helper()
	if _, err := net.InterfaceByName("vcan0"); err != nil {
		t.Skip("device vcan0 not available")
	}
}
