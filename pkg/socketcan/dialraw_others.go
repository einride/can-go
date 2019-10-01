// +build !linux !go1.12

package socketcan

import (
	"fmt"
	"net"
	"runtime"
)

func dialRaw(interfaceName string) (net.Conn, error) {
	return nil, fmt.Errorf("SocketCAN not supported on OS %s and runtime %s", runtime.GOOS, runtime.Version())
}
