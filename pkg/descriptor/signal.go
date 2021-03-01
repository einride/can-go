package descriptor

import (
	"math"

	"go.einride.tech/can"
)

// Signal describes a CAN signal.
type Signal struct {
	// Description of the signal.
	Name string
	// Start bit.
	Start uint16
	// Length in bits.
	Length uint16
	// IsBigEndian is true if the signal is big-endian.
	IsBigEndian bool
	// IsSigned is true if the signal uses raw signed values.
	IsSigned bool
	// IsMultiplexer is true if the signal is the multiplexor of a multiplexed message.
	IsMultiplexer bool
	// IsMultiplexed is true if the signal is multiplexed.
	IsMultiplexed bool
	// MultiplexerValue is the value of the multiplexer when this signal is present.
	MultiplexerValue uint
	// Offset for real-world transform.
	Offset float64
	// Scale for real-world transform.
	Scale float64
	// Min real-world value.
	Min float64
	// Max real-world value.
	Max float64
	// Unit of the signal.
	Unit string
	// Description of the signal.
	Description string
	// ValueDescriptions of the signal.
	ValueDescriptions []*ValueDescription
	// ReceiverNodes is the list of names of the nodes receiving the signal.
	ReceiverNodes []string
	// DefaultValue of the signal.
	DefaultValue int
}

// ValueDescription returns the value description for the provided value.
func (s *Signal) ValueDescription(value int) (string, bool) {
	for _, vd := range s.ValueDescriptions {
		if vd.Value == value {
			return vd.Description, true
		}
	}
	return "", false
}

// ToPhysical converts a raw signal value to its physical value.
func (s *Signal) ToPhysical(value float64) float64 {
	result := value
	result *= s.Scale
	result += s.Offset
	if s.Min != 0 || s.Max != 0 {
		result = math.Max(math.Min(result, s.Max), s.Min)
	}
	return result
}

// FromPhysical converts a physical signal value to its raw value.
func (s *Signal) FromPhysical(physical float64) float64 {
	result := physical
	if s.Min != 0 || s.Max != 0 {
		result = math.Max(math.Min(result, s.Max), s.Min)
	}
	result -= s.Offset
	result /= s.Scale
	// perform saturated cast
	if s.IsSigned {
		result = math.Max(float64(s.MinSigned()), math.Min(float64(s.MaxSigned()), result))
	} else {
		result = math.Max(0, math.Min(float64(s.MaxUnsigned()), result))
	}
	return result
}

// UnmarshalPhysical returns the physical value of the signal in the provided CAN frame.
func (s *Signal) UnmarshalPhysical(d can.Data) float64 {
	switch {
	case uint8(s.Length) == 1:
		if d.Bit(uint8(s.Start)) {
			return 1
		}
		return 0
	case s.IsSigned:
		var value int64
		if s.IsBigEndian {
			value = d.SignedBitsBigEndian(uint8(s.Start), uint8(s.Length))
		} else {
			value = d.SignedBitsLittleEndian(uint8(s.Start), uint8(s.Length))
		}
		return s.ToPhysical(float64(value))
	default:
		var value uint64
		if s.IsBigEndian {
			value = d.UnsignedBitsBigEndian(uint8(s.Start), uint8(s.Length))
		} else {
			value = d.UnsignedBitsLittleEndian(uint8(s.Start), uint8(s.Length))
		}
		return s.ToPhysical(float64(value))
	}
}

// UnmarshalPhysicalPayload returns the physical value of the signal in the provided CAN frame.
func (s *Signal) UnmarshalPhysicalPayload(p *can.Payload) float64 {
	switch {
	case uint8(s.Length) == 1:
		if p.Bit(s.Start) {
			return 1
		}
		return 0
	case s.IsSigned:
		var value int64
		if s.IsBigEndian {
			value = p.SignedBitsBigEndian(s.Start, s.Length)
		} else {
			value = p.SignedBitsLittleEndian(s.Start, s.Length)
		}
		return s.ToPhysical(float64(value))
	default:
		var value uint64
		if s.IsBigEndian {
			value = p.UnsignedBitsBigEndian(s.Start, s.Length)
		} else {
			value = p.UnsignedBitsLittleEndian(s.Start, s.Length)
		}
		return s.ToPhysical(float64(value))
	}
}

// UnmarshalUnsigned returns the unsigned value of the signal in the provided CAN frame.
func (s *Signal) UnmarshalUnsigned(d can.Data) uint64 {
	if s.IsBigEndian {
		return d.UnsignedBitsBigEndian(uint8(s.Start), uint8(s.Length))
	}
	return d.UnsignedBitsLittleEndian(uint8(s.Start), uint8(s.Length))
}

// UnmarshalUnsignedPayload returns the unsigned value of the signal in the provided CAN frame.
func (s *Signal) UnmarshalUnsignedPayload(p *can.Payload) uint64 {
	if s.IsBigEndian {
		return p.UnsignedBitsBigEndian(s.Start, s.Length)
	}
	return p.UnsignedBitsLittleEndian(s.Start, s.Length)
}

// UnmarshalValueDescription returns the value description of the signal in the provided CAN data.
func (s *Signal) UnmarshalValueDescription(d can.Data) (string, bool) {
	if len(s.ValueDescriptions) == 0 {
		return "", false
	}
	var intValue int
	if s.IsSigned {
		intValue = int(s.UnmarshalSigned(d))
	} else {
		intValue = int(s.UnmarshalUnsigned(d))
	}
	return s.ValueDescription(intValue)
}

// UnmarshalValueDescriptionPayload returns the value description of the signal in the provided CAN data.
func (s *Signal) UnmarshalValueDescriptionPayload(p *can.Payload) (string, bool) {
	if len(s.ValueDescriptions) == 0 {
		return "", false
	}
	var intValue int
	if s.IsSigned {
		intValue = int(s.UnmarshalSignedPayload(p))
	} else {
		intValue = int(s.UnmarshalUnsignedPayload(p))
	}
	return s.ValueDescription(intValue)
}

// UnmarshalSigned returns the signed value of the signal in the provided CAN frame.
func (s *Signal) UnmarshalSigned(d can.Data) int64 {
	if s.IsBigEndian {
		return d.SignedBitsBigEndian(uint8(s.Start), uint8(uint8(s.Length)))
	}
	return d.SignedBitsLittleEndian(uint8(s.Start), uint8(uint8(s.Length)))
}

// UnmarshalSignedPayload returns the signed value of the signal in the provided CAN frame.
func (s *Signal) UnmarshalSignedPayload(p *can.Payload) int64 {
	if s.IsBigEndian {
		return p.SignedBitsBigEndian(s.Start, s.Length)
	}
	return p.SignedBitsLittleEndian(s.Start, s.Length)
}

// UnmarshalBool returns the bool value of the signal in the provided CAN frame.
func (s *Signal) UnmarshalBool(d can.Data) bool {
	return d.Bit(uint8(s.Start))
}

// MarshalUnsigned sets the unsigned value of the signal in the provided CAN frame.
func (s *Signal) MarshalUnsigned(d *can.Data, value uint64) {
	if s.IsBigEndian {
		d.SetUnsignedBitsBigEndian(uint8(s.Start), uint8(s.Length), value)
	} else {
		d.SetUnsignedBitsLittleEndian(uint8(s.Start), uint8(s.Length), value)
	}
}

// MarshalSigned sets the signed value of the signal in the provided CAN frame.
func (s *Signal) MarshalSigned(d *can.Data, value int64) {
	if s.IsBigEndian {
		d.SetSignedBitsBigEndian(uint8(s.Start), uint8(s.Length), value)
	} else {
		d.SetSignedBitsLittleEndian(uint8(s.Start), uint8(s.Length), value)
	}
}

// MarshalBool sets the bool value of the signal in the provided CAN frame.
func (s *Signal) MarshalBool(d *can.Data, value bool) {
	d.SetBit(uint8(s.Start), value)
}

// MaxUnsigned returns the maximum unsigned value representable by the signal.
func (s *Signal) MaxUnsigned() uint64 {
	return (2 << (uint8(s.Length) - 1)) - 1
}

// MinSigned returns the minimum signed value representable by the signal.
func (s *Signal) MinSigned() int64 {
	return (2 << (uint8(s.Length) - 1) / 2) * -1
}

// MaxSigned returns the maximum signed value representable by the signal.
func (s *Signal) MaxSigned() int64 {
	return (2 << (uint8(s.Length) - 1) / 2) - 1
}

// SaturatedCastSigned performs a saturated cast of an int64 to the value domain of the signal.
func (s *Signal) SaturatedCastSigned(value int64) int64 {
	min := s.MinSigned()
	max := s.MaxSigned()
	switch {
	case value < min:
		return min
	case value > max:
		return max
	default:
		return value
	}
}

// SaturatedCastUnsigned performs a saturated cast of a uint64 to the value domain of the signal.
func (s *Signal) SaturatedCastUnsigned(value uint64) uint64 {
	max := s.MaxUnsigned()
	if value > max {
		return max
	}
	return value
}
