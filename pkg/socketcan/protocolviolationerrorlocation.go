package socketcan

type ProtocolViolationErrorLocation uint8

//go:generate stringer -type ProtocolViolationErrorLocation -trimprefix ProtocolViolationErrorLocation

const (
	ProtocolViolationErrorLocationUnspecified    ProtocolViolationErrorLocation = 0x00
	ProtocolViolationErrorLocationStartOfFrame   ProtocolViolationErrorLocation = 0x03
	ProtocolViolationErrorLocationID28To21       ProtocolViolationErrorLocation = 0x02 // standard frames: 10 - 3
	ProtocolViolationErrorLocationID20To18       ProtocolViolationErrorLocation = 0x06 // standard frames: 2 - 0
	ProtocolViolationErrorLocationSubstituteRTR  ProtocolViolationErrorLocation = 0x04 // standard frames: RTR
	ProtocolViolationErrorLocationIDExtension    ProtocolViolationErrorLocation = 0x05
	ProtocolViolationErrorLocationIDBits17To13   ProtocolViolationErrorLocation = 0x07
	ProtocolViolationErrorLocationIDBits12To05   ProtocolViolationErrorLocation = 0x0F
	ProtocolViolationErrorLocationIDBits04To00   ProtocolViolationErrorLocation = 0x0E
	ProtocolViolationErrorLocationRTR            ProtocolViolationErrorLocation = 0x0C
	ProtocolViolationErrorLocationReservedBit1   ProtocolViolationErrorLocation = 0x0D
	ProtocolViolationErrorLocationReservedBit0   ProtocolViolationErrorLocation = 0x09
	ProtocolViolationErrorLocationDataLengthCode ProtocolViolationErrorLocation = 0x0B
	ProtocolViolationErrorLocationData           ProtocolViolationErrorLocation = 0x0A
	ProtocolViolationErrorLocationCRCSequence    ProtocolViolationErrorLocation = 0x08
	ProtocolViolationErrorLocationCRCDelimiter   ProtocolViolationErrorLocation = 0x18
	ProtocolViolationErrorLocationACKSlot        ProtocolViolationErrorLocation = 0x19
	ProtocolViolationErrorLocationACKDelimiter   ProtocolViolationErrorLocation = 0x1B
	ProtocolViolationErrorLocationEndOfFrame     ProtocolViolationErrorLocation = 0x1A
	ProtocolViolationErrorLocationIntermission   ProtocolViolationErrorLocation = 0x12
)
