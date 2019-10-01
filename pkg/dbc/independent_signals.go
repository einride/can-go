package dbc

// Independent signals constants.
//
// DBC files may contain a special message with the following message name and message ID.
//
// This message will have size 0 and may contain duplicate signal names.
const (
	// IndependentSignalsMessageName is the message name used by the special independent signals message.
	IndependentSignalsMessageName Identifier = "VECTOR__INDEPENDENT_SIG_MSG"
	// IndependentSignalsMessageName is the message ID used by the special independent signals message.
	IndependentSignalsMessageID MessageID = 0xc0000000
	// IndependentSignalsMessageSize is the size used by the special independent signals message.
	IndependentSignalsMessageSize = 0
)

// IsIndependentSignalsMessage returns true if m is the special independent signals message.
func IsIndependentSignalsMessage(m *MessageDef) bool {
	return m.Name == IndependentSignalsMessageName &&
		m.MessageID == IndependentSignalsMessageID &&
		m.Size == IndependentSignalsMessageSize
}
