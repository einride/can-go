package socketcan

import (
	"context"
	"net"
	"testing"
	"time"

	"github.com/blueinnovationsgroup/can-go"
	"golang.org/x/sync/errgroup"
	"gotest.tools/v3/assert"
)

func TestDial_TCP(t *testing.T) {
	lis, err := net.Listen("tcp", "localhost:0")
	assert.NilError(t, err)
	var g errgroup.Group
	g.Go(func() error {
		conn, err := lis.Accept()
		if err != nil {
			return err
		}
		return conn.Close()
	})
	conn, err := Dial("tcp", lis.Addr().String())
	assert.NilError(t, err)
	assert.NilError(t, conn.Close())
	assert.NilError(t, g.Wait())
}

func TestDialContext_TCP(t *testing.T) {
	lis, err := net.Listen("tcp", "localhost:0")
	assert.NilError(t, err)
	var g errgroup.Group
	g.Go(func() error {
		conn, err := lis.Accept()
		if err != nil {
			return err
		}
		return conn.Close()
	})
	ctx, done := context.WithTimeout(context.Background(), time.Second)
	defer done()
	conn, err := DialContext(ctx, "tcp", lis.Addr().String())
	assert.NilError(t, err)
	assert.NilError(t, conn.Close())
	assert.NilError(t, g.Wait())
}

func TestConn_TransmitReceiveTCP(t *testing.T) {
	// Given: A TCP listener that writes a frame on an accepted connection
	lis, err := net.Listen("tcp", "localhost:0")
	assert.NilError(t, err)
	var g errgroup.Group
	frame := can.Frame{ID: 42, Length: 5, Data: can.Data{'H', 'e', 'l', 'l', 'o'}}
	g.Go(func() error {
		conn, err := lis.Accept()
		if err != nil {
			return err
		}
		tr := NewTransmitter(conn)
		ctx, done := context.WithTimeout(context.Background(), time.Second)
		defer done()
		if err := tr.TransmitFrame(ctx, frame); err != nil {
			return err
		}
		return conn.Close()
	})
	// When: We connect to the listener
	ctx, done := context.WithTimeout(context.Background(), time.Second)
	defer done()
	conn, err := DialContext(ctx, "tcp", lis.Addr().String())
	assert.NilError(t, err)
	rec := NewReceiver(conn)
	assert.Assert(t, rec.Receive())
	assert.Assert(t, !rec.HasErrorFrame())
	assert.DeepEqual(t, frame, rec.Frame())
	assert.NilError(t, conn.Close())
	assert.NilError(t, g.Wait())
}
