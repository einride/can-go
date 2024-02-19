package can

import (
	"fmt"
	"testing"
	"testing/quick"

	"gotest.tools/v3/assert"
)

func TestData_Bit(t *testing.T) {
	for i, tt := range []struct {
		data Data
		bits []struct {
			i   uint8
			bit bool
		}
	}{
		{
			data: Data{0x01, 0x23},
			bits: []struct {
				i   uint8
				bit bool
			}{
				// nibble 1: 0x1
				{bit: true, i: 0},
				{bit: false, i: 1},
				{bit: false, i: 2},
				{bit: false, i: 3},
				// nibble 2: 0x0
				{bit: false, i: 4},
				{bit: false, i: 5},
				{bit: false, i: 6},
				{bit: false, i: 7},
				// nibble 3: 0x3
				{bit: true, i: 8},
				{bit: true, i: 9},
				{bit: false, i: 10},
				{bit: false, i: 11},
				// nibble 4: 0x2
				{bit: false, i: 12},
				{bit: true, i: 13},
				{bit: false, i: 14},
				{bit: false, i: 15},
			},
		},
	} {
		i, tt := i, tt
		t.Run("Get", func(t *testing.T) {
			i, tt := i, tt
			for j, ttBit := range tt.bits {
				j, ttBit := j, ttBit
				t.Run(fmt.Sprintf("tt=%v,bit=%v", i, j), func(t *testing.T) {
					bit := tt.data.Bit(ttBit.i)
					assert.Equal(t, ttBit.bit, bit)
				})
			}
		})
		t.Run("Set", func(t *testing.T) {
			i, tt := i, tt
			t.Run(fmt.Sprintf("data=%v", i), func(t *testing.T) {
				var data Data
				for _, tBit := range tt.bits {
					data.SetBit(tBit.i, tBit.bit)
				}
				assert.DeepEqual(t, tt.data, data)
			})
		})
	}
}

func TestData_Property_SetGetBit(t *testing.T) {
	f := func(_ Data, _ uint8, bit bool) bool {
		return bit
	}
	g := func(data Data, i uint8, bit bool) bool {
		i %= 64
		data.SetBit(i, bit)
		return data.Bit(i)
	}
	assert.NilError(t, quick.CheckEqual(f, g, nil))
}

func TestData_LittleEndian(t *testing.T) {
	for i, tt := range []struct {
		data    Data
		signals []struct {
			start    uint8
			length   uint8
			unsigned uint64
			signed   int64
		}
	}{
		{
			data: Data{0x80, 0x01},
			signals: []struct {
				start    uint8
				length   uint8
				unsigned uint64
				signed   int64
			}{
				{start: 7, length: 2, unsigned: 0x3, signed: -1},
			},
		},
		{
			data: Data{0x01, 0x02, 0x03},
			signals: []struct {
				start    uint8
				length   uint8
				unsigned uint64
				signed   int64
			}{
				{start: 0, length: 24, unsigned: 0x030201, signed: 197121},
			},
		},
		{
			data: Data{0x40, 0x23, 0x01, 0x12},
			signals: []struct {
				start    uint8
				length   uint8
				unsigned uint64
				signed   int64
			}{
				{start: 24, length: 8, unsigned: 0x12, signed: 18},
				{start: 8, length: 8, unsigned: 0x23, signed: 35},
				{start: 4, length: 16, unsigned: 0x1234, signed: 4660},
			},
		},
	} {
		i, tt := i, tt
		t.Run(fmt.Sprintf("UnsignedBits:%v", i), func(t *testing.T) {
			for j, signal := range tt.signals {
				j, signal := j, signal
				t.Run(fmt.Sprintf("signal:%v", j), func(t *testing.T) {
					actual := tt.data.UnsignedBitsLittleEndian(signal.start, signal.length)
					assert.Equal(t, signal.unsigned, actual)
				})
			}
		})
		t.Run(fmt.Sprintf("SignedBits:%v", i), func(t *testing.T) {
			for j, signal := range tt.signals {
				j, signal := j, signal
				t.Run(fmt.Sprintf("signal:%v", j), func(t *testing.T) {
					actual := tt.data.SignedBitsLittleEndian(signal.start, signal.length)
					assert.Equal(t, signal.signed, actual)
				})
			}
		})
		t.Run(fmt.Sprintf("SetUnsignedBits:%v", i), func(t *testing.T) {
			var data Data
			for j, signal := range tt.signals {
				j, signal := j, signal
				t.Run(fmt.Sprintf("data:%v", j), func(_ *testing.T) {
					data.SetUnsignedBitsLittleEndian(signal.start, signal.length, signal.unsigned)
				})
			}
			assert.DeepEqual(t, tt.data, data)
		})
		t.Run(fmt.Sprintf("SetSignedBits:%v", i), func(t *testing.T) {
			var data Data
			for j, signal := range tt.signals {
				j, signal := j, signal
				t.Run(fmt.Sprintf("data:%v", j), func(_ *testing.T) {
					data.SetSignedBitsLittleEndian(signal.start, signal.length, signal.signed)
				})
			}
			assert.DeepEqual(t, tt.data, data)
		})
	}
}

func TestData_BigEndian(t *testing.T) {
	for i, tt := range []struct {
		data    Data
		signals []struct {
			start    uint8
			length   uint8
			unsigned uint64
			signed   int64
		}
	}{
		{
			data: Data{0x3f, 0xf7, 0x0d, 0xc4, 0x0c, 0x93, 0xff, 0xff},
			signals: []struct {
				start    uint8
				length   uint8
				unsigned uint64
				signed   int64
			}{
				{start: 7, length: 3, unsigned: 0x1, signed: 1},
				{start: 4, length: 1, unsigned: 0x1, signed: -1},
				{start: 55, length: 16, unsigned: 0xffff, signed: -1},
				{start: 39, length: 16, unsigned: 0xc93, signed: 3219},
				{start: 23, length: 16, unsigned: 0xdc4, signed: 3524},
				{start: 3, length: 12, unsigned: 0xff7, signed: -9},
			},
		},
		{
			data: Data{0x3f, 0xe4, 0x0e, 0xb6, 0x0c, 0xba, 0x00, 0x05},
			signals: []struct {
				start    uint8
				length   uint8
				unsigned uint64
				signed   int64
			}{
				{start: 7, length: 3, unsigned: 0x1, signed: 1},
				{start: 4, length: 1, unsigned: 0x1, signed: -1},
				{start: 55, length: 16, unsigned: 0x5, signed: 5},
				{start: 39, length: 16, unsigned: 0xcba, signed: 3258},
				{start: 23, length: 16, unsigned: 0xeb6, signed: 3766},
				{start: 3, length: 12, unsigned: 0xfe4, signed: -28},
			},
		},
		{
			data: Data{0x30, 0x53, 0x23, 0xe5, 0x0e, 0x11, 0xff, 0xff},
			signals: []struct {
				start    uint8
				length   uint8
				unsigned uint64
				signed   int64
			}{
				{start: 7, length: 3, unsigned: 0x1, signed: 1},
				{start: 4, length: 1, unsigned: 0x1, signed: -1},
				{start: 55, length: 16, unsigned: 0xffff, signed: -1},
				{start: 39, length: 16, unsigned: 0xe11, signed: 3601},
				{start: 23, length: 16, unsigned: 0x23e5, signed: 9189},
				{start: 3, length: 12, unsigned: 0x53, signed: 83},
			},
		},
	} {
		i, tt := i, tt
		t.Run(fmt.Sprintf("UnsignedBits:%v", i), func(t *testing.T) {
			for j, signal := range tt.signals {
				j, signal := j, signal
				t.Run(fmt.Sprintf("signal:%v", j), func(t *testing.T) {
					actual := tt.data.UnsignedBitsBigEndian(signal.start, signal.length)
					assert.Equal(t, signal.unsigned, actual)
				})
			}
		})
		t.Run(fmt.Sprintf("SignedBits:%v", i), func(t *testing.T) {
			for j, signal := range tt.signals {
				j, signal := j, signal
				t.Run(fmt.Sprintf("signal:%v", j), func(t *testing.T) {
					actual := tt.data.SignedBitsBigEndian(signal.start, signal.length)
					assert.Equal(t, signal.signed, actual)
				})
			}
		})
		t.Run(fmt.Sprintf("SetUnsignedBits:%v", i), func(t *testing.T) {
			var data Data
			for j, signal := range tt.signals {
				j, signal := j, signal
				t.Run(fmt.Sprintf("data:%v", j), func(_ *testing.T) {
					data.SetUnsignedBitsBigEndian(signal.start, signal.length, signal.unsigned)
				})
			}
			assert.DeepEqual(t, tt.data, data)
		})
		t.Run(fmt.Sprintf("SetSignedBits:%v", i), func(t *testing.T) {
			var data Data
			for j, signal := range tt.signals {
				j, signal := j, signal
				t.Run(fmt.Sprintf("data:%v", j), func(_ *testing.T) {
					data.SetSignedBitsBigEndian(signal.start, signal.length, signal.signed)
				})
			}
			assert.DeepEqual(t, tt.data, data)
		})
	}
}

func TestInvertEndian_Property_Idempotent(t *testing.T) {
	for i := uint8(0); i < 64; i++ {
		assert.Equal(t, i, invertEndian(invertEndian(i)))
	}
}

func TestPackUnpackBigEndian(t *testing.T) {
	f := func(data Data) Data {
		return data
	}
	g := func(data Data) Data {
		data.UnpackBigEndian(data.PackBigEndian())
		return data
	}
	assert.NilError(t, quick.CheckEqual(f, g, nil))
}

func TestPackUnpackLittleEndian(t *testing.T) {
	f := func(data Data) Data {
		return data
	}
	g := func(data Data) Data {
		data.UnpackLittleEndian(data.PackLittleEndian())
		return data
	}
	assert.NilError(t, quick.CheckEqual(f, g, nil))
}

func TestData_CheckBitRange(t *testing.T) {
	// example case that big-endian signals and little-endian signals use different indexing
	assert.NilError(t, CheckBitRangeBigEndian(8, 55, 16))
	assert.ErrorContains(t, CheckBitRangeLittleEndian(8, 55, 16), "bit range out of bounds")
}

func BenchmarkData_UnpackLittleEndian(b *testing.B) {
	var data Data
	for i := 0; i < b.N; i++ {
		data.UnpackLittleEndian(0)
	}
}

func BenchmarkData_UnpackBigEndian(b *testing.B) {
	var data Data
	for i := 0; i < b.N; i++ {
		data.UnpackBigEndian(0)
	}
}

func BenchmarkData_PackBigEndian(b *testing.B) {
	var data Data
	for i := 0; i < b.N; i++ {
		_ = data.PackBigEndian()
	}
}

func BenchmarkData_PackLittleEndian(b *testing.B) {
	var data Data
	for i := 0; i < b.N; i++ {
		_ = data.PackLittleEndian()
	}
}
