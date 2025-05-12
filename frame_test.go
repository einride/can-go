package can

import (
	"fmt"
	"testing"
	"unsafe"

	"gotest.tools/v3/assert"
)

// If this mocks ever starts failing, the documentation needs to be updated
// to prefer pass-by-pointer over pass-by-value.
func TestFrame_Size(t *testing.T) {
	assert.Assert(t, unsafe.Sizeof(Frame{}) <= 16, "Frame size is <= 16 bytes")
}

func TestFrame_Validate_Error(t *testing.T) {
	for _, tt := range []Frame{
		{ID: MaxID + 1},
		{ID: MaxExtendedID + 1, IsExtended: true},
	} {
		t.Run(fmt.Sprintf("%v", tt), func(t *testing.T) {
			assert.Check(t, tt.Validate() != nil, "should return validation error")
		})
	}
}
