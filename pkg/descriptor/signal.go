package descriptor

import (
	"math"

	"github.com/blueinnovationsgroup/can-go"
)

// Signal describes a CAN signal.
type Signal struct {
	// Description of the signal.
	Name string
	// Start bit.
	Start uint8
	// Length in bits.
	Length uint8
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
func (s *Signal) ValueDescription(value int64) (string, bool) {
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
	case s.Length == 1:
		if d.Bit(s.Start) {
			return 1
		}
		return 0
	case s.IsSigned:
		var value int64
		if s.IsBigEndian {
			value = d.SignedBitsBigEndian(s.Start, s.Length)
		} else {
			value = d.SignedBitsLittleEndian(s.Start, s.Length)
		}
		return s.ToPhysical(float64(value))
	default:
		var value uint64
		if s.IsBigEndian {
			value = d.UnsignedBitsBigEndian(s.Start, s.Length)
		} else {
			value = d.UnsignedBitsLittleEndian(s.Start, s.Length)
		}
		return s.ToPhysical(float64(value))
	}
}

// UnmarshalUnsigned returns the unsigned value of the signal in the provided CAN frame.
func (s *Signal) UnmarshalUnsigned(d can.Data) uint64 {
	if s.IsBigEndian {
		return d.UnsignedBitsBigEndian(s.Start, s.Length)
	}
	return d.UnsignedBitsLittleEndian(s.Start, s.Length)
}

// UnmarshalValueDescription returns the value description of the signal in the provided CAN data.
func (s *Signal) UnmarshalValueDescription(d can.Data) (string, bool) {
	if len(s.ValueDescriptions) == 0 {
		return "", false
	}
	var intValue int64
	if s.IsSigned {
		intValue = s.UnmarshalSigned(d)
	} else {
		intValue = int64(s.UnmarshalUnsigned(d))
	}
	return s.ValueDescription(intValue)
}

// UnmarshalSigned returns the signed value of the signal in the provided CAN frame.
func (s *Signal) UnmarshalSigned(d can.Data) int64 {
	if s.IsBigEndian {
		return d.SignedBitsBigEndian(s.Start, s.Length)
	}
	return d.SignedBitsLittleEndian(s.Start, s.Length)
}

// UnmarshalBool returns the bool value of the signal in the provided CAN frame.
func (s *Signal) UnmarshalBool(d can.Data) bool {
	return d.Bit(s.Start)
}

// MarshalUnsigned sets the unsigned value of the signal in the provided CAN frame.
func (s *Signal) MarshalUnsigned(d *can.Data, value uint64) {
	if s.IsBigEndian {
		d.SetUnsignedBitsBigEndian(s.Start, s.Length, value)
	} else {
		d.SetUnsignedBitsLittleEndian(s.Start, s.Length, value)
	}
}

// MarshalSigned sets the signed value of the signal in the provided CAN frame.
func (s *Signal) MarshalSigned(d *can.Data, value int64) {
	if s.IsBigEndian {
		d.SetSignedBitsBigEndian(s.Start, s.Length, value)
	} else {
		d.SetSignedBitsLittleEndian(s.Start, s.Length, value)
	}
}

// MarshalBool sets the bool value of the signal in the provided CAN frame.
func (s *Signal) MarshalBool(d *can.Data, value bool) {
	d.SetBit(s.Start, value)
}

// MaxUnsigned returns the maximum unsigned value representable by the signal.
func (s *Signal) MaxUnsigned() uint64 {
	return (2 << (s.Length - 1)) - 1
}

// MinSigned returns the minimum signed value representable by the signal.
func (s *Signal) MinSigned() int64 {
	return (2 << (s.Length - 1) / 2) * -1
}

// MaxSigned returns the maximum signed value representable by the signal.
func (s *Signal) MaxSigned() int64 {
	return (2 << (s.Length - 1) / 2) - 1
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
