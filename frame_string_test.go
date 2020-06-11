package can

import (
	"fmt"
	"testing"

	"gotest.tools/v3/assert"
	is "gotest.tools/v3/assert/cmp"
)

func TestFrame_String(t *testing.T) {
	for _, tt := range []struct {
		frame Frame
		str   string
	}{
		{
			frame: Frame{
				ID:     0x62e,
				Length: 2,
				Data:   Data{0x10, 0x44},
			},
			str: "62E#1044",
		},
		{
			frame: Frame{
				ID:       0x410,
				IsRemote: true,
				Length:   3,
			},
			str: "410#R3",
		},
		{
			frame: Frame{
				ID:     0xd2,
				Length: 2,
				Data:   Data{0xf0, 0x31},
			},
			str: "0D2#F031",
		},
		{
			frame: Frame{ID: 0xee},
			str:   "0EE#",
		},
		{
			frame: Frame{ID: 0},
			str:   "000#",
		},
		{
			frame: Frame{ID: 0, IsExtended: true},
			str:   "00000000#",
		},
		{
			frame: Frame{ID: 0x1234abcd, IsExtended: true},
			str:   "1234ABCD#",
		},
	} {
		tt := tt
		t.Run(fmt.Sprintf("String|frame=%v,str=%v", tt.frame, tt.str), func(t *testing.T) {
			assert.Check(t, is.Equal(tt.str, tt.frame.String()))
		})
		t.Run(fmt.Sprintf("UnmarshalString|frame=%v,str=%v", tt.frame, tt.str), func(t *testing.T) {
			var actual Frame
			if err := actual.UnmarshalString(tt.str); err != nil {
				t.Fatal(err)
			}
			assert.Check(t, is.DeepEqual(actual, tt.frame))
		})
	}
}

func TestParseFrame_Errors(t *testing.T) {
	for _, tt := range []string{
		"foo",                    // invalid
		"foo#",                   // invalid ID
		"0D23#F031",              // invalid ID length
		"62E#104400000000000000", // invalid data length
	} {
		tt := tt
		t.Run(fmt.Sprintf("str=%v", tt), func(t *testing.T) {
			var frame Frame
			err := frame.UnmarshalString(tt)
			assert.ErrorContains(t, err, "invalid")
			assert.Check(t, is.DeepEqual(Frame{}, frame))
		})
	}
}
