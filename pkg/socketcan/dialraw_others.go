//go:build !linux || !go1.12

package socketcan

import (
	"fmt"
	"net"
	"runtime"
)

type dialOpts struct {
}

func dialRaw(interfaceName string, opt ...DialOption) (net.Conn, error) {
	return nil, fmt.Errorf("SocketCAN not supported on OS %s and runtime %s", runtime.GOOS, runtime.Version())
}

func WithReceiveErrorFrames() DialOption {
	return func(o *dialOpts) {
	}
}
