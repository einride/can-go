package dbc

import "fmt"

// AccessType represents the access type of an environment variable.
type AccessType string

const (
	AccessTypeUnrestricted AccessType = "DUMMY_NODE_VECTOR0"
	AccessTypeRead         AccessType = "DUMMY_NODE_VECTOR1"
	AccessTypeWrite        AccessType = "DUMMY_NODE_VECTOR2"
	AccessTypeReadWrite    AccessType = "DUMMY_NODE_VECTOR3"
)

// Validate returns an error for invalid access types.
func (a AccessType) Validate() error {
	switch a {
	case AccessTypeUnrestricted:
	case AccessTypeRead:
	case AccessTypeWrite:
	case AccessTypeReadWrite:
	default:
		return fmt.Errorf("invalid access type: %v", a)
	}
	return nil
}
