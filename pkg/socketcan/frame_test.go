package socketcan

import (
	"testing"
	"testing/quick"

	"github.com/blueinnovationsgroup/can-go"
	"gotest.tools/v3/assert"
)

func TestFrame_MarshalUnmarshalBinary_Property_Idempotent(t *testing.T) {
	f := func(data [lengthOfFrame]byte) [lengthOfFrame]byte {
		data[5], data[6], data[7] = 0, 0, 0 // padding+reserved fields
		return data
	}
	g := func(data [lengthOfFrame]byte) [lengthOfFrame]byte {
		var f frame
		f.unmarshalBinary(data[:])
		var newData [lengthOfFrame]byte
		f.marshalBinary(newData[:])
		return newData
	}
	assert.NilError(t, quick.CheckEqual(f, g, nil))
}

func TestFrame_EncodeDecode(t *testing.T) {
	for _, tt := range []struct {
		msg            string
		frame          can.Frame
		socketCANFrame frame
	}{
		{
			msg: "data",
			frame: can.Frame{
				ID:     0x00000001,
				Length: 8,
				Data:   can.Data{1, 2, 3, 4, 5, 6, 7, 8},
			},
			socketCANFrame: frame{
				idAndFlags:     0x00000001,
				dataLengthCode: 8,
				data:           [8]byte{1, 2, 3, 4, 5, 6, 7, 8},
			},
		},
		{
			msg: "extended",
			frame: can.Frame{
				ID:         0x00000001,
				IsExtended: true,
			},
			socketCANFrame: frame{
				idAndFlags: 0x80000001,
			},
		},
		{
			msg: "remote",
			frame: can.Frame{
				ID:       0x00000001,
				IsRemote: true,
			},
			socketCANFrame: frame{
				idAndFlags: 0x40000001,
			},
		},
		{
			msg: "extended and remote",
			frame: can.Frame{
				ID:         0x00000001,
				IsExtended: true,
				IsRemote:   true,
			},
			socketCANFrame: frame{
				idAndFlags: 0xc0000001,
			},
		},
	} {
		tt := tt
		t.Run(tt.msg, func(t *testing.T) {
			t.Run("encode", func(t *testing.T) {
				var actual frame
				actual.encodeFrame(tt.frame)
				assert.Equal(t, tt.socketCANFrame, actual)
			})
			t.Run("decode", func(t *testing.T) {
				assert.Equal(t, tt.frame, tt.socketCANFrame.decodeFrame())
			})
		})
	}
}

func TestFrame_IsError(t *testing.T) {
	assert.Assert(t, (&frame{idAndFlags: 0x20000001}).isError())
	assert.Assert(t, !(&frame{idAndFlags: 0x00000001}).isError())
}

func TestFrame_DecodeErrorFrame(t *testing.T) {
	for _, tt := range []struct {
		msg      string
		f        frame
		expected ErrorFrame
	}{
		{
			msg: "lost arbitration",
			f: frame{
				idAndFlags:     0x20000002,
				dataLengthCode: 8,
				data: [8]byte{
					42,
				},
			},
			expected: ErrorFrame{
				ErrorClass:         ErrorClassLostArbitration,
				LostArbitrationBit: 42,
			},
		},
		{
			msg: "controller",
			f: frame{
				idAndFlags:     0x20000004,
				dataLengthCode: 8,
				data: [8]byte{
					0,
					0x04,
				},
			},
			expected: ErrorFrame{
				ErrorClass:      ErrorClassController,
				ControllerError: ControllerErrorRxWarning,
			},
		},
		{
			msg: "protocol violation",
			f: frame{
				idAndFlags:     0x20000008,
				dataLengthCode: 8,
				data: [8]byte{
					0,
					0,
					0x10,
					0x02,
				},
			},
			expected: ErrorFrame{
				ErrorClass:                     ErrorClassProtocolViolation,
				ProtocolError:                  ProtocolViolationErrorBit1,
				ProtocolViolationErrorLocation: ProtocolViolationErrorLocationID28To21,
			},
		},
		{
			msg: "transceiver",
			f: frame{
				idAndFlags:     0x20000010,
				dataLengthCode: 8,
				data: [8]byte{
					0,
					0,
					0,
					0,
					0x07,
				},
			},
			expected: ErrorFrame{
				ErrorClass:       ErrorClassTransceiver,
				TransceiverError: TransceiverErrorCANHShortToGND,
			},
		},
		{
			msg: "controller-specific information",
			f: frame{
				idAndFlags:     0x20000001,
				dataLengthCode: 8,
				data: [8]byte{
					0,
					0,
					0,
					0,
					0,
					1,
					2,
					3,
				},
			},
			expected: ErrorFrame{
				ErrorClass:                    ErrorClassTxTimeout,
				ControllerSpecificInformation: [3]byte{1, 2, 3},
			},
		},
	} {
		tt := tt
		t.Run(tt.msg, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.f.decodeErrorFrame())
		})
	}
}
