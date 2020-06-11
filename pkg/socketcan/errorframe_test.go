package socketcan

import (
	"testing"

	"gotest.tools/v3/assert"
)

func TestErrorFrame_String(t *testing.T) {
	for _, tt := range []struct {
		msg      string
		f        ErrorFrame
		expected string
	}{
		{
			msg: "lost arbitration",
			f: ErrorFrame{
				ErrorClass:         ErrorClassLostArbitration,
				LostArbitrationBit: 42,
			},
			expected: "LostArbitration in bit 42 (000000)",
		},
		{
			msg: "controller",
			f: ErrorFrame{
				ErrorClass:      ErrorClassController,
				ControllerError: ControllerErrorRxBufferOverflow,
			},
			expected: "Controller: RxBufferOverflow (000000)",
		},
		{
			msg: "protocol violation",
			f: ErrorFrame{
				ErrorClass:                     ErrorClassProtocolViolation,
				ProtocolError:                  ProtocolViolationErrorFrameFormat,
				ProtocolViolationErrorLocation: ProtocolViolationErrorLocationID20To18,
			},
			expected: "ProtocolViolation: FrameFormat: location ID20To18 (000000)",
		},
		{
			msg: "transceiver",
			f: ErrorFrame{
				ErrorClass:       ErrorClassTransceiver,
				TransceiverError: TransceiverErrorCANHShortToGND,
			},
			expected: "Transceiver: CANHShortToGND (000000)",
		},
		{
			msg: "controller specific information",
			f: ErrorFrame{
				ErrorClass:                    ErrorClassTxTimeout,
				ControllerSpecificInformation: [3]byte{0x12, 0x34, 0x56},
			},
			expected: "TxTimeout (123456)",
		},
	} {
		tt := tt
		t.Run(tt.msg, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.f.String())
		})
	}
}
