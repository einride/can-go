package socketcan

import "net"

const canRawNetwork = "can"

// canRawAddr represents a CAN_RAW address.
type canRawAddr struct {
	device string
}

var _ net.Addr = &canRawAddr{}

func (a *canRawAddr) Network() string {
	return canRawNetwork
}

func (a *canRawAddr) String() string {
	return a.device
}
