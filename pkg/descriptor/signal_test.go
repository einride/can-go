package descriptor

import (
	"encoding/hex"
	"math"
	"testing"

	"go.einride.tech/can"
	"gotest.tools/v3/assert"
)

func TestSignal_Decode_UnsignedBigEndian(t *testing.T) {
	s := &Signal{
		Name:        "TestSignal",
		IsSigned:    false,
		IsBigEndian: true,
		Offset:      -1,
		Scale:       3.0517578125e-05,
		Length:      10,
		Start:       32,
		Min:         0,
		Max:         1,
	}
	const value uint64 = 180

	// Testing can.Data
	var data can.Data
	data.SetUnsignedBitsBigEndian(uint8(s.Start), uint8(s.Length), value)
	actual := s.Decode(data)
	assert.DeepEqual(t, s.Offset+float64(value)*s.Scale, actual)

	// Testing payload
	p, _ := can.PayloadFromHex(hex.EncodeToString(data[:]))
	actual = s.DecodePayload(&p)
	assert.DeepEqual(t, s.Offset+float64(value)*s.Scale, actual)
}

func TestSignal_Decode_SignedBigEndian(t *testing.T) {
	s := &Signal{
		Name:        "TestSignal",
		IsSigned:    true,
		IsBigEndian: true,
		Offset:      -1,
		Scale:       3.0517578125e-05,
		Length:      10,
		Start:       32,
		Min:         -1,
		Max:         1,
	}
	const value int64 = -180

	// Testing can.Data
	var data can.Data
	data.SetSignedBitsBigEndian(uint8(s.Start), uint8(s.Length), value)
	actual := s.Decode(data)
	assert.DeepEqual(t, s.Offset+float64(value)*s.Scale, actual)

	// Testing payload
	p, _ := can.PayloadFromHex(hex.EncodeToString(data[:]))
	actual = s.DecodePayload(&p)
	assert.DeepEqual(t, s.Offset+float64(value)*s.Scale, actual)
}

func TestSignal_Decode_UnsignedLittleEndian(t *testing.T) {
	s := &Signal{
		Name:        "TestSignal",
		IsSigned:    false,
		IsBigEndian: false,
		Offset:      -1,
		Scale:       3.0517578125e-05,
		Length:      10,
		Start:       32,
		Min:         0,
		Max:         1,
	}
	const value uint64 = 180

	// Testing can.Data
	var data can.Data
	data.SetUnsignedBitsLittleEndian(uint8(s.Start), uint8(s.Length), value)
	actual := s.Decode(data)
	assert.DeepEqual(t, s.Offset+float64(value)*s.Scale, actual)

	// Testing payload
	p, _ := can.PayloadFromHex(hex.EncodeToString(data[:]))
	actual = s.DecodePayload(&p)
	assert.DeepEqual(t, s.Offset+float64(value)*s.Scale, actual)
}

func TestSignal_Decode_SignedLittleEndian(t *testing.T) {
	s := &Signal{
		Name:        "TestSignal",
		IsSigned:    true,
		IsBigEndian: false,
		Offset:      -1,
		Scale:       3.0517578125e-05,
		Length:      10,
		Start:       32,
		Min:         -1,
		Max:         1,
	}
	const value int64 = -180

	// Testing can.Data
	var data can.Data
	data.SetSignedBitsLittleEndian(uint8(s.Start), uint8(s.Length), value)
	actual := s.Decode(data)
	assert.DeepEqual(t, s.Offset+float64(value)*s.Scale, actual)

	// Testing payload
	p, _ := can.PayloadFromHex(hex.EncodeToString(data[:]))
	actual = s.DecodePayload(&p)
	assert.DeepEqual(t, s.Offset+float64(value)*s.Scale, actual)
}

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
	data.SetSignedBitsBigEndian(uint8(s.Start), uint8(s.Length), value)
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
	expected.SetUnsignedBitsBigEndian(uint8(s.Start), uint8(s.Length), value)
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
	expected.SetSignedBitsBigEndian(uint8(s.Start), uint8(s.Length), value)
	var actual can.Data
	s.MarshalSigned(&actual, value)
	assert.DeepEqual(t, expected, actual)
}
