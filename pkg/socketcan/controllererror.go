package socketcan

type ControllerError uint8

//go:generate stringer -type ControllerError -trimprefix ControllerError

const (
	ControllerErrorUnspecified      ControllerError = 0x00
	ControllerErrorRxBufferOverflow ControllerError = 0x01
	ControllerErrorTxBufferOverflow ControllerError = 0x02
	ControllerErrorRxWarning        ControllerError = 0x04
	ControllerErrorTxWarning        ControllerError = 0x08
	ControllerErrorRxPassive        ControllerError = 0x10
	ControllerErrorTxPassive        ControllerError = 0x20 // at least one error counter exceeds 127
	ControllerErrorActive           ControllerError = 0x40
)
