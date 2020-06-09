package descriptor

import "strings"

// SendType represents the send type of a message.
type SendType uint8

//go:generate stringer -type SendType -trimprefix SendType

const (
	// SendTypeNone means the send type is unknown or not specified.
	SendTypeNone SendType = iota
	// SendTypeCyclic means the message is sent cyclically.
	SendTypeCyclic
	// SendTypeEvent means the message is only sent upon event or request.
	SendTypeEvent
)

// UnmarshalString sets the value of *s from the provided string.
func (s *SendType) UnmarshalString(str string) error {
	// TODO: Decide on conventions and make this more strict
	switch strings.ToLower(str) {
	case "cyclic", "cyclicifactive", "periodic", "fixedperiodic", "enabledperiodic", "eventperiodic":
		*s = SendTypeCyclic
	case "event", "onevent":
		*s = SendTypeEvent
	default:
		*s = SendTypeNone
	}
	return nil
}
