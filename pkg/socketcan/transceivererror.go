package socketcan

type TransceiverError uint8

//go:generate stringer -type TransceiverError -trimprefix TransceiverError

const (
	TransceiverErrorUnspecified     TransceiverError = 0x00
	TransceiverErrorCANHNoWire      TransceiverError = 0x04
	TransceiverErrorCANHShortToBat  TransceiverError = 0x05
	TransceiverErrorCANHShortToVCC  TransceiverError = 0x06
	TransceiverErrorCANHShortToGND  TransceiverError = 0x07
	TransceiverErrorCANLNoWire      TransceiverError = 0x40
	TransceiverErrorCANLShortToBat  TransceiverError = 0x50
	TransceiverErrorCANLShortToVcc  TransceiverError = 0x60
	TransceiverErrorCANLShortToGND  TransceiverError = 0x70
	TransceiverErrorCANLShortToCANH TransceiverError = 0x80
)
