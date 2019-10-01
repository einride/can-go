// Package reinterpret provides primitives for reinterpreting arbitrary-length values as signed or unsigned.
package reinterpret

// AsSigned reinterprets the provided unsigned value as a signed value.
func AsSigned(unsigned uint64, bits uint8) int64 {
	switch bits {
	case 8:
		return int64(int8(uint8(unsigned)))
	case 16:
		return int64(int16(uint16(unsigned)))
	case 32:
		return int64(int32(uint32(unsigned)))
	case 64:
		return int64(unsigned)
	default:
		// calculate bit mask for sign bit
		signBitMask := uint64(1 << (bits - 1))
		// check if sign bit is set
		isNegative := unsigned&signBitMask > 0
		if !isNegative {
			// sign bit not set means we can reinterpret the value as-is
			return int64(unsigned)
		}
		// calculate bit mask for extracting value bits (all bits except the sign bit)
		valueBitMask := signBitMask - 1
		// calculate two's complement of the value bits
		value := ((^unsigned) & valueBitMask) + 1
		// result is the negative value of the two's complement
		return -1 * int64(value)
	}
}

// AsUnsigned reinterprets the provided signed value as an unsigned value.
func AsUnsigned(signed int64, bits uint8) uint64 {
	switch bits {
	case 8:
		return uint64(uint8(int8(signed)))
	case 16:
		return uint64(uint16(int16(signed)))
	case 32:
		return uint64(uint32(int32(signed)))
	case 64:
		return uint64(signed)
	default:
		// calculate bit mask for extracting relevant bits
		valueBitMask := uint64(1<<bits) - 1
		// extract relevant bits
		return uint64(signed) & valueBitMask
	}
}
