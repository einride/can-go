package can

import (
	"fmt"
	"reflect"
	"testing"
)

type signals struct {
	start    uint16
	length   uint16
	unsigned uint64
	signed   int64
}

func TestPackLittleEndian(t *testing.T) {
	// 302064448
	// 10010000000010010001101000000
	data := []byte{0x40, 0x23, 0x01, 0x12}
	payload := Payload{Data: data}

	dataLittleEndian := payload.PackLittleEndian()
	fmt.Println(dataLittleEndian)
	fmt.Printf("%b\n", dataLittleEndian)
}

func TestPackBigEndian(t *testing.T) {
	// 4621538819433299968
	// 100000000100011000000010001001000000000000000000000000000000000
	data := []byte{0x40, 0x23, 0x01, 0x12}
	payload := Payload{Data: data}

	dataBigEndian := payload.PackBigEndian()
	fmt.Println(dataBigEndian)
	fmt.Printf("%b\n", dataBigEndian)
}

func TestUnsignedLittleEndian(t *testing.T) {
	// 18
	data := []byte{0x40, 0x23, 0x01, 0x12}
	payload := Payload{Data: data}
	signal := signals{start: 24, length: 8, unsigned: 0x12, signed: 18}
	fmt.Println(payload.UnsignedBitsLittleEndian(signal.start, signal.length))
}

func TestUnsignedBigEndian(t *testing.T) {
	// 3219
	data := []byte{0x3f, 0xf7, 0x0d, 0xc4, 0x0c, 0x93, 0xff, 0xff}
	payload := Payload{Data: data}
	signal := signals{start: 39, length: 16, unsigned: 0xc93, signed: 3219}
	fmt.Println(payload.UnsignedBitsBigEndian(signal.start, signal.length))
}

func TestSignedLittleEndian(t *testing.T) {
	// -1
	data := []byte{0x80, 0x01}
	payload := Payload{Data: data}
	signal := signals{start: 7, length: 2, unsigned: 0x3, signed: -1}
	fmt.Println(payload.SignedBitsLittleEndian(signal.start, signal.length))
}

func TestSignedBigEndian(t *testing.T) {
	// -9
	data := []byte{0x3f, 0xf7, 0x0d, 0xc4, 0x0c, 0x93, 0xff, 0xff}
	payload := Payload{Data: data}
	signal := signals{start: 3, length: 12, unsigned: 0xff7, signed: -9}
	fmt.Println(payload.SignedBitsBigEndian(signal.start, signal.length))
}

func TestReverse(t *testing.T) {
	payloadHexes := []string{
		"17399b6c89d22805003f", // Even Length
		"17399b6c89d22805003f17399b6c89d22805003f17399b6c89d22805003f", // Long Even
		"17399b6c89d2280500", // Odd Length
		"17399b6c89d228050017399b6c89d228050017399b6c89d2280500", // Long Odd
	}

	for _, hex := range payloadHexes {
		p, _ := PayloadFromHex(hex)

		twiceReversed := reverse(reverse(p.Data))
		if !reflect.DeepEqual(p.Data, twiceReversed) {
			t.Errorf("Reversed slices are not equal, input: %v, output: %v", p.Data, twiceReversed)
		}
	}
}

func Benchmark4BytesPayload_PackLittleEndian(b *testing.B) {
	data := []byte{0x40, 0x23, 0x01, 0x12}
	payload := Payload{Data: data}
	for i := 0; i < b.N; i++ {
		_ = payload.PackLittleEndian()
	}
}

func Benchmark4BytesPayload_PackBigEndian(b *testing.B) {
	data := []byte{0x40, 0x23, 0x01, 0x12}
	payload := Payload{Data: data}
	for i := 0; i < b.N; i++ {
		_ = payload.PackBigEndian()
	}
}

func Benchmark4BytesPayload_UnsignedBitsLittleEndian(b *testing.B) {
	data := []byte{0x40, 0x23, 0x01, 0x12}
	payload := Payload{Data: data}
	for i := 0; i < b.N; i++ {
		_ = payload.UnsignedBitsLittleEndian(0, 16)
	}
}

func Benchmark4BytesPayload_UnsignedBitsBigEndian(b *testing.B) {
	data := []byte{0x40, 0x23, 0x01, 0x12}
	payload := Payload{Data: data}
	for i := 0; i < b.N; i++ {
		_ = payload.UnsignedBitsBigEndian(0, 16)
	}
}
