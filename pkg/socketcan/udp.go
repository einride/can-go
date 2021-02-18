package socketcan

import (
	"fmt"
	"net"
	"strconv"
	"time"

	"golang.org/x/net/ipv4"
	"golang.org/x/net/nettest"
)

// udpTxRx emulates a single `net.Conn` that can be used for both transmitting
// and receiving UDP multicast packets.
type udpTxRx struct {
	tx        *ipv4.PacketConn
	rx        *ipv4.PacketConn
	groupAddr *net.UDPAddr
}

func (utr *udpTxRx) Close() error {
	if err := utr.tx.Close(); err != nil {
		_ = utr.rx.Close()
		return err
	}
	return utr.rx.Close()
}

func (utr *udpTxRx) LocalAddr() net.Addr {
	return utr.rx.LocalAddr()
}

func (utr *udpTxRx) SetDeadline(t time.Time) error {
	if err := utr.rx.SetReadDeadline(t); err != nil {
		return err
	}
	return utr.tx.SetWriteDeadline(t)
}

func (utr *udpTxRx) SetReadDeadline(t time.Time) error {
	return utr.rx.SetReadDeadline(t)
}

func (utr *udpTxRx) SetWriteDeadline(t time.Time) error {
	return utr.tx.SetWriteDeadline(t)
}

func (utr *udpTxRx) Read(b []byte) (n int, err error) {
	n, _, _, err = utr.rx.ReadFrom(b)
	return
}

func (utr *udpTxRx) Write(b []byte) (n int, err error) {
	return utr.tx.WriteTo(b, nil, nil)
}

func (utr *udpTxRx) RemoteAddr() net.Addr {
	return utr.groupAddr
}

func udpTransceiver(network, address string) (*udpTxRx, error) {
	if network != udp {
		return nil, fmt.Errorf("[%v] is not a udp network", network)
	}
	ifi, err := getMulticastInterface()
	if err != nil {
		return nil, fmt.Errorf("new UDP transceiver: %w", err)
	}
	rx, groupAddr, err := udpReceiver(address, ifi)
	if err != nil {
		return nil, fmt.Errorf("new UDP transceiver: %w", err)
	}
	tx, err := udpTransmitter(groupAddr, ifi)
	if err != nil {
		return nil, fmt.Errorf("new UDP transceiver: %w", err)
	}
	return &udpTxRx{rx: rx, tx: tx, groupAddr: groupAddr}, nil
}

func getMulticastInterface() (*net.Interface, error) {
	ifi, err := nettest.RoutedInterface("ip4", net.FlagUp|net.FlagMulticast|net.FlagLoopback)
	if err == nil {
		return ifi, nil
	}
	return nettest.RoutedInterface("ip4", net.FlagUp|net.FlagMulticast)
}

func hostPortToUDPAddr(hostport string) (*net.UDPAddr, error) {
	host, portStr, err := net.SplitHostPort(hostport)
	if err != nil {
		return nil, fmt.Errorf("convert hostport to udp addr: %w", err)
	}
	port, err := strconv.Atoi(portStr)
	if err != nil {
		return nil, fmt.Errorf("convert hostport to udp addr: %w", err)
	}
	ip := net.ParseIP(host)
	return &net.UDPAddr{Port: port, IP: ip}, nil
}

func setMulticastOpts(p *ipv4.PacketConn, ifi *net.Interface, groupAddr net.Addr) error {
	if err := p.JoinGroup(ifi, groupAddr); err != nil {
		return err
	}
	if err := p.SetMulticastInterface(ifi); err != nil {
		return err
	}
	if err := p.SetMulticastLoopback(true); err != nil {
		return err
	}
	if err := p.SetMulticastTTL(0); err != nil {
		return err
	}
	return p.SetTOS(0x0)
}

func udpReceiver(address string, ifi *net.Interface) (*ipv4.PacketConn, *net.UDPAddr, error) {
	c, err := net.ListenPacket("udp4", address)
	if err != nil {
		return nil, nil, fmt.Errorf("create udp receiver: %w", err)
	}
	groupAddr, err := hostPortToUDPAddr(address)
	if err != nil {
		return nil, nil, fmt.Errorf("create udp receiver: %w", err)
	}
	// If requested port is 0, one is provided when creating the packet listener
	if groupAddr.Port == 0 {
		localAddr, err := hostPortToUDPAddr(c.LocalAddr().String())
		if err != nil {
			return nil, nil, fmt.Errorf("create udp receiver: %w", err)
		}
		groupAddr.Port = localAddr.Port
	}
	rx := ipv4.NewPacketConn(c)
	if err := setMulticastOpts(rx, ifi, groupAddr); err != nil {
		return nil, nil, fmt.Errorf("new UDP transceiver: %w", err)
	}
	return rx, groupAddr, nil
}

func udpTransmitter(groupAddr *net.UDPAddr, ifi *net.Interface) (*ipv4.PacketConn, error) {
	c, err := net.DialUDP("udp4", nil, groupAddr)
	if err != nil {
		return nil, fmt.Errorf("new UDP transmitter: %w", err)
	}
	tx := ipv4.NewPacketConn(c)
	if err := tx.SetMulticastInterface(ifi); err != nil {
		return nil, fmt.Errorf("new UDP transmitter: %w", err)
	}
	return tx, nil
}
