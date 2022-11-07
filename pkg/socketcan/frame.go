package socketcan

import (
	"encoding/binary"

	"go.einride.tech/can"
)

const (
	// lengthOfFrame is the length of a SocketCAN frame in bytes.
	lengthOfFrame = 16
	// maxLengthOfData is the max length of a SocketCAN frame payload in bytes.
	maxLengthOfData = 8
	// indexOfID is the index of the first byte of the frame ID.
	indexOfID = 0
	// lengthOfID is the length of a frame ID in bytes.
	lengthOfID = 4
	// indexOfDataLengthCode is the index of the first byte of the frame dataLengthCode.
	indexOfDataLengthCode = indexOfID + lengthOfID
	// lengthOfDataLengthCode is the length of a frame dataLengthCode in bytes.
	lengthOfDataLengthCode = 1
	// indexOfPadding is the index of the first byte of frame padding.
	indexOfPadding = indexOfDataLengthCode + lengthOfDataLengthCode
	// lengthOfPadding is the length of frame padding in bytes.
	lengthOfPadding = 3
	// indexOfData is the index of the first byte of data in a frame.
	indexOfData = indexOfPadding + lengthOfPadding
)

// error frame flag indices.
const (
	// indexOfLostArbitrationBit is the byte index of the lost arbitration bit in an error frame.
	indexOfLostArbitrationBit = 0
	// indexOfControllerError is the byte index of the controller error in an error frame.
	indexOfControllerError = 1
	// indexOfProtocolError is the byte index of the protocol error in an error frame.
	indexOfProtocolError = 2
	// indexOfProtocolViolationErrorLocation is the byte index of the protocol error location in an error frame.
	indexOfProtocolViolationErrorLocation = 3
	// indexOfTransceiverError is the byte index of the transceiver error in an error frame.
	indexOfTransceiverError = 4
	// indexOfControllerSpecificInformation is the starting byte of controller specific information in an error frame.
	indexOfControllerSpecificInformation = 5
	// LengthOfControllerSpecificInformation is the number of error frame bytes with controller-specific information.
	LengthOfControllerSpecificInformation = 3
)

var _ [lengthOfFrame]struct{} = [indexOfData + maxLengthOfData]struct{}{}

// id flags (copied from x/sys/unix).
const (
	idFlagExtended = 0x80000000
	idFlagError    = 0x20000000
	idFlagRemote   = 0x40000000
	idMaskExtended = 0x1fffffff
	idMaskStandard = 0x7ff
)

// FrameInterceptor provides a hook to intercept the transmission of a CAN frame.
// The interceptor is called if and only if the frame transmission/receival is a success.
type FrameInterceptor func(fr can.Frame)

// frame represents a SocketCAN frame.
//
// The format specified in the Linux SocketCAN kernel module:
//
//	struct can_frame {
//	        canid_t can_id;  /* 32 bit CAN_ID + EFF/RTR/ERR flags */
//	        __u8    can_dlc; /* frame payload length in byte (0 .. 8) */
//	        __u8    __pad;   /* padding */
//	        __u8    __res0;  /* reserved / padding */
//	        __u8    __res1;  /* reserved / padding */
//	        __u8    data[8] __attribute__((aligned(8)));
//	};
type frame struct {
	// idAndFlags is the combined CAN ID and flags.
	idAndFlags uint32
	// dataLengthCode is the frame payload length in bytes.
	dataLengthCode uint8
	// padding+reserved fields
	_ [3]byte
	// bytes contains the frame payload.
	data [8]byte
}

func (f *frame) unmarshalBinary(b []byte) {
	_ = b[lengthOfFrame-1] // bounds check
	f.idAndFlags = binary.LittleEndian.Uint32(b[indexOfID : indexOfID+lengthOfID])
	f.dataLengthCode = b[indexOfDataLengthCode]
	copy(f.data[:], b[indexOfData:lengthOfFrame])
}

func (f *frame) marshalBinary(b []byte) {
	_ = b[lengthOfFrame-1] // bounds check
	binary.LittleEndian.PutUint32(b[indexOfID:indexOfID+lengthOfID], f.idAndFlags)
	b[indexOfDataLengthCode] = f.dataLengthCode
	copy(b[indexOfData:], f.data[:])
}

func (f *frame) decodeFrame() can.Frame {
	return can.Frame{
		ID:         f.id(),
		Length:     f.dataLengthCode,
		Data:       f.data,
		IsExtended: f.isExtended(),
		IsRemote:   f.isRemote(),
	}
}

func (f *frame) encodeFrame(cf can.Frame) {
	f.idAndFlags = cf.ID
	if cf.IsRemote {
		f.idAndFlags |= idFlagRemote
	}
	if cf.IsExtended {
		f.idAndFlags |= idFlagExtended
	}
	f.dataLengthCode = cf.Length
	f.data = cf.Data
}

func (f *frame) isExtended() bool {
	return f.idAndFlags&idFlagExtended > 0
}

func (f *frame) isRemote() bool {
	return f.idAndFlags&idFlagRemote > 0
}

func (f *frame) isError() bool {
	return f.idAndFlags&idFlagError > 0
}

func (f *frame) id() uint32 {
	if f.isExtended() {
		return f.idAndFlags & idMaskExtended
	}
	return f.idAndFlags & idMaskStandard
}

func (f *frame) decodeErrorFrame() ErrorFrame {
	return ErrorFrame{
		ErrorClass:                     f.errorClass(),
		LostArbitrationBit:             f.lostArbitrationBit(),
		ControllerError:                f.controllerError(),
		ProtocolError:                  f.protocolError(),
		ProtocolViolationErrorLocation: f.protocolErrorLocation(),
		TransceiverError:               f.transceiverError(),
		ControllerSpecificInformation:  f.controllerSpecificInformation(),
	}
}

func (f *frame) errorClass() ErrorClass {
	return ErrorClass(f.idAndFlags &^ idFlagError)
}

func (f *frame) lostArbitrationBit() uint8 {
	return f.data[indexOfLostArbitrationBit]
}

func (f *frame) controllerError() ControllerError {
	return ControllerError(f.data[indexOfControllerError])
}

func (f *frame) protocolError() ProtocolViolationError {
	return ProtocolViolationError(f.data[indexOfProtocolError])
}

func (f *frame) protocolErrorLocation() ProtocolViolationErrorLocation {
	return ProtocolViolationErrorLocation(f.data[indexOfProtocolViolationErrorLocation])
}

func (f *frame) transceiverError() TransceiverError {
	return TransceiverError(f.data[indexOfTransceiverError])
}

func (f *frame) controllerSpecificInformation() [LengthOfControllerSpecificInformation]byte {
	var ret [LengthOfControllerSpecificInformation]byte
	start := indexOfControllerSpecificInformation
	end := start + LengthOfControllerSpecificInformation
	copy(ret[:], f.data[start:end])
	return ret
}
