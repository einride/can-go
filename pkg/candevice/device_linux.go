//go:build linux && go1.18

package candevice

import (
	"fmt"
	"net"
	"unsafe"

	"github.com/mdlayher/netlink"
	"github.com/mdlayher/netlink/nlenc"
	"golang.org/x/sys/unix"
)

const (
	CanLinkType  = "can"
	VcanLinkType = "vcan"
)

const (
	StateErrorActive  = unix.CAN_STATE_ERROR_ACTIVE
	StateErrorWarning = unix.CAN_STATE_ERROR_WARNING
	StateErrorPassive = unix.CAN_STATE_ERROR_PASSIVE
	StateBusOff       = unix.CAN_STATE_BUS_OFF
	StateStopped      = unix.CAN_STATE_STOPPED
	StateSleeping     = unix.CAN_STATE_SLEEPING
	StateMax          = unix.CAN_STATE_MAX

	sizeOfBitTiming        = int(unsafe.Sizeof(BitTiming{}))
	sizeOfBitTimingConst   = int(unsafe.Sizeof(BitTimingConst{}))
	sizeOfClock            = int(unsafe.Sizeof(Clock{}))
	sizeOfCtrlMode         = int(unsafe.Sizeof(CtrlMode{}))
	sizeOfBusErrorCounters = int(unsafe.Sizeof(BusErrorCounters{}))
	sizeOfStats            = int(unsafe.Sizeof(Stats{}))
)

type Device struct {
	ifname string
	index  int32
	li     linkInfoMsg
	ifi    ifInfoMsg
}

// Creates a handle to a CAN device specified by name, e.g. can0.
func New(deviceName string) (*Device, error) {
	iface, err := net.InterfaceByName(deviceName)
	if err != nil {
		return nil, err
	}
	d := &Device{
		index: int32(iface.Index),
	}
	if err := d.updateInfo(); err != nil {
		return nil, err
	}
	return d, nil
}

func (d *Device) IsUp() (bool, error) {
	if err := d.updateInfo(); err != nil {
		return false, err
	}
	if d.ifi.Flags&unix.IFF_UP != 0 {
		return true, nil
	}
	return false, nil
}

// Corresponds to "ip link set up".
func (d *Device) SetUp() error {
	c, err := netlink.Dial(unix.NETLINK_ROUTE, &netlink.Config{})
	if err != nil {
		return fmt.Errorf("couldn't dial netlink socket: %w", err)
	}
	defer c.Close()

	ifi := &ifInfoMsg{
		unix.IfInfomsg{
			Index:  d.index,
			Flags:  unix.IFF_UP,
			Change: unix.IFF_UP,
		},
	}
	req, err := d.newRequest(unix.RTM_NEWLINK, ifi)
	if err != nil {
		return fmt.Errorf("couldn't create netlink request: %w", err)
	}

	res, err := c.Execute(req)
	if err != nil {
		return fmt.Errorf("couldn't set link up: %w", err)
	}
	if len(res) > 1 {
		return fmt.Errorf("expected 1 message, got %d", len(res))
	}
	return nil
}

// Corresponds to "ip link set down".
func (d *Device) SetDown() error {
	c, err := netlink.Dial(unix.NETLINK_ROUTE, &netlink.Config{})
	if err != nil {
		return fmt.Errorf("couldn't dial netlink socket: %w", err)
	}
	defer c.Close()

	ifi := &ifInfoMsg{
		unix.IfInfomsg{
			Index:  d.index,
			Flags:  0,
			Change: unix.IFF_UP,
		},
	}
	req, err := d.newRequest(unix.RTM_NEWLINK, ifi)
	if err != nil {
		return fmt.Errorf("couldn't create netlink request: %w", err)
	}

	res, err := c.Execute(req)
	if err != nil {
		return fmt.Errorf("couldn't set link down: %w", err)
	}
	if len(res) > 1 {
		return fmt.Errorf("expected 1 message, got %d", len(res))
	}
	return nil
}

func (d *Device) Bitrate() (uint32, error) {
	if err := d.updateInfo(); err != nil {
		return 0, fmt.Errorf("couldn't retrieve bitrate: %w", err)
	}
	return d.li.info.BitTiming.Bitrate, nil
}

func (d *Device) SetListenOnlyMode(mode bool) error {
	c, err := netlink.Dial(unix.NETLINK_ROUTE, &netlink.Config{})
	if err != nil {
		return fmt.Errorf("couldn't dial netlink socket: %w", err)
	}
	defer c.Close()

	ifi := &ifInfoMsg{
		unix.IfInfomsg{Index: d.index},
	}
	req, err := d.newRequest(unix.RTM_NEWLINK, ifi)
	if err != nil {
		return fmt.Errorf("couldn't create netlink request: %w", err)
	}

	li := &linkInfoMsg{
		linkType: CanLinkType,
	}

	li.info, err = d.getCurrentParametersForSet()
	if err != nil {
		return fmt.Errorf("couldn't get current parameters: %w", err)
	}

	if mode {
		li.info.CtrlMode.Mask |= unix.CAN_CTRLMODE_LISTENONLY
		li.info.CtrlMode.Flags |= unix.CAN_CTRLMODE_LISTENONLY
	} else {
		li.info.CtrlMode.Mask |= unix.CAN_CTRLMODE_LISTENONLY
		li.info.CtrlMode.Flags = 0
	}

	ae := netlink.NewAttributeEncoder()
	ae.Nested(unix.IFLA_LINKINFO, li.encode)
	liData, err := ae.Encode()
	if err != nil {
		return fmt.Errorf("couldn't encode message: %w", err)
	}

	req.Data = append(req.Data, liData...)

	res, err := c.Execute(req)
	if err != nil {
		return fmt.Errorf("couldn't set listen-only mode: %w", err)
	}
	if len(res) > 1 {
		return fmt.Errorf("expected 1 message, got %d", len(res))
	}
	return nil
}

func (d *Device) SetBitrate(bitrate uint32) error {
	c, err := netlink.Dial(unix.NETLINK_ROUTE, &netlink.Config{})
	if err != nil {
		return fmt.Errorf("couldn't dial netlink socket: %w", err)
	}
	defer c.Close()

	ifi := &ifInfoMsg{
		unix.IfInfomsg{Index: d.index},
	}
	req, err := d.newRequest(unix.RTM_NEWLINK, ifi)
	if err != nil {
		return fmt.Errorf("couldn't create netlink request: %w", err)
	}

	li := &linkInfoMsg{
		linkType: CanLinkType,
	}

	li.info, err = d.getCurrentParametersForSet()
	if err != nil {
		return fmt.Errorf("couldn't get current parameters: %w", err)
	}

	li.info.BitTiming.Bitrate = bitrate
	ae := netlink.NewAttributeEncoder()
	ae.Nested(unix.IFLA_LINKINFO, li.encode)
	liData, err := ae.Encode()
	if err != nil {
		return fmt.Errorf("couldn't encode message: %w", err)
	}
	req.Data = append(req.Data, liData...)

	res, err := c.Execute(req)
	if err != nil {
		return fmt.Errorf("couldn't set bitrate: %w", err)
	}
	if len(res) > 1 {
		return fmt.Errorf("expected 1 message, got %d", len(res))
	}
	return nil
}

type Info struct {
	BitTiming        BitTiming
	BitTimingConst   BitTimingConst
	Clock            Clock
	CtrlMode         CtrlMode
	BusErrorCounters BusErrorCounters
	Type             string

	State     uint32
	RestartMs uint32
}

func (d *Device) Info() (Info, error) {
	if err := d.updateInfo(); err != nil {
		return Info{}, err
	}
	return d.li.info, nil
}

func (d *Device) updateInfo() error {
	c, err := netlink.Dial(unix.NETLINK_ROUTE, &netlink.Config{})
	if err != nil {
		return fmt.Errorf("couldn't dial netlink socket: %w", err)
	}
	defer c.Close()

	ifi := &ifInfoMsg{
		unix.IfInfomsg{Index: d.index},
	}
	req, err := d.newRequest(unix.RTM_GETLINK, ifi)
	if err != nil {
		return fmt.Errorf("couldn't create netlink request: %w", err)
	}

	res, err := c.Execute(req)
	if err != nil {
		return fmt.Errorf("couldn't retrieve link info: %w", err)
	}
	if len(res) > 1 {
		return fmt.Errorf("expected 1 message, got %d", len(res))
	}

	if err := d.unmarshalBinary(res[0].Data); err != nil {
		return fmt.Errorf("couldn't decode info: %w", err)
	}
	return nil
}

func (d *Device) getCurrentParametersForSet() (Info, error) {
	i, err := d.Info()
	if err != nil {
		return Info{}, err
	}

	return Info{BitTiming: BitTiming{unix.CANBitTiming{Bitrate: i.BitTiming.Bitrate}}, CtrlMode: i.CtrlMode}, nil
}

func (d *Device) newRequest(typ netlink.HeaderType, ifi *ifInfoMsg) (netlink.Message, error) {
	req := netlink.Message{
		Header: netlink.Header{
			Flags: netlink.Request | netlink.Acknowledge,
			Type:  typ,
		},
	}
	msg := ifi.marshalBinary()
	req.Data = append(req.Data, msg...)
	return req, nil
}

func (d *Device) unmarshalBinary(data []byte) error {
	if err := d.ifi.unmarshalBinary(data[:unix.SizeofIfInfomsg]); err != nil {
		return fmt.Errorf("couldn't unmarshal ifInfoMsg: %w", err)
	}

	ad, err := netlink.NewAttributeDecoder(data[unix.SizeofIfInfomsg:])
	if err != nil {
		return err
	}
	if d.ifi.Type != unix.ARPHRD_CAN {
		return fmt.Errorf("not a CAN interface")
	}
	for ad.Next() {
		switch ad.Type() {
		case unix.IFLA_IFNAME:
			d.ifname = ad.String()
		case unix.IFLA_LINKINFO:
			ad.Nested(d.li.decode)
			d.li.info.Type = d.li.linkType
		default:
		}
	}
	if err := ad.Err(); err != nil {
		return fmt.Errorf("couldn't decode link: %w", err)
	}
	return nil
}

type ifInfoMsg struct {
	unix.IfInfomsg
}

func (ifi *ifInfoMsg) marshalBinary() []byte {
	buf := make([]byte, unix.SizeofIfInfomsg)
	buf[0] = ifi.Family
	buf[1] = 0 // reserved
	nlenc.PutUint16(buf[2:4], ifi.Type)
	nlenc.PutInt32(buf[4:8], ifi.Index)
	nlenc.PutUint32(buf[8:12], ifi.Flags)
	nlenc.PutUint32(buf[12:16], ifi.Change)
	return buf
}

func (ifi *ifInfoMsg) unmarshalBinary(data []byte) error {
	if len(data) != unix.SizeofIfInfomsg {
		return fmt.Errorf(
			"data is not a valid ifInfoMsg, expected: %d bytes, got: %d bytes",
			unix.SizeofIfInfomsg,
			len(data),
		)
	}
	ifi.Family = nlenc.Uint8(data[0:1])
	ifi.Type = nlenc.Uint16(data[2:4])
	ifi.Index = nlenc.Int32(data[4:8])
	ifi.Flags = nlenc.Uint32(data[8:12])
	ifi.Change = nlenc.Uint32(data[12:16])
	return nil
}

type BitTiming struct {
	unix.CANBitTiming
}

func (bt *BitTiming) marshalBinary() []byte {
	buf := make([]byte, sizeOfBitTiming)
	nlenc.PutUint32(buf[0:4], bt.Bitrate)
	nlenc.PutUint32(buf[4:8], bt.Sample_point)
	nlenc.PutUint32(buf[8:12], bt.Tq)
	nlenc.PutUint32(buf[12:16], bt.Prop_seg)
	nlenc.PutUint32(buf[16:20], bt.Phase_seg1)
	nlenc.PutUint32(buf[20:24], bt.Phase_seg2)
	nlenc.PutUint32(buf[24:28], bt.Sjw)
	nlenc.PutUint32(buf[28:32], bt.Brp)
	return buf
}

func (bt *BitTiming) unmarshalBinary(data []byte) error {
	if len(data) != sizeOfBitTiming {
		return fmt.Errorf(
			"data is not a valid BitTiming, expected: %d bytes, got: %d bytes",
			sizeOfBitTiming,
			len(data),
		)
	}
	bt.Bitrate = nlenc.Uint32(data[0:4])
	bt.Sample_point = nlenc.Uint32(data[4:8])
	bt.Tq = nlenc.Uint32(data[8:12])
	bt.Prop_seg = nlenc.Uint32(data[12:16])
	bt.Phase_seg1 = nlenc.Uint32(data[16:20])
	bt.Phase_seg2 = nlenc.Uint32(data[20:24])
	bt.Sjw = nlenc.Uint32(data[24:28])
	bt.Brp = nlenc.Uint32(data[28:32])
	return nil
}

type BitTimingConst struct {
	unix.CANBitTimingConst
}

func (btc *BitTimingConst) unmarshalBinary(data []byte) error {
	if len(data) != sizeOfBitTimingConst {
		return fmt.Errorf(
			"data is not a valid BitTimingConst, expected: %d bytes, got: %d bytes",
			sizeOfBitTimingConst,
			len(data),
		)
	}
	copy(btc.Name[:], data[0:16])
	btc.Tseg1_min = nlenc.Uint32(data[16:20])
	btc.Tseg1_max = nlenc.Uint32(data[20:24])
	btc.Tseg2_min = nlenc.Uint32(data[24:28])
	btc.Tseg2_max = nlenc.Uint32(data[28:32])
	btc.Sjw_max = nlenc.Uint32(data[32:36])
	btc.Brp_min = nlenc.Uint32(data[36:40])
	btc.Brp_max = nlenc.Uint32(data[40:44])
	btc.Brp_inc = nlenc.Uint32(data[44:48])
	return nil
}

type Clock struct {
	unix.CANClock
}

func (c *Clock) unmarshalBinary(data []byte) error {
	if len(data) != sizeOfClock {
		return fmt.Errorf(
			"data is not a valid Clock, expected: %d bytes, got: %d bytes",
			sizeOfClock,
			len(data),
		)
	}
	c.Freq = nlenc.Uint32(data)
	return nil
}

type CtrlMode struct {
	unix.CANCtrlMode
}

func (cm *CtrlMode) marshalBinary() []byte {
	buf := make([]byte, sizeOfCtrlMode)
	nlenc.PutUint32(buf[0:4], cm.Mask)
	nlenc.PutUint32(buf[4:8], cm.Flags)
	return buf
}

func (cm *CtrlMode) unmarshalBinary(data []byte) error {
	if len(data) != sizeOfCtrlMode {
		return fmt.Errorf(
			"data is not a valid CtrlMode, expected: %d bytes, got: %d bytes",
			sizeOfCtrlMode,
			len(data),
		)
	}
	cm.Mask = nlenc.Uint32(data[0:4])
	cm.Flags = nlenc.Uint32(data[4:8])
	return nil
}

type BusErrorCounters struct {
	unix.CANBusErrorCounters
}

func (bec *BusErrorCounters) unmarshalBinary(data []byte) error {
	if len(data) != sizeOfBusErrorCounters {
		return fmt.Errorf(
			"data is not a valid BusErrorCounters, expected: %d bytes, got: %d bytes",
			sizeOfBusErrorCounters,
			len(data),
		)
	}
	bec.Txerr = nlenc.Uint16(data[0:2])
	bec.Rxerr = nlenc.Uint16(data[2:4])
	return nil
}

type Stats struct {
	unix.CANDeviceStats
}

func (s *Stats) unmarshalBinary(data []byte) error {
	if len(data) != sizeOfStats {
		return fmt.Errorf(
			"data is not a valid Stats, expected: %d bytes, got: %d bytes",
			sizeOfStats,
			len(data),
		)
	}
	s.Bus_error = nlenc.Uint32(data[0:4])
	s.Error_warning = nlenc.Uint32(data[4:8])
	s.Error_passive = nlenc.Uint32(data[8:12])
	s.Bus_off = nlenc.Uint32(data[12:16])
	s.Arbitration_lost = nlenc.Uint32(data[16:20])
	s.Restarts = nlenc.Uint32(data[20:24])
	return nil
}

func (i *Info) decode(nad *netlink.AttributeDecoder) error {
	var err error
	for nad.Next() {
		switch nad.Type() {
		case unix.IFLA_CAN_BITTIMING:
			err = i.BitTiming.unmarshalBinary(nad.Bytes())
		case unix.IFLA_CAN_BITTIMING_CONST:
			err = i.BitTimingConst.unmarshalBinary(nad.Bytes())
		case unix.IFLA_CAN_CLOCK:
			err = i.Clock.unmarshalBinary(nad.Bytes())
		case unix.IFLA_CAN_CTRLMODE:
			err = i.CtrlMode.unmarshalBinary(nad.Bytes())
		case unix.IFLA_CAN_BERR_COUNTER:
			err = i.BusErrorCounters.unmarshalBinary(nad.Bytes())
		default:
		}
		if err != nil {
			return err
		}
	}
	return nil
}

// TODO: add more structures as needed.
func (i *Info) encode(nae *netlink.AttributeEncoder) error {
	nae.Bytes(unix.IFLA_CAN_BITTIMING, i.BitTiming.marshalBinary())
	nae.Bytes(unix.IFLA_CAN_CTRLMODE, i.CtrlMode.marshalBinary())
	return nil
}

type linkInfoMsg struct {
	linkType string
	info     Info
	stats    Stats
}

func (li *linkInfoMsg) decode(nad *netlink.AttributeDecoder) error {
	var err error
	for nad.Next() {
		switch nad.Type() {
		case unix.IFLA_INFO_KIND:
			li.linkType = nad.String()
			if (li.linkType != CanLinkType) && (li.linkType != VcanLinkType) {
				return fmt.Errorf("not a CAN interface")
			}
		case unix.IFLA_INFO_DATA:
			nad.Nested(li.info.decode)
		case unix.IFLA_INFO_XSTATS:
			err = li.stats.unmarshalBinary(nad.Bytes())
		default:
		}
		if err != nil {
			return err
		}
	}

	return nil
}

func (li *linkInfoMsg) encode(nae *netlink.AttributeEncoder) error {
	nae.String(unix.IFLA_INFO_KIND, li.linkType)
	nae.Nested(unix.IFLA_INFO_DATA, li.info.encode)
	return nil
}
