package socketcan

import (
	"context"
	"fmt"
	"net"

	"go.einride.tech/can"
)

type TransmitterOption func(*transmitterOpts)

type transmitterOpts struct {
	frameInterceptor FrameInterceptor
}

// Transmitter transmits CAN frames.
type Transmitter struct {
	opts transmitterOpts
	conn net.Conn
}

// NewTransmitter creates a new transmitter that transmits CAN frames to the provided io.Writer.
func NewTransmitter(conn net.Conn, opt ...TransmitterOption) *Transmitter {
	opts := transmitterOpts{}
	for _, f := range opt {
		f(&opts)
	}
	return &Transmitter{
		conn: conn,
		opts: opts,
	}
}

// TransmitMessage transmits a CAN message.
func (t *Transmitter) TransmitMessage(ctx context.Context, m can.Message) error {
	f, err := m.MarshalFrame()
	if err != nil {
		return fmt.Errorf("transmit message: %w", err)
	}
	return t.TransmitFrame(ctx, f)
}

// TransmitFrame transmits a CAN frame.
func (t *Transmitter) TransmitFrame(ctx context.Context, f can.Frame) error {
	var scf frame
	scf.encodeFrame(f)
	data := make([]byte, lengthOfFrame)
	scf.marshalBinary(data)
	if deadline, ok := ctx.Deadline(); ok {
		if err := t.conn.SetWriteDeadline(deadline); err != nil {
			return fmt.Errorf("transmit frame: %w", err)
		}
	}
	if _, err := t.conn.Write(data); err != nil {
		return fmt.Errorf("transmit frame: %w", err)
	}
	if t.opts.frameInterceptor != nil {
		t.opts.frameInterceptor(f)
	}
	return nil
}

// Close the transmitter's underlying connection.
func (t *Transmitter) Close() error {
	return t.conn.Close()
}

// TransmitterFrameInterceptor returns a TransmitterOption that sets the FrameInterceptor for the
// transmitter. Only one frame interceptor can be installed.
func TransmitterFrameInterceptor(i FrameInterceptor) TransmitterOption {
	return func(o *transmitterOpts) {
		o.frameInterceptor = i
	}
}
