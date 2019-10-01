// +build linux
// +build go1.12

package socketcan

import (
	"context"
	"fmt"
	"net"
	"runtime"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.einride.tech/can"
	"golang.org/x/sync/errgroup"
)

func TestDial_CANRaw(t *testing.T) {
	requireVCAN0(t)
	conn, err := Dial("can", "vcan0")
	require.NoError(t, err)
	require.NoError(t, conn.Close())
}

func TestDialContext_CANRaw(t *testing.T) {
	requireVCAN0(t)
	ctx, done := context.WithTimeout(context.Background(), time.Second)
	defer done()
	conn, err := DialContext(ctx, "can", "vcan0")
	require.NoError(t, err)
	require.NoError(t, conn.Close())
}

func TestConn_DialFail(t *testing.T) {
	t.Run("bad file name", func(t *testing.T) {
		_, err := Dial("can", "badFileName#")
		require.Error(t, err, "Dial to a can device that does not exist should not succeed")
	})
	t.Run("timeout", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		cancel()
		_, err := DialContext(ctx, "can", "vcan0")
		require.Error(t, err, "DialContext with closed context should not succeed")
	})
}

func TestConn_Addr(t *testing.T) {
	requireVCAN0(t)
	conn, err := Dial("can", "vcan0")
	require.NoError(t, err)
	require.Nil(t, conn.LocalAddr()) // SocketCAN connections don't have a local connection
	require.Equal(t, "can", conn.RemoteAddr().Network())
	require.Equal(t, "vcan0", conn.RemoteAddr().String())
}

func TestConn_SetDeadline(t *testing.T) {
	requireVCAN0(t)
	// Given that a vcan device exists and that I can open a connection to it
	receiver, err := Dial("can", "vcan0")
	require.NoError(t, err)
	// When I set the can
	timeout := 20 * time.Millisecond
	require.NoError(t, receiver.SetDeadline(time.Now().Add(timeout)))
	// Then I expect a read without a corresponding write to time out
	data := make([]byte, lengthOfFrame)
	n, err := receiver.Read(data)
	require.Equal(t, 0, n)
	require.Error(t, err)
	// When I clear the timeouts
	require.NoError(t, receiver.SetDeadline(time.Time{}))
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
		require.NoError(t, receiver.Close())
		require.Error(t, <-errChan)
	}
}

func TestConn_ReadWrite(t *testing.T) {
	requireVCAN0(t)
	// given a reader and writer
	reader, err := Dial("can", "vcan0")
	require.NoError(t, err)
	writer, err := Dial("can", "vcan0")
	require.NoError(t, err)
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
	require.NoError(t, tr.TransmitFrame(ctx, writeFrame))
	require.NoError(t, writer.Close())
	// then the written and read frames should be identical
	require.NoError(t, g.Wait())
	require.Equal(t, writeFrame, readFrame)
}

func TestConn_WriteOnClosedFails(t *testing.T) {
	requireVCAN0(t)
	conn, err := Dial("can", "vcan0")
	require.NoError(t, err)
	tr := NewTransmitter(conn)
	ctx, done := context.WithTimeout(context.Background(), time.Second)
	defer done()
	require.NoError(t, tr.TransmitFrame(ctx, can.Frame{}))
	// When I close the connection and then write to it
	require.NoError(t, conn.Close())
	// Then it should fail
	require.Error(t, tr.TransmitFrame(ctx, can.Frame{}), "WriteFrame on a closed Conn should fail")
}

func TestConn_ReadOnClose(t *testing.T) {
	requireVCAN0(t)
	t.Run("close then read", func(t *testing.T) {
		conn, err := Dial("can", "vcan0")
		require.NoError(t, err)
		// When I close the connection and then read from it
		require.NoError(t, conn.Close())
		rec := NewReceiver(conn)
		require.False(t, rec.Receive())
		require.Error(t, rec.Err())
	})
	t.Run("read then close", func(t *testing.T) {
		conn, err := Dial("can", "vcan0")
		require.NoError(t, err)
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
		require.NoError(t, conn.Close())
		// Then the read operation should fail
		require.NoError(t, g.Wait())
		require.Error(t, receiveErr)
	})
}

func requireVCAN0(t *testing.T) {
	if _, err := net.InterfaceByName("vcan0"); err != nil {
		t.Skip("device vcan0 not available")
	}
}
