package socketcan

type ErrorClass uint32

//go:generate stringer -type ErrorClass -trimprefix ErrorClass

const (
	ErrorClassTxTimeout         ErrorClass = 0x00000001
	ErrorClassLostArbitration   ErrorClass = 0x00000002
	ErrorClassController        ErrorClass = 0x00000004
	ErrorClassProtocolViolation ErrorClass = 0x00000008
	ErrorClassTransceiver       ErrorClass = 0x00000010
	ErrorClassNoAck             ErrorClass = 0x00000020
	ErrorClassBusOff            ErrorClass = 0x00000040
	ErrorClassBusError          ErrorClass = 0x00000080
	ErrorClassRestarted         ErrorClass = 0x00000100
)
