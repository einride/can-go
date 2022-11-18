//go:build linux && go1.12

package socketcan

import (
	"fmt"
	"net"
	"os"

	"golang.org/x/sys/unix"
)

func dialRaw(device string) (conn net.Conn, err error) {
	defer func() {
		if err != nil {
			err = &net.OpError{Op: "dial", Net: canRawNetwork, Addr: &canRawAddr{device: device}, Err: err}
		}
	}()
	ifi, err := net.InterfaceByName(device)
	if err != nil {
		return nil, fmt.Errorf("interface %s: %w", device, err)
	}
	fd, err := unix.Socket(unix.AF_CAN, unix.SOCK_RAW, unix.CAN_RAW)
	if err != nil {
		return nil, fmt.Errorf("socket: %w", err)
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
