package dbc

import "fmt"

// SignalValueType represents an extended signal value type.
type SignalValueType uint64

const (
	SignalValueTypeInt     SignalValueType = 0
	SignalValueTypeFloat32 SignalValueType = 1
	SignalValueTypeFloat64 SignalValueType = 2
)

// Validate returns an error for invalid signal value types.
func (s SignalValueType) Validate() error {
	switch s {
	case SignalValueTypeInt:
	case SignalValueTypeFloat32:
	case SignalValueTypeFloat64:
	default:
		return fmt.Errorf("invalid signal value type: %v", s)
	}
	return nil
}
