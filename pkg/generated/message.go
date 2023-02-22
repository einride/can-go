// Package generated provides primitives for working with code-generated CAN messages.
package generated

import (
	"fmt"

	"github.com/blueinnovationsgroup/can-go"
	"github.com/blueinnovationsgroup/can-go/pkg/descriptor"
)

// Message represents a code-generated CAN message.
type Message interface {
	can.Message
	fmt.Stringer

	// Descriptor returns the message descriptor.
	Descriptor() *descriptor.Message

	// Reset the message signals to their default values.
	Reset()

	// Frame returns a CAN frame representing the message.
	//
	// A generated message ensures that its signals are valid and is always convertible to a CAN frame.
	Frame() can.Frame
}
