//go:build linux && go1.12

package socketcan

import (
	"fmt"
	"net"
	"os"

	"golang.org/x/sys/unix"
)

type dialOpts struct {
	errorFrameMask *int
}

func dialRaw(device string, opt ...DialOption) (conn net.Conn, err error) {
	defer func() {
		if err != nil {
			err = &net.OpError{Op: "dial", Net: canRawNetwork, Addr: &canRawAddr{device: device}, Err: err}
		}
	}()
	opts := dialOpts{}
	for _, f := range opt {
		f(&opts)
	}
	ifi, err := net.InterfaceByName(device)
	if err != nil {
		return nil, fmt.Errorf("interface %s: %w", device, err)
	}
	fd, err := unix.Socket(unix.AF_CAN, unix.SOCK_RAW, unix.CAN_RAW)
	if err != nil {
		return nil, fmt.Errorf("socket: %w", err)
	}
	if opts.errorFrameMask != nil {
		if err := unix.SetsockoptInt(fd, unix.SOL_CAN_RAW, unix.CAN_RAW_ERR_FILTER, *opts.errorFrameMask); err != nil {
			return nil, fmt.Errorf("set error filter: %w", err)
		}
	}
	// put fd in non-blocking mode so the created file will be registered by the runtime poller (Go >= 1.12)
	if err := unix.SetNonblock(fd, true); err != nil {
		return nil, fmt.Errorf("set nonblock: %w", err)
	}
	if err := unix.Bind(fd, &unix.SockaddrCAN{Ifindex: ifi.Index}); err != nil {
		return nil, fmt.Errorf("bind: %w", err)
	}
	return &fileConn{ra: &canRawAddr{device: device}, f: os.NewFile(uintptr(fd), "can")}, nil
}

// WithReceiveErrorFrames returns a DialOption which enables
// can error frame receiving on can port.
func WithReceiveErrorFrames() DialOption {
	return func(o *dialOpts) {
		canErrMask := unix.CAN_ERR_MASK
		o.errorFrameMask = &canErrMask
	}
}
