package descriptor

import (
	"path"
	"strings"
)

// Database represents a CAN database.
type Database struct {
	// SourceFile of the database.
	//
	// Example:
	//  github.com/einride/can-databases/dbc/j1939.dbc
	SourceFile string
	// Version of the database.
	Version string
	// Messages in the database.
	Messages []*Message
	// Nodes in the database.
	Nodes []*Node
}

func (d *Database) Node(nodeName string) (*Node, bool) {
	for _, n := range d.Nodes {
		if n.Name == nodeName {
			return n, true
		}
	}
	return nil, false
}

func (d *Database) Message(id uint32) (*Message, bool) {
	for _, m := range d.Messages {
		if m.ID == id {
			return m, true
		}
	}
	return nil, false
}

func (d *Database) Signal(messageID uint32, signalName string) (*Signal, bool) {
	message, ok := d.Message(messageID)
	if !ok {
		return nil, false
	}
	for _, s := range message.Signals {
		if s.Name == signalName {
			return s, true
		}
	}
	return nil, false
}

// Description returns the name of the Database.
func (d *Database) Name() string {
	return strings.TrimSuffix(path.Base(d.SourceFile), path.Ext(d.SourceFile))
}
