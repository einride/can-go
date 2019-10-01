package dbc

import "fmt"

// ObjectType identifies the type of a DBC object.
type ObjectType string

const (
	ObjectTypeUnspecified         ObjectType = ""
	ObjectTypeNetworkNode         ObjectType = "BU_"
	ObjectTypeMessage             ObjectType = "BO_"
	ObjectTypeSignal              ObjectType = "SG_"
	ObjectTypeEnvironmentVariable ObjectType = "EV_"
)

// Validate returns an error for invalid object types.
func (o ObjectType) Validate() error {
	switch o {
	case ObjectTypeUnspecified:
	case ObjectTypeNetworkNode:
	case ObjectTypeMessage:
	case ObjectTypeSignal:
	case ObjectTypeEnvironmentVariable:
	default:
		return fmt.Errorf("invalid object type: %v", o)
	}
	return nil
}
