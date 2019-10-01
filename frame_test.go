package can

import (
	"fmt"
	"testing"
	"unsafe"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// If this mocks ever starts failing, the documentation needs to be updated
// to prefer pass-by-pointer over pass-by-value.
func TestFrame_Size(t *testing.T) {
	require.True(t, unsafe.Sizeof(Frame{}) <= 16, "Frame size is <= 16 bytes")
}

func TestFrame_Validate_Error(t *testing.T) {
	for _, tt := range []Frame{
		{ID: MaxID + 1},
		{ID: MaxExtendedID + 1, IsExtended: true},
	} {
		tt := tt
		t.Run(fmt.Sprintf("%v", tt), func(t *testing.T) {
			assert.NotNil(t, tt.Validate(), "should return validation error")
		})
	}
}
