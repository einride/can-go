package socketcan

import (
	"context"
	"net"
)

const udp = "udp"

type DialOption func(*dialOpts)

// Dial connects to the address on the named net.
//
// Linux only: If net is "can" it creates a SocketCAN connection to the device
// (address is interpreted as a device name).
//
// If net is "udp" it assumes UDP multicast and sets up 2 connections, one for
// receiving and one for transmitting.
// See: https://golang.org/pkg/net/#Dial
func Dial(network, address string, opt ...DialOption) (net.Conn, error) {
	switch network {
	case udp:
		return udpTransceiver(network, address)
	case canRawNetwork:
		return dialRaw(address, opt...) // platform-specific
	default:
		return net.Dial(network, address)
	}
}

// DialContext connects to the address on the named net using
// the provided context.
//
// Linux only: If net is "can" it creates a SocketCAN connection to the device
// (address is interpreted as a device name).
//
// See: https://golang.org/pkg/net/#Dialer.DialContext
func DialContext(ctx context.Context, network, address string, opt ...DialOption) (net.Conn, error) {
	switch network {
	case canRawNetwork:
		return dialCtx(ctx, func() (net.Conn, error) {
			return dialRaw(address, opt...)
		})
	case udp:
		return dialCtx(ctx, func() (net.Conn, error) {
			return udpTransceiver(network, address)
		})
	default:
		var d net.Dialer
		return d.DialContext(ctx, network, address)
	}
}

func dialCtx(ctx context.Context, connProvider func() (net.Conn, error)) (net.Conn, error) {
	resultChan := make(chan struct {
		conn net.Conn
		err  error
	})
	go func() {
		conn, err := connProvider()
		resultChan <- struct {
			conn net.Conn
			err  error
		}{conn: conn, err: err}
	}()
	// wait for connection or timeout
	select {
	case result := <-resultChan:
		return result.conn, result.err
	case <-ctx.Done():
		// timeout - make sure we clean up the connection
		// error handling not possible since we've already returned
		go func() {
			result := <-resultChan
			if result.conn != nil {
				_ = result.conn.Close()
			}
		}()
		return nil, ctx.Err()
	}
}
