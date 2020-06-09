package socketcan

type ProtocolViolationError uint8

//go:generate stringer -type ProtocolViolationError -trimprefix ProtocolViolationError

const (
	ProtocolViolationErrorUnspecified ProtocolViolationError = 0x00
	ProtocolViolationErrorSingleBit   ProtocolViolationError = 0x01
	ProtocolViolationErrorFrameFormat ProtocolViolationError = 0x02
	ProtocolViolationErrorBitStuffing ProtocolViolationError = 0x04
	ProtocolViolationErrorBit0        ProtocolViolationError = 0x08 // unable to send dominant bit
	ProtocolViolationErrorBit1        ProtocolViolationError = 0x10 // unable to send recessive bit
	ProtocolViolationErrorBusOverload ProtocolViolationError = 0x20
	ProtocolViolationErrorActive      ProtocolViolationError = 0x40 // active error announcement
	ProtocolViolationErrorTx          ProtocolViolationError = 0x80 // error occurred on transmission
)
