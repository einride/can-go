package descriptor

import (
	"time"

	"go.einride.tech/can"
)

// Message describes a CAN message.
type Message struct {
	// Description of the message.
	Name string
	// ID of the message.
	ID uint32
	// IsExtended is true if the message is an extended CAN message.
	IsExtended bool
	// Length in bytes.
	Length uint16
	// SendType is the message's send type.
	SendType SendType
	// Description of the message.
	Description string
	// Signals in the message payload.
	Signals []*Signal
	// SenderNode is the name of the node sending the message.
	SenderNode string
	// CycleTime is the cycle time of a cyclic message.
	CycleTime time.Duration
	// DelayTime is the allowed delay between cyclic message sends.
	DelayTime time.Duration
}

// MultiplexerSignal returns the message's multiplexer signal.
func (m *Message) MultiplexerSignal() (*Signal, bool) {
	for _, s := range m.Signals {
		if s.IsMultiplexer {
			return s, true
		}
	}
	return nil, false
}

// Decode decodes a can Payload into a decoded signal array
func (m *Message) Decode(p *can.Payload) []DecodedSignal {

	var data can.Data
	if m.Length <= 8 {
		copy(data[:], p.Data)
	}

	numSignals := len(m.Signals)

	signals := make([]DecodedSignal, numSignals)
	for i, signal := range m.Signals {

		var valueDesc string
		var value float64
		if m.Length > 8 {
			valueDesc, _ = signal.UnmarshalValueDescriptionPayload(p)
			value = signal.UnmarshalPhysicalPayload(p)
		} else {
			valueDesc, _ = signal.UnmarshalValueDescription(data)
			value = signal.UnmarshalPhysical(data)
		}

		s := DecodedSignal{
			Value:       value,
			Description: valueDesc,
			Signal:      signal,
		}

		signals[i] = s

	}
	return signals

}
