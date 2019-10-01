package dbc

import "fmt"

// EnvironmentVariableType represents the type of an environment variable.
type EnvironmentVariableType uint64

const (
	EnvironmentVariableTypeInteger EnvironmentVariableType = 0
	EnvironmentVariableTypeFloat   EnvironmentVariableType = 1
	EnvironmentVariableTypeString  EnvironmentVariableType = 2
)

// Validate returns an error for invalid environment variable types.
func (e EnvironmentVariableType) Validate() error {
	switch e {
	case EnvironmentVariableTypeInteger:
	case EnvironmentVariableTypeFloat:
	case EnvironmentVariableTypeString:
	default:
		return fmt.Errorf("invalid environment variable type: %v", e)
	}
	return nil
}
