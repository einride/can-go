package descriptor

import (
	"math"
	"testing"

	"github.com/blueinnovationsgroup/can-go"
	"gotest.tools/v3/assert"
)

func TestSignal_FromPhysical_SaturatedCast(t *testing.T) {
	s := &Signal{
		Name:   "TestSignal",
		Offset: -1,
		Scale:  3.0517578125e-05,
		Min:    -1,
		Max:    1,
		Length: 16,
	}
	// without a saturated cast, the result would be math.MaxUint16 + 1, which would wrap around to 0
	assert.Equal(t, uint16(math.MaxUint16), uint16(s.FromPhysical(180)))
}

func TestSignal_SaturatedCastSigned(t *testing.T) {
	s := &Signal{
		Name:     "TestSignal",
		IsSigned: true,
		Length:   6,
	}
	assert.Equal(t, int64(31), s.SaturatedCastSigned(254))
	assert.Equal(t, int64(-32), s.SaturatedCastSigned(-255))
}

func TestSignal_SaturatedCastUnsigned(t *testing.T) {
	s := &Signal{
		Name:   "TestSignal",
		Length: 6,
	}
	assert.Equal(t, uint64(63), s.SaturatedCastUnsigned(255))
}

func TestSignal_UnmarshalSigned_BigEndian(t *testing.T) {
	s := &Signal{
		Name:        "TestSignal",
		IsSigned:    true,
		IsBigEndian: true,
		Length:      8,
		Start:       32,
	}
	const value int64 = -8
	var data can.Data
	data.SetSignedBitsBigEndian(s.Start, s.Length, value)
	assert.Equal(t, value, s.UnmarshalSigned(data))
}

func TestSignal_MarshalUnsigned_BigEndian(t *testing.T) {
	s := &Signal{
		Name:        "TestSignal",
		IsBigEndian: true,
		Length:      8,
		Start:       32,
	}
	const value uint64 = 8
	var expected can.Data
	expected.SetUnsignedBitsBigEndian(s.Start, s.Length, value)
	var actual can.Data
	s.MarshalUnsigned(&actual, value)
	assert.DeepEqual(t, expected, actual)
}

func TestSignal_MarshalSigned_BigEndian(t *testing.T) {
	s := &Signal{
		Name:        "TestSignal",
		IsSigned:    true,
		IsBigEndian: true,
		Length:      8,
		Start:       32,
	}
	const value int64 = -8
	var expected can.Data
	expected.SetSignedBitsBigEndian(s.Start, s.Length, value)
	var actual can.Data
	s.MarshalSigned(&actual, value)
	assert.DeepEqual(t, expected, actual)
}
