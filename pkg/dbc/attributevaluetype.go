package dbc

import "fmt"

// AttributeValueType represents an attribute value type.
type AttributeValueType string

const (
	AttributeValueTypeInt    AttributeValueType = "INT"
	AttributeValueTypeHex    AttributeValueType = "HEX"
	AttributeValueTypeFloat  AttributeValueType = "FLOAT"
	AttributeValueTypeString AttributeValueType = "STRING"
	AttributeValueTypeEnum   AttributeValueType = "ENUM"
)

// Validate returns an error for invalid attribute value types.
func (a AttributeValueType) Validate() error {
	switch a {
	case AttributeValueTypeInt:
	case AttributeValueTypeHex:
	case AttributeValueTypeFloat:
	case AttributeValueTypeString:
	case AttributeValueTypeEnum:
	default:
		return fmt.Errorf("invalid attribute value type: %v", a)
	}
	return nil
}
