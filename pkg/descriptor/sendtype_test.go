package descriptor

import (
	"testing"

	"github.com/stretchr/testify/require"
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
		tt := tt
		t.Run(tt.str, func(t *testing.T) {
			var actual SendType
			require.NoError(t, actual.UnmarshalString(tt.str))
			require.Equal(t, tt.expected, actual)
		})
	}
}
