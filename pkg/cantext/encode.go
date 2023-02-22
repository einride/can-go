package cantext

import (
	"strconv"

	"github.com/blueinnovationsgroup/can-go"
	"github.com/blueinnovationsgroup/can-go/pkg/descriptor"
	"github.com/blueinnovationsgroup/can-go/pkg/generated"
)

// preAllocatedBytesPerSignal is an estimate of how many bytes each signal needs.
const preAllocatedBytesPerSignal = 40

func MessageString(m generated.Message) string {
	return string(MarshalCompact(m))
}

func MarshalCompact(m generated.Message) []byte {
	f := m.Frame()
	buf := make([]byte, 0, len(m.Descriptor().Signals)*preAllocatedBytesPerSignal)
	buf = append(buf, "{"...)
	for i, s := range m.Descriptor().Signals {
		buf = AppendSignalCompact(buf, s, f.Data)
		if i != len(m.Descriptor().Signals)-1 {
			buf = append(buf, ", "...)
		}
	}
	buf = append(buf, "}"...)
	return buf
}

func Marshal(m generated.Message) []byte {
	f := m.Frame()
	// allocate space for one "extra" signal to account for message header
	buf := make([]byte, 0, (len(m.Descriptor().Signals)+1)*preAllocatedBytesPerSignal)
	buf = append(buf, m.Descriptor().Name...)
	for _, s := range m.Descriptor().Signals {
		buf = append(buf, "\n\t"...)
		buf = AppendSignal(buf, s, f.Data)
	}
	return buf
}

func AppendSignal(buf []byte, s *descriptor.Signal, d can.Data) []byte {
	buf = append(buf, s.Name...)
	buf = append(buf, ": "...)
	switch {
	case s.Length == 1: // bool
		val := s.UnmarshalBool(d)
		buf = strconv.AppendBool(buf, val)
	case s.IsSigned: // signed
		buf = strconv.AppendFloat(buf, s.UnmarshalPhysical(d), 'g', -1, 64)
		buf = append(buf, s.Unit...)
		buf = append(buf, " ("...)
		buf = append(buf, "0x"...)
		buf = strconv.AppendUint(buf, uint64(s.UnmarshalSigned(d)), 16)
		buf = append(buf, ')')
	default: // unsigned
		buf = strconv.AppendFloat(buf, s.UnmarshalPhysical(d), 'g', -1, 64)
		buf = append(buf, s.Unit...)
		buf = append(buf, " ("...)
		buf = append(buf, "0x"...)
		buf = strconv.AppendUint(buf, s.UnmarshalUnsigned(d), 16)
		buf = append(buf, ")"...)
	}
	if vd, ok := s.UnmarshalValueDescription(d); ok {
		buf = append(buf, ' ')
		buf = append(buf, vd...)
	}
	return buf
}

func AppendSignalCompact(buf []byte, s *descriptor.Signal, d can.Data) []byte {
	buf = append(buf, s.Name...)
	buf = append(buf, ": "...)
	valueDescription, hasValueDescription := s.UnmarshalValueDescription(d)
	switch {
	case hasValueDescription:
		buf = append(buf, valueDescription...)
	case s.Length == 1: // bool
		val := s.UnmarshalBool(d)
		buf = strconv.AppendBool(buf, val)
	case s.IsSigned: // signed
		buf = strconv.AppendFloat(buf, s.UnmarshalPhysical(d), 'g', -1, 64)
		buf = append(buf, s.Unit...)
	default: // unsigned
		buf = strconv.AppendFloat(buf, s.UnmarshalPhysical(d), 'g', -1, 64)
		buf = append(buf, s.Unit...)
	}
	return buf
}

func AppendID(buf []byte, m *descriptor.Message) []byte {
	buf = append(buf, "ID: "...)
	buf = strconv.AppendUint(buf, uint64(m.ID), 10)
	buf = append(buf, " (0x"...)
	buf = strconv.AppendUint(buf, uint64(m.ID), 16)
	buf = append(buf, ")"...)
	return buf
}

func AppendSender(buf []byte, m *descriptor.Message) []byte {
	return appendAttributeString(buf, "Sender", m.SenderNode)
}

func AppendSendType(buf []byte, m *descriptor.Message) []byte {
	return appendAttributeString(buf, "SendType", m.SendType.String())
}

func AppendCycleTime(buf []byte, m *descriptor.Message) []byte {
	return appendAttributeString(buf, "CycleTime", m.CycleTime.String())
}

func AppendDelayTime(buf []byte, m *descriptor.Message) []byte {
	return appendAttributeString(buf, "DelayTime", m.DelayTime.String())
}

func AppendFrame(buf []byte, f can.Frame) []byte {
	return appendAttributeString(buf, "Frame", f.String())
}

func appendAttributeString(buf []byte, name, s string) []byte {
	buf = append(buf, name...)
	buf = append(buf, ": "...)
	buf = append(buf, s...)
	return buf
}
