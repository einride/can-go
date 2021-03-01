package can

import (
	"encoding/hex"
	"fmt"
	"math/big"
)

// const MaxDataLength = 8

// // Data holds the data in a CAN frame.
// //
// // Layout
// //
// // Individual bits in the data are numbered according to the following scheme:
// //
// //             BIT
// //             NUMBER
// //             +------+------+------+------+------+------+------+------+
// //             |   7  |   6  |   5  |   4  |   3  |   2  |   1  |   0  |
// //    BYTE     +------+------+------+------+------+------+------+------+
// //    NUMBER
// //    +-----+  +------+------+------+------+------+------+------+------+
// //    |  0  |  |   7  |   6  |   5  |   4  |   3  |   2  |   1  |   0  |
// //    +-----+  +------+------+------+------+------+------+------+------+
// //    |  1  |  |  15  |  14  |  13  |  12  |  11  |  10  |   9  |   8  |
// //    +-----+  +------+------+------+------+------+------+------+------+
// //    |  2  |  |  23  |  22  |  21  |  20  |  19  |  18  |  17  |  16  |
// //    +-----+  +------+------+------+------+------+------+------+------+
// //    |  3  |  |  31  |  30  |  29  |  28  |  27  |  26  |  25  |  24  |
// //    +-----+  +------+------+------+------+------+------+------+------+
// //    |  4  |  |  39  |  38  |  37  |  36  |  35  |  34  |  33  |  32  |
// //    +-----+  +------+------+------+------+------+------+------+------+
// //    |  5  |  |  47  |  46  |  45  |  44  |  43  |  42  |  41  |  40  |
// //    +-----+  +------+------+------+------+------+------+------+------+
// //    |  6  |  |  55  |  54  |  53  |  52  |  51  |  50  |  49  |  48  |
// //    +-----+  +------+------+------+------+------+------+------+------+
// //    |  7  |  |  63  |  62  |  61  |  60  |  59  |  58  |  57  |  56  |
// //    +-----+  +------+------+------+------+------+------+------+------+
// //
// // Bit ranges can be manipulated using little-endian and big-endian bit ordering.
// //
// // Little-endian bit ranges
// //
// // Example range of length 32 starting at bit 29:
// //
// //             BIT
// //             NUMBER
// //             +------+------+------+------+------+------+------+------+
// //             |   7  |   6  |   5  |   4  |   3  |   2  |   1  |   0  |
// //    BYTE     +------+------+------+------+------+------+------+------+
// //    NUMBER
// //    +-----+  +------+------+------+------+------+------+------+------+
// //    |  0  |  |   7  |   6  |   5  |   4  |   3  |   2  |   1  |   0  |
// //    +-----+  +------+------+------+------+------+------+------+------+
// //    |  1  |  |  15  |  14  |  13  |  12  |  11  |  10  |   9  |   8  |
// //    +-----+  +------+------+------+------+------+------+------+------+
// //    |  2  |  |  23  |  22  |  21  |  20  |  19  |  18  |  17  |  16  |
// //    +-----+  +------+------+------+------+------+------+------+------+
// //    |  3  |  |  <-------------LSb |  28  |  27  |  26  |  25  |  24  |
// //    +-----+  +------+------+------+------+------+------+------+------+
// //    |  4  |  |  <--------------------------------------------------  |
// //    +-----+  +------+------+------+------+------+------+------+------+
// //    |  5  |  |  <--------------------------------------------------  |
// //    +-----+  +------+------+------+------+------+------+------+------+
// //    |  6  |  |  <--------------------------------------------------  |
// //    +-----+  +------+------+------+------+------+------+------+------+
// //    |  7  |  |  63  |  62  |  61  | <-MSb--------------------------- |
// //    +-----+  +------+------+------+------+------+------+------+------+
// //
// // Big-endian bit ranges
// //
// // Example range of length 32 starting at bit 29:
// //
// //             BIT
// //             NUMBER
// //             +------+------+------+------+------+------+------+------+
// //             |   7  |   6  |   5  |   4  |   3  |   2  |   1  |   0  |
// //    BYTE     +------+------+------+------+------+------+------+------+
// //    NUMBER
// //    +-----+  +------+------+------+------+------+------+------+------+
// //    |  0  |  |   7  |   6  |   5  |   4  |   3  |   2  |   1  |   0  |
// //    +-----+  +------+------+------+------+------+------+------+------+
// //    |  1  |  |  15  |  14  |  13  |  12  |  11  |  10  |   9  |   8  |
// //    +-----+  +------+------+------+------+------+------+------+------+
// //    |  2  |  |  23  |  22  |  21  |  20  |  19  |  18  |  17  |  16  |
// //    +-----+  +------+------+------+------+------+------+------+------+
// //    |  3  |  |  31  |  30  | <-MSb---------------------------------  |
// //    +-----+  +------+------+------+------+------+------+------+------+
// //    |  4  |  |  <--------------------------------------------------  |
// //    +-----+  +------+------+------+------+------+------+------+------+
// //    |  5  |  |  <--------------------------------------------------  |
// //    +-----+  +------+------+------+------+------+------+------+------+
// //    |  6  |  |  <--------------------------------------------------  |
// //    +-----+  +------+------+------+------+------+------+------+------+
// //    |  7  |  |  <------LSb |  61  |  60  |  59  |  58  |  57  |  56  |
// //    +-----+  +------+------+------+------+------+------+------+------+
// type Data [MaxDataLength]byte

//type Packed *big.Int

type Payload struct {
	// Message length
	Length uint16

	// Binary data
	Data []byte

	// Packed little endian
	PackedLittleEndian *big.Int

	// Packed big endian
	PackedBigEndian *big.Int
}

func (p *Payload) Hex() string {
	h := hex.EncodeToString(p.Data)
	return h
}

func PayloadFromHex(hexString string) (Payload, error) {
	b, err := hex.DecodeString(hexString)
	var p Payload
	if err != nil {
		return p, err
	}
	p = Payload{Data: b}
	return p, nil
}

// UnsignedBitsLittleEndian returns the little-endian bit range [start, start+length) as an unsigned value.
func (p *Payload) UnsignedBitsLittleEndian(start, length uint16) uint64 {
	// pack bits into one continuous value

	packed := p.PackLittleEndian()
	// lsb index in the packed value is the start bit
	lsbIndex := uint(start)
	// shift away lower bits
	shifted := packed.Rsh(packed, lsbIndex)
	// mask away higher bits
	//masked := shifted & ((1 << length) - 1)
	masked := shifted.And(shifted, big.NewInt((1<<length)-1))
	// done
	return masked.Uint64()
}

// UnsignedBitsBigEndian returns the big-endian bit range [start, start+length) as an unsigned value.
func (p *Payload) UnsignedBitsBigEndian(start, length uint16) uint64 {
	// pack bits into one continuous value
	packed := p.PackBigEndian()
	// calculate msb index in the packed value
	msbIndex := p.invertEndian(start)
	// calculate lsb index in the packed value
	lsbIndex := uint(msbIndex - length + 1)
	// shift away lower bits
	shifted := packed.Rsh(packed, lsbIndex)
	// mask away higher bits
	masked := shifted.And(shifted, big.NewInt((1<<length)-1))
	// done
	return masked.Uint64()
}

// SignedBitsLittleEndian returns little-endian bit range [start, start+length) as a signed value.
func (p *Payload) SignedBitsLittleEndian(start, length uint16) int64 {
	unsigned := p.UnsignedBitsLittleEndian(start, length)
	return AsSigned(unsigned, length)
}

// SignedBitsBigEndian returns little-endian bit range [start, start+length) as a signed value.
func (p *Payload) SignedBitsBigEndian(start, length uint16) int64 {
	unsigned := p.UnsignedBitsBigEndian(start, length)
	return AsSigned(unsigned, length)
}

// SignedBitsBigEndian returns little-endian bit range [start, start+length) as a signed value.
// func (d *Data) SignedBitsBigEndian(start, length uint8) int64 {
// 	unsigned := d.UnsignedBitsBigEndian(start, length)
// 	return reinterpret.AsSigned(unsigned, length)
// }

// // SetUnsignedBitsBigEndian sets the little-endian bit range [start, start+length) to the provided unsigned value.
// func (d *Data) SetUnsignedBitsLittleEndian(start, length uint8, value uint64) {
// 	// pack bits into one continuous value
// 	packed := d.PackLittleEndian()
// 	// lsb index in the packed value is the start bit
// 	lsbIndex := start
// 	// calculate bit mask for zeroing the bit range to set
// 	unsetMask := ^uint64(((1 << length) - 1) << lsbIndex)
// 	// calculate bit mask for setting the new value
// 	setMask := value << lsbIndex
// 	// calculate the new packed value
// 	newPacked := packed&unsetMask | setMask
// 	// unpack the new packed value into the data
// 	d.UnpackLittleEndian(newPacked)
// }

// // SetUnsignedBitsBigEndian sets the big-endian bit range [start, start+length) to the provided unsigned value.
// func (d *Data) SetUnsignedBitsBigEndian(start, length uint8, value uint64) {
// 	// pack bits into one continuous value
// 	packed := d.PackBigEndian()
// 	// calculate msb index in the packed value
// 	msbIndex := invertEndian(start)
// 	// calculate lsb index in the packed value
// 	lsbIndex := msbIndex - length + 1
// 	// calculate bit mask for zeroing the bit range to set
// 	unsetMask := ^uint64(((1 << length) - 1) << lsbIndex)
// 	// calculate bit mask for setting the new value
// 	setMask := value << lsbIndex
// 	// calculate the new packed value
// 	newPacked := packed&unsetMask | setMask
// 	// unpack the new packed value into the data
// 	d.UnpackBigEndian(newPacked)
// }

// // SetSignedBitsLittleEndian sets the little-endian bit range [start, start+length) to the provided signed value.
// func (d *Data) SetSignedBitsLittleEndian(start, length uint8, value int64) {
// 	d.SetUnsignedBitsLittleEndian(start, length, reinterpret.AsUnsigned(value, length))
// }

// // SetSignedBitsBigEndian sets the big-endian bit range [start, start+length) to the provided signed value.
// func (d *Data) SetSignedBitsBigEndian(start, length uint8, value int64) {
// 	d.SetUnsignedBitsBigEndian(start, length, reinterpret.AsUnsigned(value, length))
// }

// Bit returns the value of the i:th bit in the data as a bool.
func (p *Payload) Bit(i uint16) bool {
	if i > 8*p.Length-1 {
		return false
	}
	// calculate which byte the bit belongs to
	byteIndex := i / 8
	// calculate bit mask for extracting the bit
	bitMask := uint8(1 << (i % 8))
	// mocks the bit
	bit := p.Data[byteIndex]&bitMask > 0
	// done
	return bit
}

// SetBit sets the value of the i:th bit in the data.
func (p *Payload) SetBit(i uint16, value bool) {
	if i > 8*p.Length-1 {
		return
	}
	byteIndex := i / 8
	bitIndex := i % 8
	if value {
		p.Data[byteIndex] |= uint8(1 << bitIndex)
	} else {
		p.Data[byteIndex] &= ^uint8(1 << bitIndex)
	}
}

// PackLittleEndian packs the data into a contiguous uint64 value for little-endian signals.
// func (d *Data) PackLittleEndian() uint64 {
// 	var packed uint64
// 	packed |= uint64(d[0]) << (0 * 8)
// 	packed |= uint64(d[1]) << (1 * 8)
// 	packed |= uint64(d[2]) << (2 * 8)
// 	packed |= uint64(d[3]) << (3 * 8)
// 	packed |= uint64(d[4]) << (4 * 8)
// 	packed |= uint64(d[5]) << (5 * 8)
// 	packed |= uint64(d[6]) << (6 * 8)
// 	packed |= uint64(d[7]) << (7 * 8)
// 	return packed
// }

func (p *Payload) PackLittleEndian() *big.Int {
	if p.PackedLittleEndian == nil {
		// Initialize a big.Int called packed
		//var packed, _ = new(big.Int).SetString(strings.Repeat("00000000", int(p.Length)), 2)
		packed, _ := new(big.Int).SetString(CanBitsLittleEndian(p.Data), 2)

		// for i := 0; i < int(p.Length); i++ {
		// 	//packed |= uint8(packed >> (i * 8))
		// 	fmt.Println(fmt.Sprintf("%08b", p.Data[i]))
		// }
		//packed.SetBytes(p.Data)
		p.PackedLittleEndian = packed
	}
	return new(big.Int).Set(p.PackedLittleEndian)
}

// CanBitsLittleEndian creates the zig zag pattern of bits feeding into a can.Data frame
func CanBitsLittleEndian(bytes []byte) string {
	var bits string
	for _, n := range bytes {
		bn := fmt.Sprintf("%08b", n)

		// Need to reverse the binary strings because of the definition of bit order
		bits += reverse(bn)
	}
	return reverse(bits)
}

// CanBitsBigEndian creates the zig zag pattern of bits feeding into a can.Data frame
func CanBitsBigEndian(bytes []byte) string {
	var bits string
	for _, n := range bytes {
		bn := fmt.Sprintf("%08b", n)
		bits += bn
	}
	return bits
}

// function to reverse strings
func reverse(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

// // PackBigEndian packs the data into a contiguous uint64 value for big-endian signals.
// func (d *Data) PackBigEndian() uint64 {
// 	var packed uint64
// 	packed |= uint64(d[0]) << (7 * 8)
// 	packed |= uint64(d[1]) << (6 * 8)
// 	packed |= uint64(d[2]) << (5 * 8)
// 	packed |= uint64(d[3]) << (4 * 8)
// 	packed |= uint64(d[4]) << (3 * 8)
// 	packed |= uint64(d[5]) << (2 * 8)
// 	packed |= uint64(d[6]) << (1 * 8)
// 	packed |= uint64(d[7]) << (0 * 8)
// 	return packed
// }

func (p *Payload) PackBigEndian() *big.Int {
	if p.PackedBigEndian == nil {
		var packed, _ = new(big.Int).SetString(CanBitsBigEndian(p.Data), 2)
		p.PackedBigEndian = packed
	}
	return new(big.Int).Set(p.PackedBigEndian)
}

// // UnpackLittleEndian sets the value of d.Bytes by unpacking the provided value as sequential little-endian bits.
// func (d *Data) UnpackLittleEndian(packed uint64) {
// 	d[0] = uint8(packed >> (0 * 8))
// 	d[1] = uint8(packed >> (1 * 8))
// 	d[2] = uint8(packed >> (2 * 8))
// 	d[3] = uint8(packed >> (3 * 8))
// 	d[4] = uint8(packed >> (4 * 8))
// 	d[5] = uint8(packed >> (5 * 8))
// 	d[6] = uint8(packed >> (6 * 8))
// 	d[7] = uint8(packed >> (7 * 8))
// }

// // UnpackBigEndian sets the value of d.Bytes by unpacking the provided value as sequential big-endian bits.
// func (d *Data) UnpackBigEndian(packed uint64) {
// 	d[0] = uint8(packed >> (7 * 8))
// 	d[1] = uint8(packed >> (6 * 8))
// 	d[2] = uint8(packed >> (5 * 8))
// 	d[3] = uint8(packed >> (4 * 8))
// 	d[4] = uint8(packed >> (3 * 8))
// 	d[5] = uint8(packed >> (2 * 8))
// 	d[6] = uint8(packed >> (1 * 8))
// 	d[7] = uint8(packed >> (0 * 8))
// }

// invertEndian converts from big-endian to little-endian bit indexing and vice versa.
func (p *Payload) invertEndian(i uint16) uint16 {
	row := i / 8
	col := i % 8
	oppositeRow := uint16(len(p.Data)) - row - 1
	bitIndex := (oppositeRow * 8) + col
	return bitIndex
}

// // CheckBitRangeLittleEndian checks that a little-endian bit range fits in the data.
// func CheckBitRangeLittleEndian(frameLength, rangeStart, rangeLength uint8) error {
// 	lsbIndex := rangeStart
// 	msbIndex := rangeStart + rangeLength - 1
// 	upperBound := frameLength * 8
// 	if msbIndex >= upperBound {
// 		return fmt.Errorf("bit range out of bounds [0, %v): [%v, %v]", upperBound, lsbIndex, msbIndex)
// 	}
// 	return nil
// }

// // CheckBitRangeBigEndian checks that a big-endian bit range fits in the data.
// func CheckBitRangeBigEndian(frameLength, rangeStart, rangeLength uint8) error {
// 	upperBound := frameLength * 8
// 	if rangeStart >= upperBound {
// 		return fmt.Errorf("bit range starts out of bounds [0, %v): %v", upperBound, rangeStart)
// 	}
// 	msbIndex := invertEndian(rangeStart)
// 	lsbIndex := msbIndex - rangeLength + 1
// 	end := invertEndian(lsbIndex)
// 	if end >= upperBound {
// 		return fmt.Errorf("bit range ends out of bounds [0, %v): %v", upperBound, end)
// 	}
// 	return nil
// }

// // CheckValue checks that a value fits in a number of bits.
// func CheckValue(value uint64, bits uint8) error {
// 	upperBound := uint64(1 << bits)
// 	if value >= upperBound {
// 		return fmt.Errorf("value out of bounds [0, %v): %v", upperBound, value)
// 	}
// 	return nil
// }

// AsSigned reinterprets the provided unsigned value as a signed value.
func AsSigned(unsigned uint64, bits uint16) int64 {
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
