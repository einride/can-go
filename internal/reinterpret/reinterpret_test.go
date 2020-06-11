package reinterpret

import (
	"fmt"
	"testing"

	"gotest.tools/v3/assert"
)

func TestReinterpretSign(t *testing.T) {
	for _, tt := range []struct {
		unsigned uint64
		length   uint8
		signed   int64
	}{
		// -1, byte aligned
		{unsigned: 0xf, length: 4, signed: -1},
		{unsigned: 0xff, length: 8, signed: -1},
		{unsigned: 0xfff, length: 12, signed: -1},
		{unsigned: 0xffff, length: 16, signed: -1},
		{unsigned: 0xfffff, length: 20, signed: -1},
		{unsigned: 0xffffff, length: 24, signed: -1},
		{unsigned: 0xfffffff, length: 28, signed: -1},
		{unsigned: 0xffffffff, length: 32, signed: -1},
		{unsigned: 0xfffffffff, length: 36, signed: -1},
		{unsigned: 0xffffffffff, length: 40, signed: -1},
		{unsigned: 0xfffffffffff, length: 44, signed: -1},
		{unsigned: 0xffffffffffff, length: 48, signed: -1},
		{unsigned: 0xfffffffffffff, length: 52, signed: -1},
		{unsigned: 0xffffffffffffff, length: 56, signed: -1},
		{unsigned: 0xfffffffffffffff, length: 60, signed: -1},
		{unsigned: 0xffffffffffffffff, length: 64, signed: -1},
		// 3 bits
		{unsigned: 0x0, length: 3, signed: 0},
		{unsigned: 0x1, length: 3, signed: 1},
		{unsigned: 0x2, length: 3, signed: 2},
		{unsigned: 0x3, length: 3, signed: 3},
		{unsigned: 0x4, length: 3, signed: -4},
		{unsigned: 0x5, length: 3, signed: -3},
		{unsigned: 0x6, length: 3, signed: -2},
		{unsigned: 0x7, length: 3, signed: -1},
		// 4 bits
		{unsigned: 0x0, length: 4, signed: 0},
		{unsigned: 0x1, length: 4, signed: 1},
		{unsigned: 0x2, length: 4, signed: 2},
		{unsigned: 0x3, length: 4, signed: 3},
		{unsigned: 0x4, length: 4, signed: 4},
		{unsigned: 0x5, length: 4, signed: 5},
		{unsigned: 0x6, length: 4, signed: 6},
		{unsigned: 0x7, length: 4, signed: 7},
		{unsigned: 0x8, length: 4, signed: -8},
		{unsigned: 0x9, length: 4, signed: -7},
		{unsigned: 0xa, length: 4, signed: -6},
		{unsigned: 0xb, length: 4, signed: -5},
		{unsigned: 0xc, length: 4, signed: -4},
		{unsigned: 0xd, length: 4, signed: -3},
		{unsigned: 0xe, length: 4, signed: -2},
		{unsigned: 0xf, length: 4, signed: -1},
	} {
		tt := tt
		t.Run(fmt.Sprintf("%+v", tt), func(t *testing.T) {
			assert.Equal(t, tt.signed, AsSigned(tt.unsigned, tt.length))
			assert.Equal(t, tt.unsigned, AsUnsigned(tt.signed, tt.length))
			assert.Equal(t, tt.signed, AsSigned(AsUnsigned(tt.signed, tt.length), tt.length))
			assert.Equal(t, tt.unsigned, AsUnsigned(AsSigned(tt.unsigned, tt.length), tt.length))
		})
	}
}
