package can

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strconv"
)

type jsonFrame struct {
	ID       uint32  `json:"id"`
	Data     *string `json:"data"`
	Length   *uint8  `json:"length"`
	Extended *bool   `json:"extended"`
	Remote   *bool   `json:"remote"`
}

// JSON returns the JSON-encoding of f, using hex-encoding for the data.
//
// Examples:
//
//	{"id":32,"data":"0102030405060708"}
//	{"id":32,"extended":true,"remote":true,"length":4}
func (f Frame) JSON() string {
	switch {
	case f.IsRemote && f.IsExtended:
		return `{"id":` + strconv.Itoa(int(f.ID)) +
			`,"extended":true,"remote":true,"length":` +
			strconv.Itoa(int(f.Length)) + `}`
	case f.IsRemote:
		return `{"id":` + strconv.Itoa(int(f.ID)) +
			`,"remote":true,"length":` +
			strconv.Itoa(int(f.Length)) + `}`
	case f.IsExtended && f.Length == 0:
		return `{"id":` + strconv.Itoa(int(f.ID)) + `,"extended":true}`
	case f.IsExtended:
		return `{"id":` + strconv.Itoa(int(f.ID)) +
			`,"data":"` + hex.EncodeToString(f.Data[:f.Length]) + `"` +
			`,"extended":true}`
	case f.Length == 0:
		return `{"id":` + strconv.Itoa(int(f.ID)) + `}`
	default:
		return `{"id":` + strconv.Itoa(int(f.ID)) +
			`,"data":"` + hex.EncodeToString(f.Data[:f.Length]) + `"}`
	}
}

// MarshalJSON returns the JSON-encoding of f, using hex-encoding for the data.
//
// See JSON for an example of the JSON schema.
func (f Frame) MarshalJSON() ([]byte, error) {
	return []byte(f.JSON()), nil
}

// UnmarshalJSON sets *f using the provided JSON-encoded values.
//
// See MarshalJSON for an example of the expected JSON schema.
//
// The result should be checked with Validate to guard against invalid JSON data.
func (f *Frame) UnmarshalJSON(jsonData []byte) error {
	jf := jsonFrame{}
	if err := json.Unmarshal(jsonData, &jf); err != nil {
		return err
	}
	if jf.Data != nil {
		data, err := hex.DecodeString(*jf.Data)
		if err != nil {
			return fmt.Errorf("failed to hex-decode CAN data: %v: %w", string(jsonData), err)
		}
		f.Data = Data{}
		copy(f.Data[:], data)
		f.Length = uint8(len(data))
	} else {
		f.Data = Data{}
		f.Length = 0
	}
	f.ID = jf.ID
	if jf.Remote != nil {
		f.IsRemote = *jf.Remote
	} else {
		f.IsRemote = false
	}
	if f.IsRemote {
		if jf.Length == nil {
			return fmt.Errorf("missing length field for remote JSON frame: %v", string(jsonData))
		}
		f.Length = *jf.Length
	}
	if jf.Extended != nil {
		f.IsExtended = *jf.Extended
	} else {
		f.IsExtended = false
	}
	return nil
}
