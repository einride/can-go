package descriptor

import (
	"testing"

	"gotest.tools/v3/assert"
)

func TestSendType_UnmarshalString(t *testing.T) {
	for _, tt := range []struct {
		str      string
		expected SendType
	}{
		{str: "Cyclic", expected: SendTypeCyclic},
		{str: "Periodic", expected: SendTypeCyclic},
		{str: "OnEvent", expected: SendTypeEvent},
		{str: "Event", expected: SendTypeEvent},
	} {
		t.Run(tt.str, func(t *testing.T) {
			var actual SendType
			assert.NilError(t, actual.UnmarshalString(tt.str))
			assert.Equal(t, tt.expected, actual)
		})
	}
}
