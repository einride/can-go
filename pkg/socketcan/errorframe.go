package socketcan

import (
	"encoding/hex"
	"fmt"
)

type ErrorFrame struct {
	// Class is the error class
	ErrorClass ErrorClass
	// LostArbitrationBit contains the bit number when the error class is LostArbitration.
	LostArbitrationBit uint8
	// ControllerError contains error information when the error class is Controller.
	ControllerError ControllerError
	// ProtocolViolationError contains error information when the error class is Protocol.
	ProtocolError ProtocolViolationError
	// ProtocolViolationErrorLocation contains error location when the error class is Protocol.
	ProtocolViolationErrorLocation ProtocolViolationErrorLocation
	// TransceiverError contains error information when the error class is Transceiver.
	TransceiverError TransceiverError
	// ControllerSpecificInformation contains controller-specific additional error information.
	ControllerSpecificInformation [3]byte
}

func (e *ErrorFrame) String() string {
	switch e.ErrorClass {
	case ErrorClassLostArbitration:
		return fmt.Sprintf(
			"%s in bit %d (%s)",
			e.ErrorClass,
			e.LostArbitrationBit,
			hex.EncodeToString(e.ControllerSpecificInformation[:]),
		)
	case ErrorClassController:
		return fmt.Sprintf(
			"%s: %s (%v)",
			e.ErrorClass,
			e.ControllerError,
			hex.EncodeToString(e.ControllerSpecificInformation[:]),
		)
	case ErrorClassProtocolViolation:
		return fmt.Sprintf(
			"%s: %s: location %s (%v)",
			e.ErrorClass,
			e.ProtocolError,
			e.ProtocolViolationErrorLocation,
			hex.EncodeToString(e.ControllerSpecificInformation[:]),
		)
	case ErrorClassTransceiver:
		return fmt.Sprintf(
			"%s: %s (%v)",
			e.ErrorClass,
			e.TransceiverError,
			hex.EncodeToString(e.ControllerSpecificInformation[:]),
		)
	default:
		return fmt.Sprintf(
			"%s (%v)",
			e.ErrorClass,
			hex.EncodeToString(e.ControllerSpecificInformation[:]),
		)
	}
}
