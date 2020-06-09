package dbc

import "fmt"

// MessageID represents a message ID.
type MessageID uint32

// ID constants.
const (
	// maxID is the largest valid standard CAN ID.
	maxID = 0x7ff
	// maxExtendedID is the largest valid extended CAN ID.
	maxExtendedID = 0x1fffffff
)

// messageIDExtendedFlag is a bit flag that is set for extended message IDs.
const messageIDExtendedFlag MessageID = 0x80000000

// messageIDIndependentSignals is a special message ID used for the "independent signals" message.
const messageIDIndependentSignals MessageID = 0xc0000000

// IsExtended returns true if the message ID is an extended CAN ID.
func (m MessageID) IsExtended() bool {
	return m != messageIDIndependentSignals && m&messageIDExtendedFlag > 0
}

// ToCAN returns the CAN id value of the message ID (i.e. with bit flags removed).
func (m MessageID) ToCAN() uint32 {
	return uint32(m &^ messageIDExtendedFlag)
}

// Validate returns an error for invalid message IDs.
func (m MessageID) Validate() error {
	if m == messageIDIndependentSignals {
		return nil
	}
	if m.IsExtended() && m.ToCAN() > maxExtendedID {
		return fmt.Errorf("invalid extended ID: %v", m)
	}
	if !m.IsExtended() && m.ToCAN() > maxID {
		return fmt.Errorf("invalid standard ID: %v", m)
	}
	return nil
}
