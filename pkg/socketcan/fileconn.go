package socketcan

import (
	"errors"
	"net"
	"os"
	"time"
)

// file is an interface for mocking file operations performed by fileConn.
type file interface {
	Read([]byte) (int, error)
	Write([]byte) (int, error)
	SetDeadline(time.Time) error
	SetReadDeadline(time.Time) error
	SetWriteDeadline(time.Time) error
	Close() error
}

// fileConn provides a net.Conn API for file-like types.
type fileConn struct {
	// f is the file to provide a net.Conn API for.
	f file
	// net is the connection's network.
	net string
	// la is the connection's local address, if any.
	la net.Addr
	// ra is the connection's remote address, if any.
	ra net.Addr
}

var _ net.Conn = &fileConn{}

func (c *fileConn) Read(b []byte) (int, error) {
	n, err := c.f.Read(b)
	if err != nil {
		return n, &net.OpError{Op: "read", Net: c.net, Source: c.la, Addr: c.ra, Err: unwrapPathError(err)}
	}
	return n, nil
}

func (c *fileConn) Write(b []byte) (int, error) {
	n, err := c.f.Write(b)
	if err != nil {
		return n, &net.OpError{Op: "write", Net: c.net, Source: c.la, Addr: c.ra, Err: unwrapPathError(err)}
	}
	return n, nil
}

func (c *fileConn) LocalAddr() net.Addr {
	return c.la
}

func (c *fileConn) RemoteAddr() net.Addr {
	return c.ra
}

func (c *fileConn) SetDeadline(t time.Time) error {
	if err := c.f.SetDeadline(t); err != nil {
		return &net.OpError{Op: "set deadline", Net: c.net, Source: c.la, Addr: c.ra, Err: unwrapPathError(err)}
	}
	return nil
}

func (c *fileConn) SetReadDeadline(t time.Time) error {
	if err := c.f.SetReadDeadline(t); err != nil {
		return &net.OpError{Op: "set read deadline", Net: c.net, Source: c.la, Addr: c.ra, Err: unwrapPathError(err)}
	}
	return nil
}

func (c *fileConn) SetWriteDeadline(t time.Time) error {
	if err := c.f.SetWriteDeadline(t); err != nil {
		return &net.OpError{Op: "set write deadline", Net: c.net, Source: c.la, Addr: c.ra, Err: unwrapPathError(err)}
	}
	return nil
}

func (c *fileConn) Close() error {
	if err := c.f.Close(); err != nil {
		return &net.OpError{Op: "close", Net: c.net, Source: c.la, Addr: c.ra, Err: unwrapPathError(err)}
	}
	return nil
}

// unwrapPathError unwraps one level of *os.PathError from the provided error.
func unwrapPathError(err error) error {
	var pe *os.PathError
	if errors.As(err, &pe) {
		return pe.Err
	}
	return err
}
