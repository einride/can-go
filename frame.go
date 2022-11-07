package can

import (
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"
)

const (
	idBits         = 11
	extendedIDBits = 29
)

// CAN format constants.
const (
	MaxID         = 0x7ff
	MaxExtendedID = 0x1fffffff
)

// Frame represents a CAN frame.
//
// A Frame is intentionally designed to fit into 16 bytes on common architectures
// and is therefore amenable to pass-by-value and judicious copying.
type Frame struct {
	// ID is the CAN ID
	ID uint32
	// Length is the number of bytes of data in the frame.
	Length uint8
	// Data is the frame data.
	Data Data
	// IsRemote is true for remote frames.
	IsRemote bool
	// IsExtended is true for extended frames, i.e. frames with 29-bit IDs.
	IsExtended bool
}

// Validate returns an error if the Frame is not a valid CAN frame.
func (f *Frame) Validate() error {
	// Validate: ID
	if f.IsExtended && f.ID > MaxExtendedID {
		return fmt.Errorf(
			"invalid extended CAN id: %v does not fit in %v bits",
			f.ID,
			extendedIDBits,
		)
	} else if !f.IsExtended && f.ID > MaxID {
		return fmt.Errorf(
			"invalid standard CAN id: %v does not fit in %v bits",
			f.ID,
			idBits,
		)
	}
	// Validate: Data
	if f.Length > MaxDataLength {
		return fmt.Errorf("invalid data length: %v", f.Length)
	}
	return nil
}

// String returns an ASCII representation the CAN frame.
//
// Format:
//
//	([0-9A-F]{3}|[0-9A-F]{3})#(R[0-8]?|[0-9A-F]{0,16})
//
// The format is compatible with the candump(1) log file format.
func (f Frame) String() string {
	var id string
	if f.IsExtended {
		id = fmt.Sprintf("%08X", f.ID)
	} else {
		id = fmt.Sprintf("%03X", f.ID)
	}
	if f.IsRemote && f.Length == 0 {
		return id + "#R"
	} else if f.IsRemote {
		return id + "#R" + strconv.Itoa(int(f.Length))
	}
	return id + "#" + strings.ToUpper(hex.EncodeToString(f.Data[:f.Length]))
}

// UnmarshalString sets *f using the provided ASCII representation of a Frame.
func (f *Frame) UnmarshalString(s string) error {
	// Split split into parts
	parts := strings.Split(s, "#")
	if len(parts) != 2 {
		return fmt.Errorf("invalid frame format: %v", s)
	}
	idPart, dataPart := parts[0], parts[1]
	var frame Frame
	// Parse: IsExtended
	if len(idPart) != 3 && len(idPart) != 8 {
		return fmt.Errorf("invalid ID length: %v", s)
	}
	frame.IsExtended = len(idPart) == 8
	// Parse: ID
	id, err := strconv.ParseUint(idPart, 16, 32)
	if err != nil {
		return fmt.Errorf("invalid frame ID: %v", s)
	}
	frame.ID = uint32(id)
	if len(dataPart) == 0 {
		*f = frame
		return nil
	}
	// Parse: IsRemote
	if dataPart[0] == 'R' {
		frame.IsRemote = true
		if len(dataPart) > 2 {
			return fmt.Errorf("invalid remote length: %v", s)
		} else if len(dataPart) == 2 {
			dataLength, err := strconv.Atoi(dataPart[1:2])
			if err != nil {
				return fmt.Errorf("invalid remote length: %v: %w", s, err)
			}
			frame.Length = uint8(dataLength)
		}
		*f = frame
		return nil
	}
	// Parse: Length
	if len(dataPart) > 16 || len(dataPart)%2 != 0 {
		return fmt.Errorf("invalid data length: %v", s)
	}
	frame.Length = uint8(len(dataPart) / 2)
	// Parse: Data
	decodedData, err := hex.DecodeString(dataPart)
	if err != nil {
		return fmt.Errorf("invalid data: %v: %w", s, err)
	}
	copy(frame.Data[:], decodedData)
	*f = frame
	return nil
}
