// Package examplecan provides primitives for encoding and decoding example CAN messages.
//
// Source: testdata/dbc/example/example.dbc
package examplecan

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"sync"
	"time"

	"go.einride.tech/can"
	"go.einride.tech/can/pkg/candebug"
	"go.einride.tech/can/pkg/canrunner"
	"go.einride.tech/can/pkg/cantext"
	"go.einride.tech/can/pkg/descriptor"
	"go.einride.tech/can/pkg/generated"
	"go.einride.tech/can/pkg/socketcan"
)

// prevent unused imports
var (
	_ = context.Background
	_ = fmt.Print
	_ = net.Dial
	_ = http.Error
	_ = sync.Mutex{}
	_ = time.Now
	_ = socketcan.Dial
	_ = candebug.ServeMessagesHTTP
	_ = canrunner.Run
)

// Generated code. DO NOT EDIT.
// EmptyMessageReader provides read access to a EmptyMessage message.
type EmptyMessageReader interface {
}

// EmptyMessageWriter provides write access to a EmptyMessage message.
type EmptyMessageWriter interface {
	// CopyFrom copies all values from EmptyMessage.
	CopyFrom(*EmptyMessage) *EmptyMessage
}

type EmptyMessage struct {
}

func NewEmptyMessage() *EmptyMessage {
	m := &EmptyMessage{}
	m.Reset()
	return m
}

func (m *EmptyMessage) Reset() {
}

func (m *EmptyMessage) CopyFrom(o *EmptyMessage) *EmptyMessage {
	_ = m.UnmarshalFrame(o.Frame())
	return m
}

// Descriptor returns the EmptyMessage descriptor.
func (m *EmptyMessage) Descriptor() *descriptor.Message {
	return Messages().EmptyMessage.Message
}

// String returns a compact string representation of the message.
func (m *EmptyMessage) String() string {
	return cantext.MessageString(m)
}

// Frame returns a CAN frame representing the message.
func (m *EmptyMessage) Frame() can.Frame {
	md := Messages().EmptyMessage
	f := can.Frame{ID: md.ID, IsExtended: md.IsExtended, Length: md.Length}
	return f
}

// MarshalFrame encodes the message as a CAN frame.
func (m *EmptyMessage) MarshalFrame() (can.Frame, error) {
	return m.Frame(), nil
}

// UnmarshalFrame decodes the message from a CAN frame.
func (m *EmptyMessage) UnmarshalFrame(f can.Frame) error {
	md := Messages().EmptyMessage
	switch {
	case f.ID != md.ID:
		return fmt.Errorf(
			"unmarshal EmptyMessage: expects ID 1 (got %s with ID %d)", f.String(), f.ID,
		)
	case f.Length != md.Length:
		return fmt.Errorf(
			"unmarshal EmptyMessage: expects length 0 (got %s with length %d)", f.String(), f.Length,
		)
	case f.IsRemote:
		return fmt.Errorf(
			"unmarshal EmptyMessage: expects non-remote frame (got remote frame %s)", f.String(),
		)
	case f.IsExtended != md.IsExtended:
		return fmt.Errorf(
			"unmarshal EmptyMessage: expects standard ID (got %s with extended ID)", f.String(),
		)
	}
	return nil
}

// DriverHeartbeatReader provides read access to a DriverHeartbeat message.
type DriverHeartbeatReader interface {
	// Command returns the value of the Command signal.
	Command() DriverHeartbeat_Command
}

// DriverHeartbeatWriter provides write access to a DriverHeartbeat message.
type DriverHeartbeatWriter interface {
	// CopyFrom copies all values from DriverHeartbeat.
	CopyFrom(*DriverHeartbeat) *DriverHeartbeat
	// SetCommand sets the value of the Command signal.
	SetCommand(DriverHeartbeat_Command) *DriverHeartbeat
}

type DriverHeartbeat struct {
	xxx_Command DriverHeartbeat_Command
}

func NewDriverHeartbeat() *DriverHeartbeat {
	m := &DriverHeartbeat{}
	m.Reset()
	return m
}

func (m *DriverHeartbeat) Reset() {
	m.xxx_Command = 0
}

func (m *DriverHeartbeat) CopyFrom(o *DriverHeartbeat) *DriverHeartbeat {
	_ = m.UnmarshalFrame(o.Frame())
	return m
}

// Descriptor returns the DriverHeartbeat descriptor.
func (m *DriverHeartbeat) Descriptor() *descriptor.Message {
	return Messages().DriverHeartbeat.Message
}

// String returns a compact string representation of the message.
func (m *DriverHeartbeat) String() string {
	return cantext.MessageString(m)
}

func (m *DriverHeartbeat) Command() DriverHeartbeat_Command {
	return m.xxx_Command
}

func (m *DriverHeartbeat) SetCommand(v DriverHeartbeat_Command) *DriverHeartbeat {
	m.xxx_Command = DriverHeartbeat_Command(Messages().DriverHeartbeat.Command.SaturatedCastUnsigned(uint64(v)))
	return m
}

// DriverHeartbeat_Command models the Command signal of the DriverHeartbeat message.
type DriverHeartbeat_Command uint8

// Value descriptions for the Command signal of the DriverHeartbeat message.
const (
	DriverHeartbeat_Command_None         DriverHeartbeat_Command = 0
	DriverHeartbeat_Command_Sync         DriverHeartbeat_Command = 1
	DriverHeartbeat_Command_Reboot       DriverHeartbeat_Command = 2
	DriverHeartbeat_Command_HeadlightsOn DriverHeartbeat_Command = 3
)

func (v DriverHeartbeat_Command) String() string {
	switch v {
	case 0:
		return "None"
	case 1:
		return "Sync"
	case 2:
		return "Reboot"
	case 3:
		return "Headlights On"
	default:
		return fmt.Sprintf("DriverHeartbeat_Command(%d)", v)
	}
}

// Frame returns a CAN frame representing the message.
func (m *DriverHeartbeat) Frame() can.Frame {
	md := Messages().DriverHeartbeat
	f := can.Frame{ID: md.ID, IsExtended: md.IsExtended, Length: md.Length}
	md.Command.MarshalUnsigned(&f.Data, uint64(m.xxx_Command))
	return f
}

// MarshalFrame encodes the message as a CAN frame.
func (m *DriverHeartbeat) MarshalFrame() (can.Frame, error) {
	return m.Frame(), nil
}

// UnmarshalFrame decodes the message from a CAN frame.
func (m *DriverHeartbeat) UnmarshalFrame(f can.Frame) error {
	md := Messages().DriverHeartbeat
	switch {
	case f.ID != md.ID:
		return fmt.Errorf(
			"unmarshal DriverHeartbeat: expects ID 100 (got %s with ID %d)", f.String(), f.ID,
		)
	case f.Length != md.Length:
		return fmt.Errorf(
			"unmarshal DriverHeartbeat: expects length 1 (got %s with length %d)", f.String(), f.Length,
		)
	case f.IsRemote:
		return fmt.Errorf(
			"unmarshal DriverHeartbeat: expects non-remote frame (got remote frame %s)", f.String(),
		)
	case f.IsExtended != md.IsExtended:
		return fmt.Errorf(
			"unmarshal DriverHeartbeat: expects standard ID (got %s with extended ID)", f.String(),
		)
	}
	m.xxx_Command = DriverHeartbeat_Command(md.Command.UnmarshalUnsigned(f.Data))
	return nil
}

// MotorCommandReader provides read access to a MotorCommand message.
type MotorCommandReader interface {
	// Steer returns the physical value of the Steer signal.
	Steer() float64
	// RawSteer returns the raw (encoded) value of the Steer signal.
	RawSteer() int8
	// Drive returns the physical value of the Drive signal.
	Drive() float64
	// RawDrive returns the raw (encoded) value of the Drive signal.
	RawDrive() uint8
}

// MotorCommandWriter provides write access to a MotorCommand message.
type MotorCommandWriter interface {
	// CopyFrom copies all values from MotorCommand.
	CopyFrom(*MotorCommand) *MotorCommand
	// SetSteer sets the physical value of the Steer signal.
	SetSteer(float64) *MotorCommand
	// SetRawSteer sets the raw (encoded) value of the Steer signal.
	SetRawSteer(int8) *MotorCommand
	// SetDrive sets the physical value of the Drive signal.
	SetDrive(float64) *MotorCommand
	// SetRawDrive sets the raw (encoded) value of the Drive signal.
	SetRawDrive(uint8) *MotorCommand
}

type MotorCommand struct {
	xxx_Steer int8
	xxx_Drive uint8
}

func NewMotorCommand() *MotorCommand {
	m := &MotorCommand{}
	m.Reset()
	return m
}

func (m *MotorCommand) Reset() {
	m.xxx_Steer = 0
	m.xxx_Drive = 0
}

func (m *MotorCommand) CopyFrom(o *MotorCommand) *MotorCommand {
	_ = m.UnmarshalFrame(o.Frame())
	return m
}

// Descriptor returns the MotorCommand descriptor.
func (m *MotorCommand) Descriptor() *descriptor.Message {
	return Messages().MotorCommand.Message
}

// String returns a compact string representation of the message.
func (m *MotorCommand) String() string {
	return cantext.MessageString(m)
}

func (m *MotorCommand) Steer() float64 {
	return Messages().MotorCommand.Steer.ToPhysical(float64(m.xxx_Steer))
}

func (m *MotorCommand) SetSteer(v float64) *MotorCommand {
	m.xxx_Steer = int8(Messages().MotorCommand.Steer.FromPhysical(v))
	return m
}

func (m *MotorCommand) RawSteer() int8 {
	return m.xxx_Steer
}

func (m *MotorCommand) SetRawSteer(v int8) *MotorCommand {
	m.xxx_Steer = int8(Messages().MotorCommand.Steer.SaturatedCastSigned(int64(v)))
	return m
}

func (m *MotorCommand) Drive() float64 {
	return Messages().MotorCommand.Drive.ToPhysical(float64(m.xxx_Drive))
}

func (m *MotorCommand) SetDrive(v float64) *MotorCommand {
	m.xxx_Drive = uint8(Messages().MotorCommand.Drive.FromPhysical(v))
	return m
}

func (m *MotorCommand) RawDrive() uint8 {
	return m.xxx_Drive
}

func (m *MotorCommand) SetRawDrive(v uint8) *MotorCommand {
	m.xxx_Drive = uint8(Messages().MotorCommand.Drive.SaturatedCastUnsigned(uint64(v)))
	return m
}

// Frame returns a CAN frame representing the message.
func (m *MotorCommand) Frame() can.Frame {
	md := Messages().MotorCommand
	f := can.Frame{ID: md.ID, IsExtended: md.IsExtended, Length: md.Length}
	md.Steer.MarshalSigned(&f.Data, int64(m.xxx_Steer))
	md.Drive.MarshalUnsigned(&f.Data, uint64(m.xxx_Drive))
	return f
}

// MarshalFrame encodes the message as a CAN frame.
func (m *MotorCommand) MarshalFrame() (can.Frame, error) {
	return m.Frame(), nil
}

// UnmarshalFrame decodes the message from a CAN frame.
func (m *MotorCommand) UnmarshalFrame(f can.Frame) error {
	md := Messages().MotorCommand
	switch {
	case f.ID != md.ID:
		return fmt.Errorf(
			"unmarshal MotorCommand: expects ID 101 (got %s with ID %d)", f.String(), f.ID,
		)
	case f.Length != md.Length:
		return fmt.Errorf(
			"unmarshal MotorCommand: expects length 1 (got %s with length %d)", f.String(), f.Length,
		)
	case f.IsRemote:
		return fmt.Errorf(
			"unmarshal MotorCommand: expects non-remote frame (got remote frame %s)", f.String(),
		)
	case f.IsExtended != md.IsExtended:
		return fmt.Errorf(
			"unmarshal MotorCommand: expects standard ID (got %s with extended ID)", f.String(),
		)
	}
	m.xxx_Steer = int8(md.Steer.UnmarshalSigned(f.Data))
	m.xxx_Drive = uint8(md.Drive.UnmarshalUnsigned(f.Data))
	return nil
}

// SensorSonarsReader provides read access to a SensorSonars message.
type SensorSonarsReader interface {
	// Mux returns the value of the Mux signal.
	Mux() uint8
	// ErrCount returns the value of the ErrCount signal.
	ErrCount() uint16
	// Left returns the physical value of the Left signal.
	Left() float64
	// RawLeft returns the raw (encoded) value of the Left signal.
	RawLeft() uint16
	// NoFiltLeft returns the physical value of the NoFiltLeft signal.
	NoFiltLeft() float64
	// RawNoFiltLeft returns the raw (encoded) value of the NoFiltLeft signal.
	RawNoFiltLeft() uint16
	// Middle returns the physical value of the Middle signal.
	Middle() float64
	// RawMiddle returns the raw (encoded) value of the Middle signal.
	RawMiddle() uint16
	// NoFiltMiddle returns the physical value of the NoFiltMiddle signal.
	NoFiltMiddle() float64
	// RawNoFiltMiddle returns the raw (encoded) value of the NoFiltMiddle signal.
	RawNoFiltMiddle() uint16
	// Right returns the physical value of the Right signal.
	Right() float64
	// RawRight returns the raw (encoded) value of the Right signal.
	RawRight() uint16
	// NoFiltRight returns the physical value of the NoFiltRight signal.
	NoFiltRight() float64
	// RawNoFiltRight returns the raw (encoded) value of the NoFiltRight signal.
	RawNoFiltRight() uint16
	// Rear returns the physical value of the Rear signal.
	Rear() float64
	// RawRear returns the raw (encoded) value of the Rear signal.
	RawRear() uint16
	// NoFiltRear returns the physical value of the NoFiltRear signal.
	NoFiltRear() float64
	// RawNoFiltRear returns the raw (encoded) value of the NoFiltRear signal.
	RawNoFiltRear() uint16
}

// SensorSonarsWriter provides write access to a SensorSonars message.
type SensorSonarsWriter interface {
	// CopyFrom copies all values from SensorSonars.
	CopyFrom(*SensorSonars) *SensorSonars
	// SetMux sets the value of the Mux signal.
	SetMux(uint8) *SensorSonars
	// SetErrCount sets the value of the ErrCount signal.
	SetErrCount(uint16) *SensorSonars
	// SetLeft sets the physical value of the Left signal.
	SetLeft(float64) *SensorSonars
	// SetRawLeft sets the raw (encoded) value of the Left signal.
	SetRawLeft(uint16) *SensorSonars
	// SetNoFiltLeft sets the physical value of the NoFiltLeft signal.
	SetNoFiltLeft(float64) *SensorSonars
	// SetRawNoFiltLeft sets the raw (encoded) value of the NoFiltLeft signal.
	SetRawNoFiltLeft(uint16) *SensorSonars
	// SetMiddle sets the physical value of the Middle signal.
	SetMiddle(float64) *SensorSonars
	// SetRawMiddle sets the raw (encoded) value of the Middle signal.
	SetRawMiddle(uint16) *SensorSonars
	// SetNoFiltMiddle sets the physical value of the NoFiltMiddle signal.
	SetNoFiltMiddle(float64) *SensorSonars
	// SetRawNoFiltMiddle sets the raw (encoded) value of the NoFiltMiddle signal.
	SetRawNoFiltMiddle(uint16) *SensorSonars
	// SetRight sets the physical value of the Right signal.
	SetRight(float64) *SensorSonars
	// SetRawRight sets the raw (encoded) value of the Right signal.
	SetRawRight(uint16) *SensorSonars
	// SetNoFiltRight sets the physical value of the NoFiltRight signal.
	SetNoFiltRight(float64) *SensorSonars
	// SetRawNoFiltRight sets the raw (encoded) value of the NoFiltRight signal.
	SetRawNoFiltRight(uint16) *SensorSonars
	// SetRear sets the physical value of the Rear signal.
	SetRear(float64) *SensorSonars
	// SetRawRear sets the raw (encoded) value of the Rear signal.
	SetRawRear(uint16) *SensorSonars
	// SetNoFiltRear sets the physical value of the NoFiltRear signal.
	SetNoFiltRear(float64) *SensorSonars
	// SetRawNoFiltRear sets the raw (encoded) value of the NoFiltRear signal.
	SetRawNoFiltRear(uint16) *SensorSonars
}

type SensorSonars struct {
	xxx_Mux          uint8
	xxx_ErrCount     uint16
	xxx_Left         uint16
	xxx_NoFiltLeft   uint16
	xxx_Middle       uint16
	xxx_NoFiltMiddle uint16
	xxx_Right        uint16
	xxx_NoFiltRight  uint16
	xxx_Rear         uint16
	xxx_NoFiltRear   uint16
}

func NewSensorSonars() *SensorSonars {
	m := &SensorSonars{}
	m.Reset()
	return m
}

func (m *SensorSonars) Reset() {
	m.xxx_Mux = 0
	m.xxx_ErrCount = 0
	m.xxx_Left = 0
	m.xxx_NoFiltLeft = 0
	m.xxx_Middle = 0
	m.xxx_NoFiltMiddle = 0
	m.xxx_Right = 0
	m.xxx_NoFiltRight = 0
	m.xxx_Rear = 0
	m.xxx_NoFiltRear = 0
}

func (m *SensorSonars) CopyFrom(o *SensorSonars) *SensorSonars {
	_ = m.UnmarshalFrame(o.Frame())
	return m
}

// Descriptor returns the SensorSonars descriptor.
func (m *SensorSonars) Descriptor() *descriptor.Message {
	return Messages().SensorSonars.Message
}

// String returns a compact string representation of the message.
func (m *SensorSonars) String() string {
	return cantext.MessageString(m)
}

func (m *SensorSonars) Mux() uint8 {
	return m.xxx_Mux
}

func (m *SensorSonars) SetMux(v uint8) *SensorSonars {
	m.xxx_Mux = uint8(Messages().SensorSonars.Mux.SaturatedCastUnsigned(uint64(v)))
	return m
}

func (m *SensorSonars) ErrCount() uint16 {
	return m.xxx_ErrCount
}

func (m *SensorSonars) SetErrCount(v uint16) *SensorSonars {
	m.xxx_ErrCount = uint16(Messages().SensorSonars.ErrCount.SaturatedCastUnsigned(uint64(v)))
	return m
}

func (m *SensorSonars) Left() float64 {
	return Messages().SensorSonars.Left.ToPhysical(float64(m.xxx_Left))
}

func (m *SensorSonars) SetLeft(v float64) *SensorSonars {
	m.xxx_Left = uint16(Messages().SensorSonars.Left.FromPhysical(v))
	return m
}

func (m *SensorSonars) RawLeft() uint16 {
	return m.xxx_Left
}

func (m *SensorSonars) SetRawLeft(v uint16) *SensorSonars {
	m.xxx_Left = uint16(Messages().SensorSonars.Left.SaturatedCastUnsigned(uint64(v)))
	return m
}

func (m *SensorSonars) NoFiltLeft() float64 {
	return Messages().SensorSonars.NoFiltLeft.ToPhysical(float64(m.xxx_NoFiltLeft))
}

func (m *SensorSonars) SetNoFiltLeft(v float64) *SensorSonars {
	m.xxx_NoFiltLeft = uint16(Messages().SensorSonars.NoFiltLeft.FromPhysical(v))
	return m
}

func (m *SensorSonars) RawNoFiltLeft() uint16 {
	return m.xxx_NoFiltLeft
}

func (m *SensorSonars) SetRawNoFiltLeft(v uint16) *SensorSonars {
	m.xxx_NoFiltLeft = uint16(Messages().SensorSonars.NoFiltLeft.SaturatedCastUnsigned(uint64(v)))
	return m
}

func (m *SensorSonars) Middle() float64 {
	return Messages().SensorSonars.Middle.ToPhysical(float64(m.xxx_Middle))
}

func (m *SensorSonars) SetMiddle(v float64) *SensorSonars {
	m.xxx_Middle = uint16(Messages().SensorSonars.Middle.FromPhysical(v))
	return m
}

func (m *SensorSonars) RawMiddle() uint16 {
	return m.xxx_Middle
}

func (m *SensorSonars) SetRawMiddle(v uint16) *SensorSonars {
	m.xxx_Middle = uint16(Messages().SensorSonars.Middle.SaturatedCastUnsigned(uint64(v)))
	return m
}

func (m *SensorSonars) NoFiltMiddle() float64 {
	return Messages().SensorSonars.NoFiltMiddle.ToPhysical(float64(m.xxx_NoFiltMiddle))
}

func (m *SensorSonars) SetNoFiltMiddle(v float64) *SensorSonars {
	m.xxx_NoFiltMiddle = uint16(Messages().SensorSonars.NoFiltMiddle.FromPhysical(v))
	return m
}

func (m *SensorSonars) RawNoFiltMiddle() uint16 {
	return m.xxx_NoFiltMiddle
}

func (m *SensorSonars) SetRawNoFiltMiddle(v uint16) *SensorSonars {
	m.xxx_NoFiltMiddle = uint16(Messages().SensorSonars.NoFiltMiddle.SaturatedCastUnsigned(uint64(v)))
	return m
}

func (m *SensorSonars) Right() float64 {
	return Messages().SensorSonars.Right.ToPhysical(float64(m.xxx_Right))
}

func (m *SensorSonars) SetRight(v float64) *SensorSonars {
	m.xxx_Right = uint16(Messages().SensorSonars.Right.FromPhysical(v))
	return m
}

func (m *SensorSonars) RawRight() uint16 {
	return m.xxx_Right
}

func (m *SensorSonars) SetRawRight(v uint16) *SensorSonars {
	m.xxx_Right = uint16(Messages().SensorSonars.Right.SaturatedCastUnsigned(uint64(v)))
	return m
}

func (m *SensorSonars) NoFiltRight() float64 {
	return Messages().SensorSonars.NoFiltRight.ToPhysical(float64(m.xxx_NoFiltRight))
}

func (m *SensorSonars) SetNoFiltRight(v float64) *SensorSonars {
	m.xxx_NoFiltRight = uint16(Messages().SensorSonars.NoFiltRight.FromPhysical(v))
	return m
}

func (m *SensorSonars) RawNoFiltRight() uint16 {
	return m.xxx_NoFiltRight
}

func (m *SensorSonars) SetRawNoFiltRight(v uint16) *SensorSonars {
	m.xxx_NoFiltRight = uint16(Messages().SensorSonars.NoFiltRight.SaturatedCastUnsigned(uint64(v)))
	return m
}

func (m *SensorSonars) Rear() float64 {
	return Messages().SensorSonars.Rear.ToPhysical(float64(m.xxx_Rear))
}

func (m *SensorSonars) SetRear(v float64) *SensorSonars {
	m.xxx_Rear = uint16(Messages().SensorSonars.Rear.FromPhysical(v))
	return m
}

func (m *SensorSonars) RawRear() uint16 {
	return m.xxx_Rear
}

func (m *SensorSonars) SetRawRear(v uint16) *SensorSonars {
	m.xxx_Rear = uint16(Messages().SensorSonars.Rear.SaturatedCastUnsigned(uint64(v)))
	return m
}

func (m *SensorSonars) NoFiltRear() float64 {
	return Messages().SensorSonars.NoFiltRear.ToPhysical(float64(m.xxx_NoFiltRear))
}

func (m *SensorSonars) SetNoFiltRear(v float64) *SensorSonars {
	m.xxx_NoFiltRear = uint16(Messages().SensorSonars.NoFiltRear.FromPhysical(v))
	return m
}

func (m *SensorSonars) RawNoFiltRear() uint16 {
	return m.xxx_NoFiltRear
}

func (m *SensorSonars) SetRawNoFiltRear(v uint16) *SensorSonars {
	m.xxx_NoFiltRear = uint16(Messages().SensorSonars.NoFiltRear.SaturatedCastUnsigned(uint64(v)))
	return m
}

// Frame returns a CAN frame representing the message.
func (m *SensorSonars) Frame() can.Frame {
	md := Messages().SensorSonars
	f := can.Frame{ID: md.ID, IsExtended: md.IsExtended, Length: md.Length}
	md.Mux.MarshalUnsigned(&f.Data, uint64(m.xxx_Mux))
	md.ErrCount.MarshalUnsigned(&f.Data, uint64(m.xxx_ErrCount))
	if m.xxx_Mux == 0 {
		md.Left.MarshalUnsigned(&f.Data, uint64(m.xxx_Left))
	}
	if m.xxx_Mux == 1 {
		md.NoFiltLeft.MarshalUnsigned(&f.Data, uint64(m.xxx_NoFiltLeft))
	}
	if m.xxx_Mux == 0 {
		md.Middle.MarshalUnsigned(&f.Data, uint64(m.xxx_Middle))
	}
	if m.xxx_Mux == 1 {
		md.NoFiltMiddle.MarshalUnsigned(&f.Data, uint64(m.xxx_NoFiltMiddle))
	}
	if m.xxx_Mux == 0 {
		md.Right.MarshalUnsigned(&f.Data, uint64(m.xxx_Right))
	}
	if m.xxx_Mux == 1 {
		md.NoFiltRight.MarshalUnsigned(&f.Data, uint64(m.xxx_NoFiltRight))
	}
	if m.xxx_Mux == 0 {
		md.Rear.MarshalUnsigned(&f.Data, uint64(m.xxx_Rear))
	}
	if m.xxx_Mux == 1 {
		md.NoFiltRear.MarshalUnsigned(&f.Data, uint64(m.xxx_NoFiltRear))
	}
	return f
}

// MarshalFrame encodes the message as a CAN frame.
func (m *SensorSonars) MarshalFrame() (can.Frame, error) {
	return m.Frame(), nil
}

// UnmarshalFrame decodes the message from a CAN frame.
func (m *SensorSonars) UnmarshalFrame(f can.Frame) error {
	md := Messages().SensorSonars
	switch {
	case f.ID != md.ID:
		return fmt.Errorf(
			"unmarshal SensorSonars: expects ID 200 (got %s with ID %d)", f.String(), f.ID,
		)
	case f.Length != md.Length:
		return fmt.Errorf(
			"unmarshal SensorSonars: expects length 8 (got %s with length %d)", f.String(), f.Length,
		)
	case f.IsRemote:
		return fmt.Errorf(
			"unmarshal SensorSonars: expects non-remote frame (got remote frame %s)", f.String(),
		)
	case f.IsExtended != md.IsExtended:
		return fmt.Errorf(
			"unmarshal SensorSonars: expects standard ID (got %s with extended ID)", f.String(),
		)
	}
	m.xxx_Mux = uint8(md.Mux.UnmarshalUnsigned(f.Data))
	m.xxx_ErrCount = uint16(md.ErrCount.UnmarshalUnsigned(f.Data))
	if m.xxx_Mux == 0 {
		m.xxx_Left = uint16(md.Left.UnmarshalUnsigned(f.Data))
	}
	if m.xxx_Mux == 1 {
		m.xxx_NoFiltLeft = uint16(md.NoFiltLeft.UnmarshalUnsigned(f.Data))
	}
	if m.xxx_Mux == 0 {
		m.xxx_Middle = uint16(md.Middle.UnmarshalUnsigned(f.Data))
	}
	if m.xxx_Mux == 1 {
		m.xxx_NoFiltMiddle = uint16(md.NoFiltMiddle.UnmarshalUnsigned(f.Data))
	}
	if m.xxx_Mux == 0 {
		m.xxx_Right = uint16(md.Right.UnmarshalUnsigned(f.Data))
	}
	if m.xxx_Mux == 1 {
		m.xxx_NoFiltRight = uint16(md.NoFiltRight.UnmarshalUnsigned(f.Data))
	}
	if m.xxx_Mux == 0 {
		m.xxx_Rear = uint16(md.Rear.UnmarshalUnsigned(f.Data))
	}
	if m.xxx_Mux == 1 {
		m.xxx_NoFiltRear = uint16(md.NoFiltRear.UnmarshalUnsigned(f.Data))
	}
	return nil
}

// MotorStatusReader provides read access to a MotorStatus message.
type MotorStatusReader interface {
	// WheelError returns the value of the WheelError signal.
	WheelError() bool
	// SpeedKph returns the physical value of the SpeedKph signal.
	SpeedKph() float64
	// RawSpeedKph returns the raw (encoded) value of the SpeedKph signal.
	RawSpeedKph() uint16
}

// MotorStatusWriter provides write access to a MotorStatus message.
type MotorStatusWriter interface {
	// CopyFrom copies all values from MotorStatus.
	CopyFrom(*MotorStatus) *MotorStatus
	// SetWheelError sets the value of the WheelError signal.
	SetWheelError(bool) *MotorStatus
	// SetSpeedKph sets the physical value of the SpeedKph signal.
	SetSpeedKph(float64) *MotorStatus
	// SetRawSpeedKph sets the raw (encoded) value of the SpeedKph signal.
	SetRawSpeedKph(uint16) *MotorStatus
}

type MotorStatus struct {
	xxx_WheelError bool
	xxx_SpeedKph   uint16
}

func NewMotorStatus() *MotorStatus {
	m := &MotorStatus{}
	m.Reset()
	return m
}

func (m *MotorStatus) Reset() {
	m.xxx_WheelError = false
	m.xxx_SpeedKph = 0
}

func (m *MotorStatus) CopyFrom(o *MotorStatus) *MotorStatus {
	_ = m.UnmarshalFrame(o.Frame())
	return m
}

// Descriptor returns the MotorStatus descriptor.
func (m *MotorStatus) Descriptor() *descriptor.Message {
	return Messages().MotorStatus.Message
}

// String returns a compact string representation of the message.
func (m *MotorStatus) String() string {
	return cantext.MessageString(m)
}

func (m *MotorStatus) WheelError() bool {
	return m.xxx_WheelError
}

func (m *MotorStatus) SetWheelError(v bool) *MotorStatus {
	m.xxx_WheelError = v
	return m
}

func (m *MotorStatus) SpeedKph() float64 {
	return Messages().MotorStatus.SpeedKph.ToPhysical(float64(m.xxx_SpeedKph))
}

func (m *MotorStatus) SetSpeedKph(v float64) *MotorStatus {
	m.xxx_SpeedKph = uint16(Messages().MotorStatus.SpeedKph.FromPhysical(v))
	return m
}

func (m *MotorStatus) RawSpeedKph() uint16 {
	return m.xxx_SpeedKph
}

func (m *MotorStatus) SetRawSpeedKph(v uint16) *MotorStatus {
	m.xxx_SpeedKph = uint16(Messages().MotorStatus.SpeedKph.SaturatedCastUnsigned(uint64(v)))
	return m
}

// Frame returns a CAN frame representing the message.
func (m *MotorStatus) Frame() can.Frame {
	md := Messages().MotorStatus
	f := can.Frame{ID: md.ID, IsExtended: md.IsExtended, Length: md.Length}
	md.WheelError.MarshalBool(&f.Data, bool(m.xxx_WheelError))
	md.SpeedKph.MarshalUnsigned(&f.Data, uint64(m.xxx_SpeedKph))
	return f
}

// MarshalFrame encodes the message as a CAN frame.
func (m *MotorStatus) MarshalFrame() (can.Frame, error) {
	return m.Frame(), nil
}

// UnmarshalFrame decodes the message from a CAN frame.
func (m *MotorStatus) UnmarshalFrame(f can.Frame) error {
	md := Messages().MotorStatus
	switch {
	case f.ID != md.ID:
		return fmt.Errorf(
			"unmarshal MotorStatus: expects ID 400 (got %s with ID %d)", f.String(), f.ID,
		)
	case f.Length != md.Length:
		return fmt.Errorf(
			"unmarshal MotorStatus: expects length 3 (got %s with length %d)", f.String(), f.Length,
		)
	case f.IsRemote:
		return fmt.Errorf(
			"unmarshal MotorStatus: expects non-remote frame (got remote frame %s)", f.String(),
		)
	case f.IsExtended != md.IsExtended:
		return fmt.Errorf(
			"unmarshal MotorStatus: expects standard ID (got %s with extended ID)", f.String(),
		)
	}
	m.xxx_WheelError = bool(md.WheelError.UnmarshalBool(f.Data))
	m.xxx_SpeedKph = uint16(md.SpeedKph.UnmarshalUnsigned(f.Data))
	return nil
}

// IODebugReader provides read access to a IODebug message.
type IODebugReader interface {
	// TestUnsigned returns the value of the TestUnsigned signal.
	TestUnsigned() uint8
	// TestEnum returns the value of the TestEnum signal.
	TestEnum() IODebug_TestEnum
	// TestSigned returns the value of the TestSigned signal.
	TestSigned() int8
	// TestFloat returns the physical value of the TestFloat signal.
	TestFloat() float64
	// RawTestFloat returns the raw (encoded) value of the TestFloat signal.
	RawTestFloat() uint8
	// TestBoolEnum returns the value of the TestBoolEnum signal.
	TestBoolEnum() IODebug_TestBoolEnum
	// TestScaledEnum returns the physical value of the TestScaledEnum signal.
	TestScaledEnum() float64
	// RawTestScaledEnum returns the raw (encoded) value of the TestScaledEnum signal.
	RawTestScaledEnum() IODebug_TestScaledEnum
}

// IODebugWriter provides write access to a IODebug message.
type IODebugWriter interface {
	// CopyFrom copies all values from IODebug.
	CopyFrom(*IODebug) *IODebug
	// SetTestUnsigned sets the value of the TestUnsigned signal.
	SetTestUnsigned(uint8) *IODebug
	// SetTestEnum sets the value of the TestEnum signal.
	SetTestEnum(IODebug_TestEnum) *IODebug
	// SetTestSigned sets the value of the TestSigned signal.
	SetTestSigned(int8) *IODebug
	// SetTestFloat sets the physical value of the TestFloat signal.
	SetTestFloat(float64) *IODebug
	// SetRawTestFloat sets the raw (encoded) value of the TestFloat signal.
	SetRawTestFloat(uint8) *IODebug
	// SetTestBoolEnum sets the value of the TestBoolEnum signal.
	SetTestBoolEnum(IODebug_TestBoolEnum) *IODebug
	// SetTestScaledEnum sets the physical value of the TestScaledEnum signal.
	SetTestScaledEnum(float64) *IODebug
	// SetRawTestScaledEnum sets the raw (encoded) value of the TestScaledEnum signal.
	SetRawTestScaledEnum(IODebug_TestScaledEnum) *IODebug
}

type IODebug struct {
	xxx_TestUnsigned   uint8
	xxx_TestEnum       IODebug_TestEnum
	xxx_TestSigned     int8
	xxx_TestFloat      uint8
	xxx_TestBoolEnum   IODebug_TestBoolEnum
	xxx_TestScaledEnum IODebug_TestScaledEnum
}

func NewIODebug() *IODebug {
	m := &IODebug{}
	m.Reset()
	return m
}

func (m *IODebug) Reset() {
	m.xxx_TestUnsigned = 0
	m.xxx_TestEnum = 2
	m.xxx_TestSigned = 0
	m.xxx_TestFloat = 0
	m.xxx_TestBoolEnum = false
	m.xxx_TestScaledEnum = 0
}

func (m *IODebug) CopyFrom(o *IODebug) *IODebug {
	_ = m.UnmarshalFrame(o.Frame())
	return m
}

// Descriptor returns the IODebug descriptor.
func (m *IODebug) Descriptor() *descriptor.Message {
	return Messages().IODebug.Message
}

// String returns a compact string representation of the message.
func (m *IODebug) String() string {
	return cantext.MessageString(m)
}

func (m *IODebug) TestUnsigned() uint8 {
	return m.xxx_TestUnsigned
}

func (m *IODebug) SetTestUnsigned(v uint8) *IODebug {
	m.xxx_TestUnsigned = uint8(Messages().IODebug.TestUnsigned.SaturatedCastUnsigned(uint64(v)))
	return m
}

func (m *IODebug) TestEnum() IODebug_TestEnum {
	return m.xxx_TestEnum
}

func (m *IODebug) SetTestEnum(v IODebug_TestEnum) *IODebug {
	m.xxx_TestEnum = IODebug_TestEnum(Messages().IODebug.TestEnum.SaturatedCastUnsigned(uint64(v)))
	return m
}

func (m *IODebug) TestSigned() int8 {
	return m.xxx_TestSigned
}

func (m *IODebug) SetTestSigned(v int8) *IODebug {
	m.xxx_TestSigned = int8(Messages().IODebug.TestSigned.SaturatedCastSigned(int64(v)))
	return m
}

func (m *IODebug) TestFloat() float64 {
	return Messages().IODebug.TestFloat.ToPhysical(float64(m.xxx_TestFloat))
}

func (m *IODebug) SetTestFloat(v float64) *IODebug {
	m.xxx_TestFloat = uint8(Messages().IODebug.TestFloat.FromPhysical(v))
	return m
}

func (m *IODebug) RawTestFloat() uint8 {
	return m.xxx_TestFloat
}

func (m *IODebug) SetRawTestFloat(v uint8) *IODebug {
	m.xxx_TestFloat = uint8(Messages().IODebug.TestFloat.SaturatedCastUnsigned(uint64(v)))
	return m
}

func (m *IODebug) TestBoolEnum() IODebug_TestBoolEnum {
	return m.xxx_TestBoolEnum
}

func (m *IODebug) SetTestBoolEnum(v IODebug_TestBoolEnum) *IODebug {
	m.xxx_TestBoolEnum = v
	return m
}

func (m *IODebug) TestScaledEnum() float64 {
	return Messages().IODebug.TestScaledEnum.ToPhysical(float64(m.xxx_TestScaledEnum))
}

func (m *IODebug) SetTestScaledEnum(v float64) *IODebug {
	m.xxx_TestScaledEnum = IODebug_TestScaledEnum(Messages().IODebug.TestScaledEnum.FromPhysical(v))
	return m
}

func (m *IODebug) RawTestScaledEnum() IODebug_TestScaledEnum {
	return m.xxx_TestScaledEnum
}

func (m *IODebug) SetRawTestScaledEnum(v IODebug_TestScaledEnum) *IODebug {
	m.xxx_TestScaledEnum = IODebug_TestScaledEnum(Messages().IODebug.TestScaledEnum.SaturatedCastUnsigned(uint64(v)))
	return m
}

// IODebug_TestEnum models the TestEnum signal of the IODebug message.
type IODebug_TestEnum uint8

// Value descriptions for the TestEnum signal of the IODebug message.
const (
	IODebug_TestEnum_One IODebug_TestEnum = 1
	IODebug_TestEnum_Two IODebug_TestEnum = 2
)

func (v IODebug_TestEnum) String() string {
	switch v {
	case 1:
		return "One"
	case 2:
		return "Two"
	default:
		return fmt.Sprintf("IODebug_TestEnum(%d)", v)
	}
}

// IODebug_TestBoolEnum models the TestBoolEnum signal of the IODebug message.
type IODebug_TestBoolEnum bool

// Value descriptions for the TestBoolEnum signal of the IODebug message.
const (
	IODebug_TestBoolEnum_Zero IODebug_TestBoolEnum = false
	IODebug_TestBoolEnum_One  IODebug_TestBoolEnum = true
)

func (v IODebug_TestBoolEnum) String() string {
	switch bool(v) {
	case false:
		return "Zero"
	case true:
		return "One"
	}
	return fmt.Sprintf("IODebug_TestBoolEnum(%t)", v)
}

// IODebug_TestScaledEnum models the TestScaledEnum signal of the IODebug message.
type IODebug_TestScaledEnum uint8

// Value descriptions for the TestScaledEnum signal of the IODebug message.
const (
	IODebug_TestScaledEnum_Zero IODebug_TestScaledEnum = 0
	IODebug_TestScaledEnum_Two  IODebug_TestScaledEnum = 1
	IODebug_TestScaledEnum_Four IODebug_TestScaledEnum = 2
	IODebug_TestScaledEnum_Six  IODebug_TestScaledEnum = 3
)

func (v IODebug_TestScaledEnum) String() string {
	switch v {
	case 0:
		return "Zero"
	case 1:
		return "Two"
	case 2:
		return "Four"
	case 3:
		return "Six"
	default:
		return fmt.Sprintf("IODebug_TestScaledEnum(%d)", v)
	}
}

// Frame returns a CAN frame representing the message.
func (m *IODebug) Frame() can.Frame {
	md := Messages().IODebug
	f := can.Frame{ID: md.ID, IsExtended: md.IsExtended, Length: md.Length}
	md.TestUnsigned.MarshalUnsigned(&f.Data, uint64(m.xxx_TestUnsigned))
	md.TestEnum.MarshalUnsigned(&f.Data, uint64(m.xxx_TestEnum))
	md.TestSigned.MarshalSigned(&f.Data, int64(m.xxx_TestSigned))
	md.TestFloat.MarshalUnsigned(&f.Data, uint64(m.xxx_TestFloat))
	md.TestBoolEnum.MarshalBool(&f.Data, bool(m.xxx_TestBoolEnum))
	md.TestScaledEnum.MarshalUnsigned(&f.Data, uint64(m.xxx_TestScaledEnum))
	return f
}

// MarshalFrame encodes the message as a CAN frame.
func (m *IODebug) MarshalFrame() (can.Frame, error) {
	return m.Frame(), nil
}

// UnmarshalFrame decodes the message from a CAN frame.
func (m *IODebug) UnmarshalFrame(f can.Frame) error {
	md := Messages().IODebug
	switch {
	case f.ID != md.ID:
		return fmt.Errorf(
			"unmarshal IODebug: expects ID 500 (got %s with ID %d)", f.String(), f.ID,
		)
	case f.Length != md.Length:
		return fmt.Errorf(
			"unmarshal IODebug: expects length 6 (got %s with length %d)", f.String(), f.Length,
		)
	case f.IsRemote:
		return fmt.Errorf(
			"unmarshal IODebug: expects non-remote frame (got remote frame %s)", f.String(),
		)
	case f.IsExtended != md.IsExtended:
		return fmt.Errorf(
			"unmarshal IODebug: expects standard ID (got %s with extended ID)", f.String(),
		)
	}
	m.xxx_TestUnsigned = uint8(md.TestUnsigned.UnmarshalUnsigned(f.Data))
	m.xxx_TestEnum = IODebug_TestEnum(md.TestEnum.UnmarshalUnsigned(f.Data))
	m.xxx_TestSigned = int8(md.TestSigned.UnmarshalSigned(f.Data))
	m.xxx_TestFloat = uint8(md.TestFloat.UnmarshalUnsigned(f.Data))
	m.xxx_TestBoolEnum = IODebug_TestBoolEnum(md.TestBoolEnum.UnmarshalBool(f.Data))
	m.xxx_TestScaledEnum = IODebug_TestScaledEnum(md.TestScaledEnum.UnmarshalUnsigned(f.Data))
	return nil
}

// IOFloat32Reader provides read access to a IOFloat32 message.
type IOFloat32Reader interface {
	// Float32ValueNoRange returns the value of the Float32ValueNoRange signal.
	Float32ValueNoRange() float32
	// Float32WithRange returns the physical value of the Float32WithRange signal.
	Float32WithRange() float64
	// RawFloat32WithRange returns the raw (encoded) value of the Float32WithRange signal.
	RawFloat32WithRange() float32
}

// IOFloat32Writer provides write access to a IOFloat32 message.
type IOFloat32Writer interface {
	// CopyFrom copies all values from IOFloat32.
	CopyFrom(*IOFloat32) *IOFloat32
	// SetFloat32ValueNoRange sets the value of the Float32ValueNoRange signal.
	SetFloat32ValueNoRange(float32) *IOFloat32
	// SetFloat32WithRange sets the physical value of the Float32WithRange signal.
	SetFloat32WithRange(float64) *IOFloat32
	// SetRawFloat32WithRange sets the raw (encoded) value of the Float32WithRange signal.
	SetRawFloat32WithRange(float32) *IOFloat32
}

type IOFloat32 struct {
	xxx_Float32ValueNoRange float32
	xxx_Float32WithRange    float32
}

func NewIOFloat32() *IOFloat32 {
	m := &IOFloat32{}
	m.Reset()
	return m
}

func (m *IOFloat32) Reset() {
	m.xxx_Float32ValueNoRange = 0
	m.xxx_Float32WithRange = 0
}

func (m *IOFloat32) CopyFrom(o *IOFloat32) *IOFloat32 {
	_ = m.UnmarshalFrame(o.Frame())
	return m
}

// Descriptor returns the IOFloat32 descriptor.
func (m *IOFloat32) Descriptor() *descriptor.Message {
	return Messages().IOFloat32.Message
}

// String returns a compact string representation of the message.
func (m *IOFloat32) String() string {
	return cantext.MessageString(m)
}

func (m *IOFloat32) Float32ValueNoRange() float32 {
	return m.xxx_Float32ValueNoRange
}

func (m *IOFloat32) SetFloat32ValueNoRange(v float32) *IOFloat32 {
	m.xxx_Float32ValueNoRange = float32(Messages().IOFloat32.Float32ValueNoRange.SaturatedCastFloat(float64(v)))
	return m
}

func (m *IOFloat32) Float32WithRange() float64 {
	return Messages().IOFloat32.Float32WithRange.ToPhysical(float64(m.xxx_Float32WithRange))
}

func (m *IOFloat32) SetFloat32WithRange(v float64) *IOFloat32 {
	m.xxx_Float32WithRange = float32(Messages().IOFloat32.Float32WithRange.FromPhysical(v))
	return m
}

func (m *IOFloat32) RawFloat32WithRange() float32 {
	return m.xxx_Float32WithRange
}

func (m *IOFloat32) SetRawFloat32WithRange(v float32) *IOFloat32 {
	m.xxx_Float32WithRange = float32(Messages().IOFloat32.Float32WithRange.SaturatedCastFloat(float64(v)))
	return m
}

// Frame returns a CAN frame representing the message.
func (m *IOFloat32) Frame() can.Frame {
	md := Messages().IOFloat32
	f := can.Frame{ID: md.ID, IsExtended: md.IsExtended, Length: md.Length}
	md.Float32ValueNoRange.MarshalFloat(&f.Data, float64(m.xxx_Float32ValueNoRange))
	md.Float32WithRange.MarshalFloat(&f.Data, float64(m.xxx_Float32WithRange))
	return f
}

// MarshalFrame encodes the message as a CAN frame.
func (m *IOFloat32) MarshalFrame() (can.Frame, error) {
	return m.Frame(), nil
}

// UnmarshalFrame decodes the message from a CAN frame.
func (m *IOFloat32) UnmarshalFrame(f can.Frame) error {
	md := Messages().IOFloat32
	switch {
	case f.ID != md.ID:
		return fmt.Errorf(
			"unmarshal IOFloat32: expects ID 600 (got %s with ID %d)", f.String(), f.ID,
		)
	case f.Length != md.Length:
		return fmt.Errorf(
			"unmarshal IOFloat32: expects length 8 (got %s with length %d)", f.String(), f.Length,
		)
	case f.IsRemote:
		return fmt.Errorf(
			"unmarshal IOFloat32: expects non-remote frame (got remote frame %s)", f.String(),
		)
	case f.IsExtended != md.IsExtended:
		return fmt.Errorf(
			"unmarshal IOFloat32: expects standard ID (got %s with extended ID)", f.String(),
		)
	}
	m.xxx_Float32ValueNoRange = float32(md.Float32ValueNoRange.UnmarshalFloat(f.Data))
	m.xxx_Float32WithRange = float32(md.Float32WithRange.UnmarshalFloat(f.Data))
	return nil
}

type DBG interface {
	sync.Locker
	Tx() DBG_Tx
	Rx() DBG_Rx
	Run(ctx context.Context) error
}

type DBG_Rx interface {
	http.Handler // for debugging
	SensorSonars() DBG_Rx_SensorSonars
	IODebug() DBG_Rx_IODebug
	IOFloat32() DBG_Rx_IOFloat32
}

type DBG_Tx interface {
	http.Handler // for debugging
}

type DBG_Rx_SensorSonars interface {
	SensorSonarsReader
	ReceiveTime() time.Time
	SetAfterReceiveHook(h func(context.Context) error)
}

type DBG_Rx_IODebug interface {
	IODebugReader
	ReceiveTime() time.Time
	SetAfterReceiveHook(h func(context.Context) error)
}

type DBG_Rx_IOFloat32 interface {
	IOFloat32Reader
	ReceiveTime() time.Time
	SetAfterReceiveHook(h func(context.Context) error)
}

type xxx_DBG struct {
	sync.Mutex // protects all node state
	network    string
	address    string
	rx         xxx_DBG_Rx
	tx         xxx_DBG_Tx
}

var _ DBG = &xxx_DBG{}
var _ canrunner.Node = &xxx_DBG{}

func NewDBG(network, address string) DBG {
	n := &xxx_DBG{network: network, address: address}
	n.rx.parentMutex = &n.Mutex
	n.tx.parentMutex = &n.Mutex
	n.rx.xxx_SensorSonars.init()
	n.rx.xxx_SensorSonars.Reset()
	n.rx.xxx_IODebug.init()
	n.rx.xxx_IODebug.Reset()
	n.rx.xxx_IOFloat32.init()
	n.rx.xxx_IOFloat32.Reset()
	return n
}

func (n *xxx_DBG) Run(ctx context.Context) error {
	return canrunner.Run(ctx, n)
}

func (n *xxx_DBG) Rx() DBG_Rx {
	return &n.rx
}

func (n *xxx_DBG) Tx() DBG_Tx {
	return &n.tx
}

type xxx_DBG_Rx struct {
	parentMutex      *sync.Mutex
	xxx_SensorSonars xxx_DBG_Rx_SensorSonars
	xxx_IODebug      xxx_DBG_Rx_IODebug
	xxx_IOFloat32    xxx_DBG_Rx_IOFloat32
}

var _ DBG_Rx = &xxx_DBG_Rx{}

func (rx *xxx_DBG_Rx) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	rx.parentMutex.Lock()
	defer rx.parentMutex.Unlock()
	candebug.ServeMessagesHTTP(w, r, []generated.Message{
		&rx.xxx_SensorSonars,
		&rx.xxx_IODebug,
		&rx.xxx_IOFloat32,
	})
}

func (rx *xxx_DBG_Rx) SensorSonars() DBG_Rx_SensorSonars {
	return &rx.xxx_SensorSonars
}

func (rx *xxx_DBG_Rx) IODebug() DBG_Rx_IODebug {
	return &rx.xxx_IODebug
}

func (rx *xxx_DBG_Rx) IOFloat32() DBG_Rx_IOFloat32 {
	return &rx.xxx_IOFloat32
}

type xxx_DBG_Tx struct {
	parentMutex *sync.Mutex
}

var _ DBG_Tx = &xxx_DBG_Tx{}

func (tx *xxx_DBG_Tx) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	tx.parentMutex.Lock()
	defer tx.parentMutex.Unlock()
	candebug.ServeMessagesHTTP(w, r, []generated.Message{})
}

func (n *xxx_DBG) Descriptor() *descriptor.Node {
	return Nodes().DBG
}

func (n *xxx_DBG) Connect() (net.Conn, error) {
	return socketcan.Dial(n.network, n.address)
}

func (n *xxx_DBG) ReceivedMessage(id uint32) (canrunner.ReceivedMessage, bool) {
	switch id {
	case 200:
		return &n.rx.xxx_SensorSonars, true
	case 500:
		return &n.rx.xxx_IODebug, true
	case 600:
		return &n.rx.xxx_IOFloat32, true
	default:
		return nil, false
	}
}

func (n *xxx_DBG) TransmittedMessages() []canrunner.TransmittedMessage {
	return []canrunner.TransmittedMessage{}
}

type xxx_DBG_Rx_SensorSonars struct {
	SensorSonars
	receiveTime      time.Time
	afterReceiveHook func(context.Context) error
}

func (m *xxx_DBG_Rx_SensorSonars) init() {
	m.afterReceiveHook = func(context.Context) error { return nil }
}

func (m *xxx_DBG_Rx_SensorSonars) SetAfterReceiveHook(h func(context.Context) error) {
	m.afterReceiveHook = h
}

func (m *xxx_DBG_Rx_SensorSonars) AfterReceiveHook() func(context.Context) error {
	return m.afterReceiveHook
}

func (m *xxx_DBG_Rx_SensorSonars) ReceiveTime() time.Time {
	return m.receiveTime
}

func (m *xxx_DBG_Rx_SensorSonars) SetReceiveTime(t time.Time) {
	m.receiveTime = t
}

var _ canrunner.ReceivedMessage = &xxx_DBG_Rx_SensorSonars{}

type xxx_DBG_Rx_IODebug struct {
	IODebug
	receiveTime      time.Time
	afterReceiveHook func(context.Context) error
}

func (m *xxx_DBG_Rx_IODebug) init() {
	m.afterReceiveHook = func(context.Context) error { return nil }
}

func (m *xxx_DBG_Rx_IODebug) SetAfterReceiveHook(h func(context.Context) error) {
	m.afterReceiveHook = h
}

func (m *xxx_DBG_Rx_IODebug) AfterReceiveHook() func(context.Context) error {
	return m.afterReceiveHook
}

func (m *xxx_DBG_Rx_IODebug) ReceiveTime() time.Time {
	return m.receiveTime
}

func (m *xxx_DBG_Rx_IODebug) SetReceiveTime(t time.Time) {
	m.receiveTime = t
}

var _ canrunner.ReceivedMessage = &xxx_DBG_Rx_IODebug{}

type xxx_DBG_Rx_IOFloat32 struct {
	IOFloat32
	receiveTime      time.Time
	afterReceiveHook func(context.Context) error
}

func (m *xxx_DBG_Rx_IOFloat32) init() {
	m.afterReceiveHook = func(context.Context) error { return nil }
}

func (m *xxx_DBG_Rx_IOFloat32) SetAfterReceiveHook(h func(context.Context) error) {
	m.afterReceiveHook = h
}

func (m *xxx_DBG_Rx_IOFloat32) AfterReceiveHook() func(context.Context) error {
	return m.afterReceiveHook
}

func (m *xxx_DBG_Rx_IOFloat32) ReceiveTime() time.Time {
	return m.receiveTime
}

func (m *xxx_DBG_Rx_IOFloat32) SetReceiveTime(t time.Time) {
	m.receiveTime = t
}

var _ canrunner.ReceivedMessage = &xxx_DBG_Rx_IOFloat32{}

type DRIVER interface {
	sync.Locker
	Tx() DRIVER_Tx
	Rx() DRIVER_Rx
	Run(ctx context.Context) error
}

type DRIVER_Rx interface {
	http.Handler // for debugging
	SensorSonars() DRIVER_Rx_SensorSonars
	MotorStatus() DRIVER_Rx_MotorStatus
}

type DRIVER_Tx interface {
	http.Handler // for debugging
	DriverHeartbeat() DRIVER_Tx_DriverHeartbeat
	MotorCommand() DRIVER_Tx_MotorCommand
}

type DRIVER_Rx_SensorSonars interface {
	SensorSonarsReader
	ReceiveTime() time.Time
	SetAfterReceiveHook(h func(context.Context) error)
}

type DRIVER_Rx_MotorStatus interface {
	MotorStatusReader
	ReceiveTime() time.Time
	SetAfterReceiveHook(h func(context.Context) error)
}

type DRIVER_Tx_DriverHeartbeat interface {
	DriverHeartbeatReader
	DriverHeartbeatWriter
	TransmitTime() time.Time
	Transmit(ctx context.Context) error
	SetBeforeTransmitHook(h func(context.Context) error)
	// SetCyclicTransmissionEnabled enables/disables cyclic transmission.
	SetCyclicTransmissionEnabled(bool)
	// IsCyclicTransmissionEnabled returns whether cyclic transmission is enabled/disabled.
	IsCyclicTransmissionEnabled() bool
}

type DRIVER_Tx_MotorCommand interface {
	MotorCommandReader
	MotorCommandWriter
	TransmitTime() time.Time
	Transmit(ctx context.Context) error
	SetBeforeTransmitHook(h func(context.Context) error)
	// SetCyclicTransmissionEnabled enables/disables cyclic transmission.
	SetCyclicTransmissionEnabled(bool)
	// IsCyclicTransmissionEnabled returns whether cyclic transmission is enabled/disabled.
	IsCyclicTransmissionEnabled() bool
}

type xxx_DRIVER struct {
	sync.Mutex // protects all node state
	network    string
	address    string
	rx         xxx_DRIVER_Rx
	tx         xxx_DRIVER_Tx
}

var _ DRIVER = &xxx_DRIVER{}
var _ canrunner.Node = &xxx_DRIVER{}

func NewDRIVER(network, address string) DRIVER {
	n := &xxx_DRIVER{network: network, address: address}
	n.rx.parentMutex = &n.Mutex
	n.tx.parentMutex = &n.Mutex
	n.rx.xxx_SensorSonars.init()
	n.rx.xxx_SensorSonars.Reset()
	n.rx.xxx_MotorStatus.init()
	n.rx.xxx_MotorStatus.Reset()
	n.tx.xxx_DriverHeartbeat.init()
	n.tx.xxx_DriverHeartbeat.Reset()
	n.tx.xxx_MotorCommand.init()
	n.tx.xxx_MotorCommand.Reset()
	return n
}

func (n *xxx_DRIVER) Run(ctx context.Context) error {
	return canrunner.Run(ctx, n)
}

func (n *xxx_DRIVER) Rx() DRIVER_Rx {
	return &n.rx
}

func (n *xxx_DRIVER) Tx() DRIVER_Tx {
	return &n.tx
}

type xxx_DRIVER_Rx struct {
	parentMutex      *sync.Mutex
	xxx_SensorSonars xxx_DRIVER_Rx_SensorSonars
	xxx_MotorStatus  xxx_DRIVER_Rx_MotorStatus
}

var _ DRIVER_Rx = &xxx_DRIVER_Rx{}

func (rx *xxx_DRIVER_Rx) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	rx.parentMutex.Lock()
	defer rx.parentMutex.Unlock()
	candebug.ServeMessagesHTTP(w, r, []generated.Message{
		&rx.xxx_SensorSonars,
		&rx.xxx_MotorStatus,
	})
}

func (rx *xxx_DRIVER_Rx) SensorSonars() DRIVER_Rx_SensorSonars {
	return &rx.xxx_SensorSonars
}

func (rx *xxx_DRIVER_Rx) MotorStatus() DRIVER_Rx_MotorStatus {
	return &rx.xxx_MotorStatus
}

type xxx_DRIVER_Tx struct {
	parentMutex         *sync.Mutex
	xxx_DriverHeartbeat xxx_DRIVER_Tx_DriverHeartbeat
	xxx_MotorCommand    xxx_DRIVER_Tx_MotorCommand
}

var _ DRIVER_Tx = &xxx_DRIVER_Tx{}

func (tx *xxx_DRIVER_Tx) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	tx.parentMutex.Lock()
	defer tx.parentMutex.Unlock()
	candebug.ServeMessagesHTTP(w, r, []generated.Message{
		&tx.xxx_DriverHeartbeat,
		&tx.xxx_MotorCommand,
	})
}

func (tx *xxx_DRIVER_Tx) DriverHeartbeat() DRIVER_Tx_DriverHeartbeat {
	return &tx.xxx_DriverHeartbeat
}

func (tx *xxx_DRIVER_Tx) MotorCommand() DRIVER_Tx_MotorCommand {
	return &tx.xxx_MotorCommand
}

func (n *xxx_DRIVER) Descriptor() *descriptor.Node {
	return Nodes().DRIVER
}

func (n *xxx_DRIVER) Connect() (net.Conn, error) {
	return socketcan.Dial(n.network, n.address)
}

func (n *xxx_DRIVER) ReceivedMessage(id uint32) (canrunner.ReceivedMessage, bool) {
	switch id {
	case 200:
		return &n.rx.xxx_SensorSonars, true
	case 400:
		return &n.rx.xxx_MotorStatus, true
	default:
		return nil, false
	}
}

func (n *xxx_DRIVER) TransmittedMessages() []canrunner.TransmittedMessage {
	return []canrunner.TransmittedMessage{
		&n.tx.xxx_DriverHeartbeat,
		&n.tx.xxx_MotorCommand,
	}
}

type xxx_DRIVER_Rx_SensorSonars struct {
	SensorSonars
	receiveTime      time.Time
	afterReceiveHook func(context.Context) error
}

func (m *xxx_DRIVER_Rx_SensorSonars) init() {
	m.afterReceiveHook = func(context.Context) error { return nil }
}

func (m *xxx_DRIVER_Rx_SensorSonars) SetAfterReceiveHook(h func(context.Context) error) {
	m.afterReceiveHook = h
}

func (m *xxx_DRIVER_Rx_SensorSonars) AfterReceiveHook() func(context.Context) error {
	return m.afterReceiveHook
}

func (m *xxx_DRIVER_Rx_SensorSonars) ReceiveTime() time.Time {
	return m.receiveTime
}

func (m *xxx_DRIVER_Rx_SensorSonars) SetReceiveTime(t time.Time) {
	m.receiveTime = t
}

var _ canrunner.ReceivedMessage = &xxx_DRIVER_Rx_SensorSonars{}

type xxx_DRIVER_Rx_MotorStatus struct {
	MotorStatus
	receiveTime      time.Time
	afterReceiveHook func(context.Context) error
}

func (m *xxx_DRIVER_Rx_MotorStatus) init() {
	m.afterReceiveHook = func(context.Context) error { return nil }
}

func (m *xxx_DRIVER_Rx_MotorStatus) SetAfterReceiveHook(h func(context.Context) error) {
	m.afterReceiveHook = h
}

func (m *xxx_DRIVER_Rx_MotorStatus) AfterReceiveHook() func(context.Context) error {
	return m.afterReceiveHook
}

func (m *xxx_DRIVER_Rx_MotorStatus) ReceiveTime() time.Time {
	return m.receiveTime
}

func (m *xxx_DRIVER_Rx_MotorStatus) SetReceiveTime(t time.Time) {
	m.receiveTime = t
}

var _ canrunner.ReceivedMessage = &xxx_DRIVER_Rx_MotorStatus{}

type xxx_DRIVER_Tx_DriverHeartbeat struct {
	DriverHeartbeat
	transmitTime       time.Time
	beforeTransmitHook func(context.Context) error
	isCyclicEnabled    bool
	wakeUpChan         chan struct{}
	transmitEventChan  chan struct{}
}

var _ DRIVER_Tx_DriverHeartbeat = &xxx_DRIVER_Tx_DriverHeartbeat{}
var _ canrunner.TransmittedMessage = &xxx_DRIVER_Tx_DriverHeartbeat{}

func (m *xxx_DRIVER_Tx_DriverHeartbeat) init() {
	m.beforeTransmitHook = func(context.Context) error { return nil }
	m.wakeUpChan = make(chan struct{}, 1)
	m.transmitEventChan = make(chan struct{})
}

func (m *xxx_DRIVER_Tx_DriverHeartbeat) SetBeforeTransmitHook(h func(context.Context) error) {
	m.beforeTransmitHook = h
}

func (m *xxx_DRIVER_Tx_DriverHeartbeat) BeforeTransmitHook() func(context.Context) error {
	return m.beforeTransmitHook
}

func (m *xxx_DRIVER_Tx_DriverHeartbeat) TransmitTime() time.Time {
	return m.transmitTime
}

func (m *xxx_DRIVER_Tx_DriverHeartbeat) SetTransmitTime(t time.Time) {
	m.transmitTime = t
}

func (m *xxx_DRIVER_Tx_DriverHeartbeat) IsCyclicTransmissionEnabled() bool {
	return m.isCyclicEnabled
}

func (m *xxx_DRIVER_Tx_DriverHeartbeat) SetCyclicTransmissionEnabled(b bool) {
	m.isCyclicEnabled = b
	select {
	case m.wakeUpChan <- struct{}{}:
	default:
	}
}

func (m *xxx_DRIVER_Tx_DriverHeartbeat) WakeUpChan() <-chan struct{} {
	return m.wakeUpChan
}

func (m *xxx_DRIVER_Tx_DriverHeartbeat) Transmit(ctx context.Context) error {
	select {
	case m.transmitEventChan <- struct{}{}:
		return nil
	case <-ctx.Done():
		return fmt.Errorf("event-triggered transmit of DriverHeartbeat: %w", ctx.Err())
	}
}

func (m *xxx_DRIVER_Tx_DriverHeartbeat) TransmitEventChan() <-chan struct{} {
	return m.transmitEventChan
}

var _ canrunner.TransmittedMessage = &xxx_DRIVER_Tx_DriverHeartbeat{}

type xxx_DRIVER_Tx_MotorCommand struct {
	MotorCommand
	transmitTime       time.Time
	beforeTransmitHook func(context.Context) error
	isCyclicEnabled    bool
	wakeUpChan         chan struct{}
	transmitEventChan  chan struct{}
}

var _ DRIVER_Tx_MotorCommand = &xxx_DRIVER_Tx_MotorCommand{}
var _ canrunner.TransmittedMessage = &xxx_DRIVER_Tx_MotorCommand{}

func (m *xxx_DRIVER_Tx_MotorCommand) init() {
	m.beforeTransmitHook = func(context.Context) error { return nil }
	m.wakeUpChan = make(chan struct{}, 1)
	m.transmitEventChan = make(chan struct{})
}

func (m *xxx_DRIVER_Tx_MotorCommand) SetBeforeTransmitHook(h func(context.Context) error) {
	m.beforeTransmitHook = h
}

func (m *xxx_DRIVER_Tx_MotorCommand) BeforeTransmitHook() func(context.Context) error {
	return m.beforeTransmitHook
}

func (m *xxx_DRIVER_Tx_MotorCommand) TransmitTime() time.Time {
	return m.transmitTime
}

func (m *xxx_DRIVER_Tx_MotorCommand) SetTransmitTime(t time.Time) {
	m.transmitTime = t
}

func (m *xxx_DRIVER_Tx_MotorCommand) IsCyclicTransmissionEnabled() bool {
	return m.isCyclicEnabled
}

func (m *xxx_DRIVER_Tx_MotorCommand) SetCyclicTransmissionEnabled(b bool) {
	m.isCyclicEnabled = b
	select {
	case m.wakeUpChan <- struct{}{}:
	default:
	}
}

func (m *xxx_DRIVER_Tx_MotorCommand) WakeUpChan() <-chan struct{} {
	return m.wakeUpChan
}

func (m *xxx_DRIVER_Tx_MotorCommand) Transmit(ctx context.Context) error {
	select {
	case m.transmitEventChan <- struct{}{}:
		return nil
	case <-ctx.Done():
		return fmt.Errorf("event-triggered transmit of MotorCommand: %w", ctx.Err())
	}
}

func (m *xxx_DRIVER_Tx_MotorCommand) TransmitEventChan() <-chan struct{} {
	return m.transmitEventChan
}

var _ canrunner.TransmittedMessage = &xxx_DRIVER_Tx_MotorCommand{}

type IO interface {
	sync.Locker
	Tx() IO_Tx
	Rx() IO_Rx
	Run(ctx context.Context) error
}

type IO_Rx interface {
	http.Handler // for debugging
	SensorSonars() IO_Rx_SensorSonars
	MotorStatus() IO_Rx_MotorStatus
}

type IO_Tx interface {
	http.Handler // for debugging
	IODebug() IO_Tx_IODebug
}

type IO_Rx_SensorSonars interface {
	SensorSonarsReader
	ReceiveTime() time.Time
	SetAfterReceiveHook(h func(context.Context) error)
}

type IO_Rx_MotorStatus interface {
	MotorStatusReader
	ReceiveTime() time.Time
	SetAfterReceiveHook(h func(context.Context) error)
}

type IO_Tx_IODebug interface {
	IODebugReader
	IODebugWriter
	TransmitTime() time.Time
	Transmit(ctx context.Context) error
	SetBeforeTransmitHook(h func(context.Context) error)
}

type xxx_IO struct {
	sync.Mutex // protects all node state
	network    string
	address    string
	rx         xxx_IO_Rx
	tx         xxx_IO_Tx
}

var _ IO = &xxx_IO{}
var _ canrunner.Node = &xxx_IO{}

func NewIO(network, address string) IO {
	n := &xxx_IO{network: network, address: address}
	n.rx.parentMutex = &n.Mutex
	n.tx.parentMutex = &n.Mutex
	n.rx.xxx_SensorSonars.init()
	n.rx.xxx_SensorSonars.Reset()
	n.rx.xxx_MotorStatus.init()
	n.rx.xxx_MotorStatus.Reset()
	n.tx.xxx_IODebug.init()
	n.tx.xxx_IODebug.Reset()
	return n
}

func (n *xxx_IO) Run(ctx context.Context) error {
	return canrunner.Run(ctx, n)
}

func (n *xxx_IO) Rx() IO_Rx {
	return &n.rx
}

func (n *xxx_IO) Tx() IO_Tx {
	return &n.tx
}

type xxx_IO_Rx struct {
	parentMutex      *sync.Mutex
	xxx_SensorSonars xxx_IO_Rx_SensorSonars
	xxx_MotorStatus  xxx_IO_Rx_MotorStatus
}

var _ IO_Rx = &xxx_IO_Rx{}

func (rx *xxx_IO_Rx) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	rx.parentMutex.Lock()
	defer rx.parentMutex.Unlock()
	candebug.ServeMessagesHTTP(w, r, []generated.Message{
		&rx.xxx_SensorSonars,
		&rx.xxx_MotorStatus,
	})
}

func (rx *xxx_IO_Rx) SensorSonars() IO_Rx_SensorSonars {
	return &rx.xxx_SensorSonars
}

func (rx *xxx_IO_Rx) MotorStatus() IO_Rx_MotorStatus {
	return &rx.xxx_MotorStatus
}

type xxx_IO_Tx struct {
	parentMutex *sync.Mutex
	xxx_IODebug xxx_IO_Tx_IODebug
}

var _ IO_Tx = &xxx_IO_Tx{}

func (tx *xxx_IO_Tx) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	tx.parentMutex.Lock()
	defer tx.parentMutex.Unlock()
	candebug.ServeMessagesHTTP(w, r, []generated.Message{
		&tx.xxx_IODebug,
	})
}

func (tx *xxx_IO_Tx) IODebug() IO_Tx_IODebug {
	return &tx.xxx_IODebug
}

func (n *xxx_IO) Descriptor() *descriptor.Node {
	return Nodes().IO
}

func (n *xxx_IO) Connect() (net.Conn, error) {
	return socketcan.Dial(n.network, n.address)
}

func (n *xxx_IO) ReceivedMessage(id uint32) (canrunner.ReceivedMessage, bool) {
	switch id {
	case 200:
		return &n.rx.xxx_SensorSonars, true
	case 400:
		return &n.rx.xxx_MotorStatus, true
	default:
		return nil, false
	}
}

func (n *xxx_IO) TransmittedMessages() []canrunner.TransmittedMessage {
	return []canrunner.TransmittedMessage{
		&n.tx.xxx_IODebug,
	}
}

type xxx_IO_Rx_SensorSonars struct {
	SensorSonars
	receiveTime      time.Time
	afterReceiveHook func(context.Context) error
}

func (m *xxx_IO_Rx_SensorSonars) init() {
	m.afterReceiveHook = func(context.Context) error { return nil }
}

func (m *xxx_IO_Rx_SensorSonars) SetAfterReceiveHook(h func(context.Context) error) {
	m.afterReceiveHook = h
}

func (m *xxx_IO_Rx_SensorSonars) AfterReceiveHook() func(context.Context) error {
	return m.afterReceiveHook
}

func (m *xxx_IO_Rx_SensorSonars) ReceiveTime() time.Time {
	return m.receiveTime
}

func (m *xxx_IO_Rx_SensorSonars) SetReceiveTime(t time.Time) {
	m.receiveTime = t
}

var _ canrunner.ReceivedMessage = &xxx_IO_Rx_SensorSonars{}

type xxx_IO_Rx_MotorStatus struct {
	MotorStatus
	receiveTime      time.Time
	afterReceiveHook func(context.Context) error
}

func (m *xxx_IO_Rx_MotorStatus) init() {
	m.afterReceiveHook = func(context.Context) error { return nil }
}

func (m *xxx_IO_Rx_MotorStatus) SetAfterReceiveHook(h func(context.Context) error) {
	m.afterReceiveHook = h
}

func (m *xxx_IO_Rx_MotorStatus) AfterReceiveHook() func(context.Context) error {
	return m.afterReceiveHook
}

func (m *xxx_IO_Rx_MotorStatus) ReceiveTime() time.Time {
	return m.receiveTime
}

func (m *xxx_IO_Rx_MotorStatus) SetReceiveTime(t time.Time) {
	m.receiveTime = t
}

var _ canrunner.ReceivedMessage = &xxx_IO_Rx_MotorStatus{}

type xxx_IO_Tx_IODebug struct {
	IODebug
	transmitTime       time.Time
	beforeTransmitHook func(context.Context) error
	isCyclicEnabled    bool
	wakeUpChan         chan struct{}
	transmitEventChan  chan struct{}
}

var _ IO_Tx_IODebug = &xxx_IO_Tx_IODebug{}
var _ canrunner.TransmittedMessage = &xxx_IO_Tx_IODebug{}

func (m *xxx_IO_Tx_IODebug) init() {
	m.beforeTransmitHook = func(context.Context) error { return nil }
	m.wakeUpChan = make(chan struct{}, 1)
	m.transmitEventChan = make(chan struct{})
}

func (m *xxx_IO_Tx_IODebug) SetBeforeTransmitHook(h func(context.Context) error) {
	m.beforeTransmitHook = h
}

func (m *xxx_IO_Tx_IODebug) BeforeTransmitHook() func(context.Context) error {
	return m.beforeTransmitHook
}

func (m *xxx_IO_Tx_IODebug) TransmitTime() time.Time {
	return m.transmitTime
}

func (m *xxx_IO_Tx_IODebug) SetTransmitTime(t time.Time) {
	m.transmitTime = t
}

func (m *xxx_IO_Tx_IODebug) IsCyclicTransmissionEnabled() bool {
	return m.isCyclicEnabled
}

func (m *xxx_IO_Tx_IODebug) SetCyclicTransmissionEnabled(b bool) {
	m.isCyclicEnabled = b
	select {
	case m.wakeUpChan <- struct{}{}:
	default:
	}
}

func (m *xxx_IO_Tx_IODebug) WakeUpChan() <-chan struct{} {
	return m.wakeUpChan
}

func (m *xxx_IO_Tx_IODebug) Transmit(ctx context.Context) error {
	select {
	case m.transmitEventChan <- struct{}{}:
		return nil
	case <-ctx.Done():
		return fmt.Errorf("event-triggered transmit of IODebug: %w", ctx.Err())
	}
}

func (m *xxx_IO_Tx_IODebug) TransmitEventChan() <-chan struct{} {
	return m.transmitEventChan
}

var _ canrunner.TransmittedMessage = &xxx_IO_Tx_IODebug{}

type MOTOR interface {
	sync.Locker
	Tx() MOTOR_Tx
	Rx() MOTOR_Rx
	Run(ctx context.Context) error
}

type MOTOR_Rx interface {
	http.Handler // for debugging
	DriverHeartbeat() MOTOR_Rx_DriverHeartbeat
	MotorCommand() MOTOR_Rx_MotorCommand
}

type MOTOR_Tx interface {
	http.Handler // for debugging
	MotorStatus() MOTOR_Tx_MotorStatus
}

type MOTOR_Rx_DriverHeartbeat interface {
	DriverHeartbeatReader
	ReceiveTime() time.Time
	SetAfterReceiveHook(h func(context.Context) error)
}

type MOTOR_Rx_MotorCommand interface {
	MotorCommandReader
	ReceiveTime() time.Time
	SetAfterReceiveHook(h func(context.Context) error)
}

type MOTOR_Tx_MotorStatus interface {
	MotorStatusReader
	MotorStatusWriter
	TransmitTime() time.Time
	Transmit(ctx context.Context) error
	SetBeforeTransmitHook(h func(context.Context) error)
	// SetCyclicTransmissionEnabled enables/disables cyclic transmission.
	SetCyclicTransmissionEnabled(bool)
	// IsCyclicTransmissionEnabled returns whether cyclic transmission is enabled/disabled.
	IsCyclicTransmissionEnabled() bool
}

type xxx_MOTOR struct {
	sync.Mutex // protects all node state
	network    string
	address    string
	rx         xxx_MOTOR_Rx
	tx         xxx_MOTOR_Tx
}

var _ MOTOR = &xxx_MOTOR{}
var _ canrunner.Node = &xxx_MOTOR{}

func NewMOTOR(network, address string) MOTOR {
	n := &xxx_MOTOR{network: network, address: address}
	n.rx.parentMutex = &n.Mutex
	n.tx.parentMutex = &n.Mutex
	n.rx.xxx_DriverHeartbeat.init()
	n.rx.xxx_DriverHeartbeat.Reset()
	n.rx.xxx_MotorCommand.init()
	n.rx.xxx_MotorCommand.Reset()
	n.tx.xxx_MotorStatus.init()
	n.tx.xxx_MotorStatus.Reset()
	return n
}

func (n *xxx_MOTOR) Run(ctx context.Context) error {
	return canrunner.Run(ctx, n)
}

func (n *xxx_MOTOR) Rx() MOTOR_Rx {
	return &n.rx
}

func (n *xxx_MOTOR) Tx() MOTOR_Tx {
	return &n.tx
}

type xxx_MOTOR_Rx struct {
	parentMutex         *sync.Mutex
	xxx_DriverHeartbeat xxx_MOTOR_Rx_DriverHeartbeat
	xxx_MotorCommand    xxx_MOTOR_Rx_MotorCommand
}

var _ MOTOR_Rx = &xxx_MOTOR_Rx{}

func (rx *xxx_MOTOR_Rx) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	rx.parentMutex.Lock()
	defer rx.parentMutex.Unlock()
	candebug.ServeMessagesHTTP(w, r, []generated.Message{
		&rx.xxx_DriverHeartbeat,
		&rx.xxx_MotorCommand,
	})
}

func (rx *xxx_MOTOR_Rx) DriverHeartbeat() MOTOR_Rx_DriverHeartbeat {
	return &rx.xxx_DriverHeartbeat
}

func (rx *xxx_MOTOR_Rx) MotorCommand() MOTOR_Rx_MotorCommand {
	return &rx.xxx_MotorCommand
}

type xxx_MOTOR_Tx struct {
	parentMutex     *sync.Mutex
	xxx_MotorStatus xxx_MOTOR_Tx_MotorStatus
}

var _ MOTOR_Tx = &xxx_MOTOR_Tx{}

func (tx *xxx_MOTOR_Tx) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	tx.parentMutex.Lock()
	defer tx.parentMutex.Unlock()
	candebug.ServeMessagesHTTP(w, r, []generated.Message{
		&tx.xxx_MotorStatus,
	})
}

func (tx *xxx_MOTOR_Tx) MotorStatus() MOTOR_Tx_MotorStatus {
	return &tx.xxx_MotorStatus
}

func (n *xxx_MOTOR) Descriptor() *descriptor.Node {
	return Nodes().MOTOR
}

func (n *xxx_MOTOR) Connect() (net.Conn, error) {
	return socketcan.Dial(n.network, n.address)
}

func (n *xxx_MOTOR) ReceivedMessage(id uint32) (canrunner.ReceivedMessage, bool) {
	switch id {
	case 100:
		return &n.rx.xxx_DriverHeartbeat, true
	case 101:
		return &n.rx.xxx_MotorCommand, true
	default:
		return nil, false
	}
}

func (n *xxx_MOTOR) TransmittedMessages() []canrunner.TransmittedMessage {
	return []canrunner.TransmittedMessage{
		&n.tx.xxx_MotorStatus,
	}
}

type xxx_MOTOR_Rx_DriverHeartbeat struct {
	DriverHeartbeat
	receiveTime      time.Time
	afterReceiveHook func(context.Context) error
}

func (m *xxx_MOTOR_Rx_DriverHeartbeat) init() {
	m.afterReceiveHook = func(context.Context) error { return nil }
}

func (m *xxx_MOTOR_Rx_DriverHeartbeat) SetAfterReceiveHook(h func(context.Context) error) {
	m.afterReceiveHook = h
}

func (m *xxx_MOTOR_Rx_DriverHeartbeat) AfterReceiveHook() func(context.Context) error {
	return m.afterReceiveHook
}

func (m *xxx_MOTOR_Rx_DriverHeartbeat) ReceiveTime() time.Time {
	return m.receiveTime
}

func (m *xxx_MOTOR_Rx_DriverHeartbeat) SetReceiveTime(t time.Time) {
	m.receiveTime = t
}

var _ canrunner.ReceivedMessage = &xxx_MOTOR_Rx_DriverHeartbeat{}

type xxx_MOTOR_Rx_MotorCommand struct {
	MotorCommand
	receiveTime      time.Time
	afterReceiveHook func(context.Context) error
}

func (m *xxx_MOTOR_Rx_MotorCommand) init() {
	m.afterReceiveHook = func(context.Context) error { return nil }
}

func (m *xxx_MOTOR_Rx_MotorCommand) SetAfterReceiveHook(h func(context.Context) error) {
	m.afterReceiveHook = h
}

func (m *xxx_MOTOR_Rx_MotorCommand) AfterReceiveHook() func(context.Context) error {
	return m.afterReceiveHook
}

func (m *xxx_MOTOR_Rx_MotorCommand) ReceiveTime() time.Time {
	return m.receiveTime
}

func (m *xxx_MOTOR_Rx_MotorCommand) SetReceiveTime(t time.Time) {
	m.receiveTime = t
}

var _ canrunner.ReceivedMessage = &xxx_MOTOR_Rx_MotorCommand{}

type xxx_MOTOR_Tx_MotorStatus struct {
	MotorStatus
	transmitTime       time.Time
	beforeTransmitHook func(context.Context) error
	isCyclicEnabled    bool
	wakeUpChan         chan struct{}
	transmitEventChan  chan struct{}
}

var _ MOTOR_Tx_MotorStatus = &xxx_MOTOR_Tx_MotorStatus{}
var _ canrunner.TransmittedMessage = &xxx_MOTOR_Tx_MotorStatus{}

func (m *xxx_MOTOR_Tx_MotorStatus) init() {
	m.beforeTransmitHook = func(context.Context) error { return nil }
	m.wakeUpChan = make(chan struct{}, 1)
	m.transmitEventChan = make(chan struct{})
}

func (m *xxx_MOTOR_Tx_MotorStatus) SetBeforeTransmitHook(h func(context.Context) error) {
	m.beforeTransmitHook = h
}

func (m *xxx_MOTOR_Tx_MotorStatus) BeforeTransmitHook() func(context.Context) error {
	return m.beforeTransmitHook
}

func (m *xxx_MOTOR_Tx_MotorStatus) TransmitTime() time.Time {
	return m.transmitTime
}

func (m *xxx_MOTOR_Tx_MotorStatus) SetTransmitTime(t time.Time) {
	m.transmitTime = t
}

func (m *xxx_MOTOR_Tx_MotorStatus) IsCyclicTransmissionEnabled() bool {
	return m.isCyclicEnabled
}

func (m *xxx_MOTOR_Tx_MotorStatus) SetCyclicTransmissionEnabled(b bool) {
	m.isCyclicEnabled = b
	select {
	case m.wakeUpChan <- struct{}{}:
	default:
	}
}

func (m *xxx_MOTOR_Tx_MotorStatus) WakeUpChan() <-chan struct{} {
	return m.wakeUpChan
}

func (m *xxx_MOTOR_Tx_MotorStatus) Transmit(ctx context.Context) error {
	select {
	case m.transmitEventChan <- struct{}{}:
		return nil
	case <-ctx.Done():
		return fmt.Errorf("event-triggered transmit of MotorStatus: %w", ctx.Err())
	}
}

func (m *xxx_MOTOR_Tx_MotorStatus) TransmitEventChan() <-chan struct{} {
	return m.transmitEventChan
}

var _ canrunner.TransmittedMessage = &xxx_MOTOR_Tx_MotorStatus{}

type SENSOR interface {
	sync.Locker
	Tx() SENSOR_Tx
	Rx() SENSOR_Rx
	Run(ctx context.Context) error
}

type SENSOR_Rx interface {
	http.Handler // for debugging
	DriverHeartbeat() SENSOR_Rx_DriverHeartbeat
}

type SENSOR_Tx interface {
	http.Handler // for debugging
	SensorSonars() SENSOR_Tx_SensorSonars
}

type SENSOR_Rx_DriverHeartbeat interface {
	DriverHeartbeatReader
	ReceiveTime() time.Time
	SetAfterReceiveHook(h func(context.Context) error)
}

type SENSOR_Tx_SensorSonars interface {
	SensorSonarsReader
	SensorSonarsWriter
	TransmitTime() time.Time
	Transmit(ctx context.Context) error
	SetBeforeTransmitHook(h func(context.Context) error)
	// SetCyclicTransmissionEnabled enables/disables cyclic transmission.
	SetCyclicTransmissionEnabled(bool)
	// IsCyclicTransmissionEnabled returns whether cyclic transmission is enabled/disabled.
	IsCyclicTransmissionEnabled() bool
}

type xxx_SENSOR struct {
	sync.Mutex // protects all node state
	network    string
	address    string
	rx         xxx_SENSOR_Rx
	tx         xxx_SENSOR_Tx
}

var _ SENSOR = &xxx_SENSOR{}
var _ canrunner.Node = &xxx_SENSOR{}

func NewSENSOR(network, address string) SENSOR {
	n := &xxx_SENSOR{network: network, address: address}
	n.rx.parentMutex = &n.Mutex
	n.tx.parentMutex = &n.Mutex
	n.rx.xxx_DriverHeartbeat.init()
	n.rx.xxx_DriverHeartbeat.Reset()
	n.tx.xxx_SensorSonars.init()
	n.tx.xxx_SensorSonars.Reset()
	return n
}

func (n *xxx_SENSOR) Run(ctx context.Context) error {
	return canrunner.Run(ctx, n)
}

func (n *xxx_SENSOR) Rx() SENSOR_Rx {
	return &n.rx
}

func (n *xxx_SENSOR) Tx() SENSOR_Tx {
	return &n.tx
}

type xxx_SENSOR_Rx struct {
	parentMutex         *sync.Mutex
	xxx_DriverHeartbeat xxx_SENSOR_Rx_DriverHeartbeat
}

var _ SENSOR_Rx = &xxx_SENSOR_Rx{}

func (rx *xxx_SENSOR_Rx) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	rx.parentMutex.Lock()
	defer rx.parentMutex.Unlock()
	candebug.ServeMessagesHTTP(w, r, []generated.Message{
		&rx.xxx_DriverHeartbeat,
	})
}

func (rx *xxx_SENSOR_Rx) DriverHeartbeat() SENSOR_Rx_DriverHeartbeat {
	return &rx.xxx_DriverHeartbeat
}

type xxx_SENSOR_Tx struct {
	parentMutex      *sync.Mutex
	xxx_SensorSonars xxx_SENSOR_Tx_SensorSonars
}

var _ SENSOR_Tx = &xxx_SENSOR_Tx{}

func (tx *xxx_SENSOR_Tx) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	tx.parentMutex.Lock()
	defer tx.parentMutex.Unlock()
	candebug.ServeMessagesHTTP(w, r, []generated.Message{
		&tx.xxx_SensorSonars,
	})
}

func (tx *xxx_SENSOR_Tx) SensorSonars() SENSOR_Tx_SensorSonars {
	return &tx.xxx_SensorSonars
}

func (n *xxx_SENSOR) Descriptor() *descriptor.Node {
	return Nodes().SENSOR
}

func (n *xxx_SENSOR) Connect() (net.Conn, error) {
	return socketcan.Dial(n.network, n.address)
}

func (n *xxx_SENSOR) ReceivedMessage(id uint32) (canrunner.ReceivedMessage, bool) {
	switch id {
	case 100:
		return &n.rx.xxx_DriverHeartbeat, true
	default:
		return nil, false
	}
}

func (n *xxx_SENSOR) TransmittedMessages() []canrunner.TransmittedMessage {
	return []canrunner.TransmittedMessage{
		&n.tx.xxx_SensorSonars,
	}
}

type xxx_SENSOR_Rx_DriverHeartbeat struct {
	DriverHeartbeat
	receiveTime      time.Time
	afterReceiveHook func(context.Context) error
}

func (m *xxx_SENSOR_Rx_DriverHeartbeat) init() {
	m.afterReceiveHook = func(context.Context) error { return nil }
}

func (m *xxx_SENSOR_Rx_DriverHeartbeat) SetAfterReceiveHook(h func(context.Context) error) {
	m.afterReceiveHook = h
}

func (m *xxx_SENSOR_Rx_DriverHeartbeat) AfterReceiveHook() func(context.Context) error {
	return m.afterReceiveHook
}

func (m *xxx_SENSOR_Rx_DriverHeartbeat) ReceiveTime() time.Time {
	return m.receiveTime
}

func (m *xxx_SENSOR_Rx_DriverHeartbeat) SetReceiveTime(t time.Time) {
	m.receiveTime = t
}

var _ canrunner.ReceivedMessage = &xxx_SENSOR_Rx_DriverHeartbeat{}

type xxx_SENSOR_Tx_SensorSonars struct {
	SensorSonars
	transmitTime       time.Time
	beforeTransmitHook func(context.Context) error
	isCyclicEnabled    bool
	wakeUpChan         chan struct{}
	transmitEventChan  chan struct{}
}

var _ SENSOR_Tx_SensorSonars = &xxx_SENSOR_Tx_SensorSonars{}
var _ canrunner.TransmittedMessage = &xxx_SENSOR_Tx_SensorSonars{}

func (m *xxx_SENSOR_Tx_SensorSonars) init() {
	m.beforeTransmitHook = func(context.Context) error { return nil }
	m.wakeUpChan = make(chan struct{}, 1)
	m.transmitEventChan = make(chan struct{})
}

func (m *xxx_SENSOR_Tx_SensorSonars) SetBeforeTransmitHook(h func(context.Context) error) {
	m.beforeTransmitHook = h
}

func (m *xxx_SENSOR_Tx_SensorSonars) BeforeTransmitHook() func(context.Context) error {
	return m.beforeTransmitHook
}

func (m *xxx_SENSOR_Tx_SensorSonars) TransmitTime() time.Time {
	return m.transmitTime
}

func (m *xxx_SENSOR_Tx_SensorSonars) SetTransmitTime(t time.Time) {
	m.transmitTime = t
}

func (m *xxx_SENSOR_Tx_SensorSonars) IsCyclicTransmissionEnabled() bool {
	return m.isCyclicEnabled
}

func (m *xxx_SENSOR_Tx_SensorSonars) SetCyclicTransmissionEnabled(b bool) {
	m.isCyclicEnabled = b
	select {
	case m.wakeUpChan <- struct{}{}:
	default:
	}
}

func (m *xxx_SENSOR_Tx_SensorSonars) WakeUpChan() <-chan struct{} {
	return m.wakeUpChan
}

func (m *xxx_SENSOR_Tx_SensorSonars) Transmit(ctx context.Context) error {
	select {
	case m.transmitEventChan <- struct{}{}:
		return nil
	case <-ctx.Done():
		return fmt.Errorf("event-triggered transmit of SensorSonars: %w", ctx.Err())
	}
}

func (m *xxx_SENSOR_Tx_SensorSonars) TransmitEventChan() <-chan struct{} {
	return m.transmitEventChan
}

var _ canrunner.TransmittedMessage = &xxx_SENSOR_Tx_SensorSonars{}

// Nodes returns the example node descriptors.
func Nodes() *NodesDescriptor {
	return nd
}

// NodesDescriptor contains all example node descriptors.
type NodesDescriptor struct {
	DBG    *descriptor.Node
	DRIVER *descriptor.Node
	IO     *descriptor.Node
	MOTOR  *descriptor.Node
	SENSOR *descriptor.Node
}

// Messages returns the example message descriptors.
func Messages() *MessagesDescriptor {
	return md
}

// MessagesDescriptor contains all example message descriptors.
type MessagesDescriptor struct {
	EmptyMessage    *EmptyMessageDescriptor
	DriverHeartbeat *DriverHeartbeatDescriptor
	MotorCommand    *MotorCommandDescriptor
	SensorSonars    *SensorSonarsDescriptor
	MotorStatus     *MotorStatusDescriptor
	IODebug         *IODebugDescriptor
	IOFloat32       *IOFloat32Descriptor
}

// UnmarshalFrame unmarshals the provided example CAN frame.
func (md *MessagesDescriptor) UnmarshalFrame(f can.Frame) (generated.Message, error) {
	switch f.ID {
	case md.EmptyMessage.ID:
		var msg EmptyMessage
		if err := msg.UnmarshalFrame(f); err != nil {
			return nil, fmt.Errorf("unmarshal example frame: %w", err)
		}
		return &msg, nil
	case md.DriverHeartbeat.ID:
		var msg DriverHeartbeat
		if err := msg.UnmarshalFrame(f); err != nil {
			return nil, fmt.Errorf("unmarshal example frame: %w", err)
		}
		return &msg, nil
	case md.MotorCommand.ID:
		var msg MotorCommand
		if err := msg.UnmarshalFrame(f); err != nil {
			return nil, fmt.Errorf("unmarshal example frame: %w", err)
		}
		return &msg, nil
	case md.SensorSonars.ID:
		var msg SensorSonars
		if err := msg.UnmarshalFrame(f); err != nil {
			return nil, fmt.Errorf("unmarshal example frame: %w", err)
		}
		return &msg, nil
	case md.MotorStatus.ID:
		var msg MotorStatus
		if err := msg.UnmarshalFrame(f); err != nil {
			return nil, fmt.Errorf("unmarshal example frame: %w", err)
		}
		return &msg, nil
	case md.IODebug.ID:
		var msg IODebug
		if err := msg.UnmarshalFrame(f); err != nil {
			return nil, fmt.Errorf("unmarshal example frame: %w", err)
		}
		return &msg, nil
	case md.IOFloat32.ID:
		var msg IOFloat32
		if err := msg.UnmarshalFrame(f); err != nil {
			return nil, fmt.Errorf("unmarshal example frame: %w", err)
		}
		return &msg, nil
	default:
		return nil, fmt.Errorf("unmarshal example frame: ID not in database: %d", f.ID)
	}
}

type EmptyMessageDescriptor struct {
	*descriptor.Message
}

type DriverHeartbeatDescriptor struct {
	*descriptor.Message
	Command *descriptor.Signal
}

type MotorCommandDescriptor struct {
	*descriptor.Message
	Steer *descriptor.Signal
	Drive *descriptor.Signal
}

type SensorSonarsDescriptor struct {
	*descriptor.Message
	Mux          *descriptor.Signal
	ErrCount     *descriptor.Signal
	Left         *descriptor.Signal
	NoFiltLeft   *descriptor.Signal
	Middle       *descriptor.Signal
	NoFiltMiddle *descriptor.Signal
	Right        *descriptor.Signal
	NoFiltRight  *descriptor.Signal
	Rear         *descriptor.Signal
	NoFiltRear   *descriptor.Signal
}

type MotorStatusDescriptor struct {
	*descriptor.Message
	WheelError *descriptor.Signal
	SpeedKph   *descriptor.Signal
}

type IODebugDescriptor struct {
	*descriptor.Message
	TestUnsigned   *descriptor.Signal
	TestEnum       *descriptor.Signal
	TestSigned     *descriptor.Signal
	TestFloat      *descriptor.Signal
	TestBoolEnum   *descriptor.Signal
	TestScaledEnum *descriptor.Signal
}

type IOFloat32Descriptor struct {
	*descriptor.Message
	Float32ValueNoRange *descriptor.Signal
	Float32WithRange    *descriptor.Signal
}

// Database returns the example database descriptor.
func (md *MessagesDescriptor) Database() *descriptor.Database {
	return d
}

var nd = &NodesDescriptor{
	DBG:    d.Nodes[0],
	DRIVER: d.Nodes[1],
	IO:     d.Nodes[2],
	MOTOR:  d.Nodes[3],
	SENSOR: d.Nodes[4],
}

var md = &MessagesDescriptor{
	EmptyMessage: &EmptyMessageDescriptor{
		Message: d.Messages[0],
	},
	DriverHeartbeat: &DriverHeartbeatDescriptor{
		Message: d.Messages[1],
		Command: d.Messages[1].Signals[0],
	},
	MotorCommand: &MotorCommandDescriptor{
		Message: d.Messages[2],
		Steer:   d.Messages[2].Signals[0],
		Drive:   d.Messages[2].Signals[1],
	},
	SensorSonars: &SensorSonarsDescriptor{
		Message:      d.Messages[3],
		Mux:          d.Messages[3].Signals[0],
		ErrCount:     d.Messages[3].Signals[1],
		Left:         d.Messages[3].Signals[2],
		NoFiltLeft:   d.Messages[3].Signals[3],
		Middle:       d.Messages[3].Signals[4],
		NoFiltMiddle: d.Messages[3].Signals[5],
		Right:        d.Messages[3].Signals[6],
		NoFiltRight:  d.Messages[3].Signals[7],
		Rear:         d.Messages[3].Signals[8],
		NoFiltRear:   d.Messages[3].Signals[9],
	},
	MotorStatus: &MotorStatusDescriptor{
		Message:    d.Messages[4],
		WheelError: d.Messages[4].Signals[0],
		SpeedKph:   d.Messages[4].Signals[1],
	},
	IODebug: &IODebugDescriptor{
		Message:        d.Messages[5],
		TestUnsigned:   d.Messages[5].Signals[0],
		TestEnum:       d.Messages[5].Signals[1],
		TestSigned:     d.Messages[5].Signals[2],
		TestFloat:      d.Messages[5].Signals[3],
		TestBoolEnum:   d.Messages[5].Signals[4],
		TestScaledEnum: d.Messages[5].Signals[5],
	},
	IOFloat32: &IOFloat32Descriptor{
		Message:             d.Messages[6],
		Float32ValueNoRange: d.Messages[6].Signals[0],
		Float32WithRange:    d.Messages[6].Signals[1],
	},
}

var d = (*descriptor.Database)(&descriptor.Database{
	SourceFile: (string)("testdata/dbc/example/example.dbc"),
	Version:    (string)(""),
	Messages: ([]*descriptor.Message)([]*descriptor.Message{
		(*descriptor.Message)(&descriptor.Message{
			Name:        (string)("EmptyMessage"),
			ID:          (uint32)(1),
			IsExtended:  (bool)(false),
			Length:      (uint8)(0),
			SendType:    (descriptor.SendType)(0),
			Description: (string)(""),
			Signals:     ([]*descriptor.Signal)(nil),
			SenderNode:  (string)("DBG"),
			CycleTime:   (time.Duration)(0),
			DelayTime:   (time.Duration)(0),
		}),
		(*descriptor.Message)(&descriptor.Message{
			Name:        (string)("DriverHeartbeat"),
			ID:          (uint32)(100),
			IsExtended:  (bool)(false),
			Length:      (uint8)(1),
			SendType:    (descriptor.SendType)(1),
			Description: (string)("Sync message used to synchronize the controllers"),
			Signals: ([]*descriptor.Signal)([]*descriptor.Signal{
				(*descriptor.Signal)(&descriptor.Signal{
					Name:             (string)("Command"),
					Start:            (uint8)(0),
					Length:           (uint8)(8),
					IsBigEndian:      (bool)(false),
					IsSigned:         (bool)(false),
					IsFloat:          (bool)(false),
					IsMultiplexer:    (bool)(false),
					IsMultiplexed:    (bool)(false),
					MultiplexerValue: (uint)(0),
					Offset:           (float64)(0),
					Scale:            (float64)(1),
					Min:              (float64)(0),
					Max:              (float64)(0),
					Unit:             (string)(""),
					Description:      (string)(""),
					ValueDescriptions: ([]*descriptor.ValueDescription)([]*descriptor.ValueDescription{
						(*descriptor.ValueDescription)(&descriptor.ValueDescription{
							Value:       (int64)(0),
							Description: (string)("None"),
						}),
						(*descriptor.ValueDescription)(&descriptor.ValueDescription{
							Value:       (int64)(1),
							Description: (string)("Sync"),
						}),
						(*descriptor.ValueDescription)(&descriptor.ValueDescription{
							Value:       (int64)(2),
							Description: (string)("Reboot"),
						}),
						(*descriptor.ValueDescription)(&descriptor.ValueDescription{
							Value:       (int64)(3),
							Description: (string)("Headlights On"),
						}),
					}),
					ReceiverNodes: ([]string)([]string{
						(string)("SENSOR"),
						(string)("MOTOR"),
					}),
					DefaultValue: (int)(0),
				}),
			}),
			SenderNode: (string)("DRIVER"),
			CycleTime:  (time.Duration)(1000000000),
			DelayTime:  (time.Duration)(0),
		}),
		(*descriptor.Message)(&descriptor.Message{
			Name:        (string)("MotorCommand"),
			ID:          (uint32)(101),
			IsExtended:  (bool)(false),
			Length:      (uint8)(1),
			SendType:    (descriptor.SendType)(1),
			Description: (string)(""),
			Signals: ([]*descriptor.Signal)([]*descriptor.Signal{
				(*descriptor.Signal)(&descriptor.Signal{
					Name:              (string)("Steer"),
					Start:             (uint8)(0),
					Length:            (uint8)(4),
					IsBigEndian:       (bool)(false),
					IsSigned:          (bool)(true),
					IsFloat:           (bool)(false),
					IsMultiplexer:     (bool)(false),
					IsMultiplexed:     (bool)(false),
					MultiplexerValue:  (uint)(0),
					Offset:            (float64)(-5),
					Scale:             (float64)(1),
					Min:               (float64)(-5),
					Max:               (float64)(5),
					Unit:              (string)(""),
					Description:       (string)(""),
					ValueDescriptions: ([]*descriptor.ValueDescription)(nil),
					ReceiverNodes: ([]string)([]string{
						(string)("MOTOR"),
					}),
					DefaultValue: (int)(0),
				}),
				(*descriptor.Signal)(&descriptor.Signal{
					Name:              (string)("Drive"),
					Start:             (uint8)(4),
					Length:            (uint8)(4),
					IsBigEndian:       (bool)(false),
					IsSigned:          (bool)(false),
					IsFloat:           (bool)(false),
					IsMultiplexer:     (bool)(false),
					IsMultiplexed:     (bool)(false),
					MultiplexerValue:  (uint)(0),
					Offset:            (float64)(0),
					Scale:             (float64)(1),
					Min:               (float64)(0),
					Max:               (float64)(9),
					Unit:              (string)(""),
					Description:       (string)(""),
					ValueDescriptions: ([]*descriptor.ValueDescription)(nil),
					ReceiverNodes: ([]string)([]string{
						(string)("MOTOR"),
					}),
					DefaultValue: (int)(0),
				}),
			}),
			SenderNode: (string)("DRIVER"),
			CycleTime:  (time.Duration)(100000000),
			DelayTime:  (time.Duration)(0),
		}),
		(*descriptor.Message)(&descriptor.Message{
			Name:        (string)("SensorSonars"),
			ID:          (uint32)(200),
			IsExtended:  (bool)(false),
			Length:      (uint8)(8),
			SendType:    (descriptor.SendType)(1),
			Description: (string)(""),
			Signals: ([]*descriptor.Signal)([]*descriptor.Signal{
				(*descriptor.Signal)(&descriptor.Signal{
					Name:              (string)("Mux"),
					Start:             (uint8)(0),
					Length:            (uint8)(4),
					IsBigEndian:       (bool)(false),
					IsSigned:          (bool)(false),
					IsFloat:           (bool)(false),
					IsMultiplexer:     (bool)(true),
					IsMultiplexed:     (bool)(false),
					MultiplexerValue:  (uint)(0),
					Offset:            (float64)(0),
					Scale:             (float64)(1),
					Min:               (float64)(0),
					Max:               (float64)(0),
					Unit:              (string)(""),
					Description:       (string)(""),
					ValueDescriptions: ([]*descriptor.ValueDescription)(nil),
					ReceiverNodes: ([]string)([]string{
						(string)("DRIVER"),
						(string)("IO"),
					}),
					DefaultValue: (int)(0),
				}),
				(*descriptor.Signal)(&descriptor.Signal{
					Name:              (string)("ErrCount"),
					Start:             (uint8)(4),
					Length:            (uint8)(12),
					IsBigEndian:       (bool)(false),
					IsSigned:          (bool)(false),
					IsFloat:           (bool)(false),
					IsMultiplexer:     (bool)(false),
					IsMultiplexed:     (bool)(false),
					MultiplexerValue:  (uint)(0),
					Offset:            (float64)(0),
					Scale:             (float64)(1),
					Min:               (float64)(0),
					Max:               (float64)(0),
					Unit:              (string)(""),
					Description:       (string)(""),
					ValueDescriptions: ([]*descriptor.ValueDescription)(nil),
					ReceiverNodes: ([]string)([]string{
						(string)("DRIVER"),
						(string)("IO"),
					}),
					DefaultValue: (int)(0),
				}),
				(*descriptor.Signal)(&descriptor.Signal{
					Name:              (string)("Left"),
					Start:             (uint8)(16),
					Length:            (uint8)(12),
					IsBigEndian:       (bool)(false),
					IsSigned:          (bool)(false),
					IsFloat:           (bool)(false),
					IsMultiplexer:     (bool)(false),
					IsMultiplexed:     (bool)(true),
					MultiplexerValue:  (uint)(0),
					Offset:            (float64)(0),
					Scale:             (float64)(0.1),
					Min:               (float64)(0),
					Max:               (float64)(0),
					Unit:              (string)(""),
					Description:       (string)(""),
					ValueDescriptions: ([]*descriptor.ValueDescription)(nil),
					ReceiverNodes: ([]string)([]string{
						(string)("DRIVER"),
						(string)("IO"),
					}),
					DefaultValue: (int)(0),
				}),
				(*descriptor.Signal)(&descriptor.Signal{
					Name:              (string)("NoFiltLeft"),
					Start:             (uint8)(16),
					Length:            (uint8)(12),
					IsBigEndian:       (bool)(false),
					IsSigned:          (bool)(false),
					IsFloat:           (bool)(false),
					IsMultiplexer:     (bool)(false),
					IsMultiplexed:     (bool)(true),
					MultiplexerValue:  (uint)(1),
					Offset:            (float64)(0),
					Scale:             (float64)(0.1),
					Min:               (float64)(0),
					Max:               (float64)(0),
					Unit:              (string)(""),
					Description:       (string)(""),
					ValueDescriptions: ([]*descriptor.ValueDescription)(nil),
					ReceiverNodes: ([]string)([]string{
						(string)("DBG"),
					}),
					DefaultValue: (int)(0),
				}),
				(*descriptor.Signal)(&descriptor.Signal{
					Name:              (string)("Middle"),
					Start:             (uint8)(28),
					Length:            (uint8)(12),
					IsBigEndian:       (bool)(false),
					IsSigned:          (bool)(false),
					IsFloat:           (bool)(false),
					IsMultiplexer:     (bool)(false),
					IsMultiplexed:     (bool)(true),
					MultiplexerValue:  (uint)(0),
					Offset:            (float64)(0),
					Scale:             (float64)(0.1),
					Min:               (float64)(0),
					Max:               (float64)(0),
					Unit:              (string)(""),
					Description:       (string)(""),
					ValueDescriptions: ([]*descriptor.ValueDescription)(nil),
					ReceiverNodes: ([]string)([]string{
						(string)("DRIVER"),
						(string)("IO"),
					}),
					DefaultValue: (int)(0),
				}),
				(*descriptor.Signal)(&descriptor.Signal{
					Name:              (string)("NoFiltMiddle"),
					Start:             (uint8)(28),
					Length:            (uint8)(12),
					IsBigEndian:       (bool)(false),
					IsSigned:          (bool)(false),
					IsFloat:           (bool)(false),
					IsMultiplexer:     (bool)(false),
					IsMultiplexed:     (bool)(true),
					MultiplexerValue:  (uint)(1),
					Offset:            (float64)(0),
					Scale:             (float64)(0.1),
					Min:               (float64)(0),
					Max:               (float64)(0),
					Unit:              (string)(""),
					Description:       (string)(""),
					ValueDescriptions: ([]*descriptor.ValueDescription)(nil),
					ReceiverNodes: ([]string)([]string{
						(string)("DBG"),
					}),
					DefaultValue: (int)(0),
				}),
				(*descriptor.Signal)(&descriptor.Signal{
					Name:              (string)("Right"),
					Start:             (uint8)(40),
					Length:            (uint8)(12),
					IsBigEndian:       (bool)(false),
					IsSigned:          (bool)(false),
					IsFloat:           (bool)(false),
					IsMultiplexer:     (bool)(false),
					IsMultiplexed:     (bool)(true),
					MultiplexerValue:  (uint)(0),
					Offset:            (float64)(0),
					Scale:             (float64)(0.1),
					Min:               (float64)(0),
					Max:               (float64)(0),
					Unit:              (string)(""),
					Description:       (string)(""),
					ValueDescriptions: ([]*descriptor.ValueDescription)(nil),
					ReceiverNodes: ([]string)([]string{
						(string)("DRIVER"),
						(string)("IO"),
					}),
					DefaultValue: (int)(0),
				}),
				(*descriptor.Signal)(&descriptor.Signal{
					Name:              (string)("NoFiltRight"),
					Start:             (uint8)(40),
					Length:            (uint8)(12),
					IsBigEndian:       (bool)(false),
					IsSigned:          (bool)(false),
					IsFloat:           (bool)(false),
					IsMultiplexer:     (bool)(false),
					IsMultiplexed:     (bool)(true),
					MultiplexerValue:  (uint)(1),
					Offset:            (float64)(0),
					Scale:             (float64)(0.1),
					Min:               (float64)(0),
					Max:               (float64)(0),
					Unit:              (string)(""),
					Description:       (string)(""),
					ValueDescriptions: ([]*descriptor.ValueDescription)(nil),
					ReceiverNodes: ([]string)([]string{
						(string)("DBG"),
					}),
					DefaultValue: (int)(0),
				}),
				(*descriptor.Signal)(&descriptor.Signal{
					Name:              (string)("Rear"),
					Start:             (uint8)(52),
					Length:            (uint8)(12),
					IsBigEndian:       (bool)(false),
					IsSigned:          (bool)(false),
					IsFloat:           (bool)(false),
					IsMultiplexer:     (bool)(false),
					IsMultiplexed:     (bool)(true),
					MultiplexerValue:  (uint)(0),
					Offset:            (float64)(0),
					Scale:             (float64)(0.1),
					Min:               (float64)(0),
					Max:               (float64)(0),
					Unit:              (string)(""),
					Description:       (string)(""),
					ValueDescriptions: ([]*descriptor.ValueDescription)(nil),
					ReceiverNodes: ([]string)([]string{
						(string)("DRIVER"),
						(string)("IO"),
					}),
					DefaultValue: (int)(0),
				}),
				(*descriptor.Signal)(&descriptor.Signal{
					Name:              (string)("NoFiltRear"),
					Start:             (uint8)(52),
					Length:            (uint8)(12),
					IsBigEndian:       (bool)(false),
					IsSigned:          (bool)(false),
					IsFloat:           (bool)(false),
					IsMultiplexer:     (bool)(false),
					IsMultiplexed:     (bool)(true),
					MultiplexerValue:  (uint)(1),
					Offset:            (float64)(0),
					Scale:             (float64)(0.1),
					Min:               (float64)(0),
					Max:               (float64)(0),
					Unit:              (string)(""),
					Description:       (string)(""),
					ValueDescriptions: ([]*descriptor.ValueDescription)(nil),
					ReceiverNodes: ([]string)([]string{
						(string)("DBG"),
					}),
					DefaultValue: (int)(0),
				}),
			}),
			SenderNode: (string)("SENSOR"),
			CycleTime:  (time.Duration)(100000000),
			DelayTime:  (time.Duration)(0),
		}),
		(*descriptor.Message)(&descriptor.Message{
			Name:        (string)("MotorStatus"),
			ID:          (uint32)(400),
			IsExtended:  (bool)(false),
			Length:      (uint8)(3),
			SendType:    (descriptor.SendType)(1),
			Description: (string)(""),
			Signals: ([]*descriptor.Signal)([]*descriptor.Signal{
				(*descriptor.Signal)(&descriptor.Signal{
					Name:              (string)("WheelError"),
					Start:             (uint8)(0),
					Length:            (uint8)(1),
					IsBigEndian:       (bool)(false),
					IsSigned:          (bool)(false),
					IsFloat:           (bool)(false),
					IsMultiplexer:     (bool)(false),
					IsMultiplexed:     (bool)(false),
					MultiplexerValue:  (uint)(0),
					Offset:            (float64)(0),
					Scale:             (float64)(1),
					Min:               (float64)(0),
					Max:               (float64)(0),
					Unit:              (string)(""),
					Description:       (string)(""),
					ValueDescriptions: ([]*descriptor.ValueDescription)(nil),
					ReceiverNodes: ([]string)([]string{
						(string)("DRIVER"),
						(string)("IO"),
					}),
					DefaultValue: (int)(0),
				}),
				(*descriptor.Signal)(&descriptor.Signal{
					Name:              (string)("SpeedKph"),
					Start:             (uint8)(8),
					Length:            (uint8)(16),
					IsBigEndian:       (bool)(false),
					IsSigned:          (bool)(false),
					IsFloat:           (bool)(false),
					IsMultiplexer:     (bool)(false),
					IsMultiplexed:     (bool)(false),
					MultiplexerValue:  (uint)(0),
					Offset:            (float64)(0),
					Scale:             (float64)(0.001),
					Min:               (float64)(0),
					Max:               (float64)(0),
					Unit:              (string)("km/h"),
					Description:       (string)(""),
					ValueDescriptions: ([]*descriptor.ValueDescription)(nil),
					ReceiverNodes: ([]string)([]string{
						(string)("DRIVER"),
						(string)("IO"),
					}),
					DefaultValue: (int)(0),
				}),
			}),
			SenderNode: (string)("MOTOR"),
			CycleTime:  (time.Duration)(100000000),
			DelayTime:  (time.Duration)(0),
		}),
		(*descriptor.Message)(&descriptor.Message{
			Name:        (string)("IODebug"),
			ID:          (uint32)(500),
			IsExtended:  (bool)(false),
			Length:      (uint8)(6),
			SendType:    (descriptor.SendType)(2),
			Description: (string)(""),
			Signals: ([]*descriptor.Signal)([]*descriptor.Signal{
				(*descriptor.Signal)(&descriptor.Signal{
					Name:              (string)("TestUnsigned"),
					Start:             (uint8)(0),
					Length:            (uint8)(8),
					IsBigEndian:       (bool)(false),
					IsSigned:          (bool)(false),
					IsFloat:           (bool)(false),
					IsMultiplexer:     (bool)(false),
					IsMultiplexed:     (bool)(false),
					MultiplexerValue:  (uint)(0),
					Offset:            (float64)(0),
					Scale:             (float64)(1),
					Min:               (float64)(0),
					Max:               (float64)(0),
					Unit:              (string)(""),
					Description:       (string)(""),
					ValueDescriptions: ([]*descriptor.ValueDescription)(nil),
					ReceiverNodes: ([]string)([]string{
						(string)("DBG"),
					}),
					DefaultValue: (int)(0),
				}),
				(*descriptor.Signal)(&descriptor.Signal{
					Name:             (string)("TestEnum"),
					Start:            (uint8)(8),
					Length:           (uint8)(6),
					IsBigEndian:      (bool)(false),
					IsSigned:         (bool)(false),
					IsFloat:          (bool)(false),
					IsMultiplexer:    (bool)(false),
					IsMultiplexed:    (bool)(false),
					MultiplexerValue: (uint)(0),
					Offset:           (float64)(0),
					Scale:            (float64)(1),
					Min:              (float64)(0),
					Max:              (float64)(0),
					Unit:             (string)(""),
					Description:      (string)(""),
					ValueDescriptions: ([]*descriptor.ValueDescription)([]*descriptor.ValueDescription{
						(*descriptor.ValueDescription)(&descriptor.ValueDescription{
							Value:       (int64)(1),
							Description: (string)("One"),
						}),
						(*descriptor.ValueDescription)(&descriptor.ValueDescription{
							Value:       (int64)(2),
							Description: (string)("Two"),
						}),
					}),
					ReceiverNodes: ([]string)([]string{
						(string)("DBG"),
					}),
					DefaultValue: (int)(2),
				}),
				(*descriptor.Signal)(&descriptor.Signal{
					Name:              (string)("TestSigned"),
					Start:             (uint8)(16),
					Length:            (uint8)(8),
					IsBigEndian:       (bool)(false),
					IsSigned:          (bool)(true),
					IsFloat:           (bool)(false),
					IsMultiplexer:     (bool)(false),
					IsMultiplexed:     (bool)(false),
					MultiplexerValue:  (uint)(0),
					Offset:            (float64)(0),
					Scale:             (float64)(1),
					Min:               (float64)(0),
					Max:               (float64)(0),
					Unit:              (string)(""),
					Description:       (string)(""),
					ValueDescriptions: ([]*descriptor.ValueDescription)(nil),
					ReceiverNodes: ([]string)([]string{
						(string)("DBG"),
					}),
					DefaultValue: (int)(0),
				}),
				(*descriptor.Signal)(&descriptor.Signal{
					Name:              (string)("TestFloat"),
					Start:             (uint8)(24),
					Length:            (uint8)(8),
					IsBigEndian:       (bool)(false),
					IsSigned:          (bool)(false),
					IsFloat:           (bool)(false),
					IsMultiplexer:     (bool)(false),
					IsMultiplexed:     (bool)(false),
					MultiplexerValue:  (uint)(0),
					Offset:            (float64)(0),
					Scale:             (float64)(0.5),
					Min:               (float64)(0),
					Max:               (float64)(0),
					Unit:              (string)(""),
					Description:       (string)(""),
					ValueDescriptions: ([]*descriptor.ValueDescription)(nil),
					ReceiverNodes: ([]string)([]string{
						(string)("DBG"),
					}),
					DefaultValue: (int)(0),
				}),
				(*descriptor.Signal)(&descriptor.Signal{
					Name:             (string)("TestBoolEnum"),
					Start:            (uint8)(32),
					Length:           (uint8)(1),
					IsBigEndian:      (bool)(false),
					IsSigned:         (bool)(false),
					IsFloat:          (bool)(false),
					IsMultiplexer:    (bool)(false),
					IsMultiplexed:    (bool)(false),
					MultiplexerValue: (uint)(0),
					Offset:           (float64)(0),
					Scale:            (float64)(1),
					Min:              (float64)(0),
					Max:              (float64)(0),
					Unit:             (string)(""),
					Description:      (string)(""),
					ValueDescriptions: ([]*descriptor.ValueDescription)([]*descriptor.ValueDescription{
						(*descriptor.ValueDescription)(&descriptor.ValueDescription{
							Value:       (int64)(0),
							Description: (string)("Zero"),
						}),
						(*descriptor.ValueDescription)(&descriptor.ValueDescription{
							Value:       (int64)(1),
							Description: (string)("One"),
						}),
					}),
					ReceiverNodes: ([]string)([]string{
						(string)("DBG"),
					}),
					DefaultValue: (int)(0),
				}),
				(*descriptor.Signal)(&descriptor.Signal{
					Name:             (string)("TestScaledEnum"),
					Start:            (uint8)(40),
					Length:           (uint8)(2),
					IsBigEndian:      (bool)(false),
					IsSigned:         (bool)(false),
					IsFloat:          (bool)(false),
					IsMultiplexer:    (bool)(false),
					IsMultiplexed:    (bool)(false),
					MultiplexerValue: (uint)(0),
					Offset:           (float64)(0),
					Scale:            (float64)(2),
					Min:              (float64)(0),
					Max:              (float64)(6),
					Unit:             (string)(""),
					Description:      (string)(""),
					ValueDescriptions: ([]*descriptor.ValueDescription)([]*descriptor.ValueDescription{
						(*descriptor.ValueDescription)(&descriptor.ValueDescription{
							Value:       (int64)(0),
							Description: (string)("Zero"),
						}),
						(*descriptor.ValueDescription)(&descriptor.ValueDescription{
							Value:       (int64)(1),
							Description: (string)("Two"),
						}),
						(*descriptor.ValueDescription)(&descriptor.ValueDescription{
							Value:       (int64)(2),
							Description: (string)("Four"),
						}),
						(*descriptor.ValueDescription)(&descriptor.ValueDescription{
							Value:       (int64)(3),
							Description: (string)("Six"),
						}),
					}),
					ReceiverNodes: ([]string)([]string{
						(string)("DBG"),
					}),
					DefaultValue: (int)(0),
				}),
			}),
			SenderNode: (string)("IO"),
			CycleTime:  (time.Duration)(0),
			DelayTime:  (time.Duration)(0),
		}),
		(*descriptor.Message)(&descriptor.Message{
			Name:        (string)("IOFloat32"),
			ID:          (uint32)(600),
			IsExtended:  (bool)(false),
			Length:      (uint8)(8),
			SendType:    (descriptor.SendType)(0),
			Description: (string)(""),
			Signals: ([]*descriptor.Signal)([]*descriptor.Signal{
				(*descriptor.Signal)(&descriptor.Signal{
					Name:              (string)("Float32ValueNoRange"),
					Start:             (uint8)(0),
					Length:            (uint8)(32),
					IsBigEndian:       (bool)(false),
					IsSigned:          (bool)(true),
					IsFloat:           (bool)(true),
					IsMultiplexer:     (bool)(false),
					IsMultiplexed:     (bool)(false),
					MultiplexerValue:  (uint)(0),
					Offset:            (float64)(0),
					Scale:             (float64)(1),
					Min:               (float64)(0),
					Max:               (float64)(0),
					Unit:              (string)(""),
					Description:       (string)(""),
					ValueDescriptions: ([]*descriptor.ValueDescription)(nil),
					ReceiverNodes: ([]string)([]string{
						(string)("DBG"),
					}),
					DefaultValue: (int)(0),
				}),
				(*descriptor.Signal)(&descriptor.Signal{
					Name:              (string)("Float32WithRange"),
					Start:             (uint8)(32),
					Length:            (uint8)(32),
					IsBigEndian:       (bool)(false),
					IsSigned:          (bool)(true),
					IsFloat:           (bool)(true),
					IsMultiplexer:     (bool)(false),
					IsMultiplexed:     (bool)(false),
					MultiplexerValue:  (uint)(0),
					Offset:            (float64)(0),
					Scale:             (float64)(1),
					Min:               (float64)(-100),
					Max:               (float64)(100),
					Unit:              (string)(""),
					Description:       (string)(""),
					ValueDescriptions: ([]*descriptor.ValueDescription)(nil),
					ReceiverNodes: ([]string)([]string{
						(string)("DBG"),
					}),
					DefaultValue: (int)(0),
				}),
			}),
			SenderNode: (string)("IO"),
			CycleTime:  (time.Duration)(0),
			DelayTime:  (time.Duration)(0),
		}),
	}),
	Nodes: ([]*descriptor.Node)([]*descriptor.Node{
		(*descriptor.Node)(&descriptor.Node{
			Name:        (string)("DBG"),
			Description: (string)(""),
		}),
		(*descriptor.Node)(&descriptor.Node{
			Name:        (string)("DRIVER"),
			Description: (string)("The driver controller driving the car"),
		}),
		(*descriptor.Node)(&descriptor.Node{
			Name:        (string)("IO"),
			Description: (string)(""),
		}),
		(*descriptor.Node)(&descriptor.Node{
			Name:        (string)("MOTOR"),
			Description: (string)("The motor controller of the car"),
		}),
		(*descriptor.Node)(&descriptor.Node{
			Name:        (string)("SENSOR"),
			Description: (string)("The sensor controller of the car"),
		}),
	}),
})
