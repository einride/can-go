package can

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"reflect"
	"testing"
	"testing/quick"

	"gotest.tools/v3/assert"
	is "gotest.tools/v3/assert/cmp"
)

func TestFrame_JSON(t *testing.T) {
	for _, tt := range []struct {
		jsonFrame string
		frame     Frame
	}{
		{
			// Standard frame
			jsonFrame: `{"id":42,"data":"00010203"}`,
			frame: Frame{
				ID:     42,
				Length: 4,
				Data:   Data{0x00, 0x01, 0x02, 0x03},
			},
		},
		{
			// Standard frame, no data
			jsonFrame: `{"id":42}`,
			frame:     Frame{ID: 42},
		},
		{
			// Standard remote frame
			jsonFrame: `{"id":42,"remote":true,"length":4}`,
			frame: Frame{
				ID:       42,
				IsRemote: true,
				Length:   4,
			},
		},
		{
			// Extended frame
			jsonFrame: `{"id":42,"data":"0001020304050607","extended":true}`,
			frame: Frame{
				ID:         42,
				IsExtended: true,
				Length:     8,
				Data:       Data{0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07},
			},
		},
		{
			// Extended frame, no data
			jsonFrame: `{"id":42,"extended":true}`,
			frame:     Frame{ID: 42, IsExtended: true},
		},
		{
			// Extended remote frame
			jsonFrame: `{"id":42,"extended":true,"remote":true,"length":8}`,
			frame: Frame{
				ID:         42,
				IsExtended: true,
				IsRemote:   true,
				Length:     8,
			},
		},
	} {
		tt := tt
		t.Run(fmt.Sprintf("JSON|frame=%v", tt.frame), func(t *testing.T) {
			assert.Check(t, is.Equal(tt.jsonFrame, tt.frame.JSON()))
		})
		t.Run(fmt.Sprintf("UnmarshalJSON|frame=%v", tt.frame), func(t *testing.T) {
			var frame Frame
			if err := json.Unmarshal([]byte(tt.jsonFrame), &frame); err != nil {
				t.Fatal(err)
			}
			assert.Check(t, is.DeepEqual(tt.frame, frame))
		})
	}
}

func TestFrame_UnmarshalJSON_Invalid(t *testing.T) {
	var f Frame
	t.Run("invalid JSON", func(t *testing.T) {
		data := `foobar`
		assert.Check(t, f.UnmarshalJSON([]uint8(data)) != nil)
	})
	t.Run("invalid payload", func(t *testing.T) {
		data := `{"id":1,"data":"foobar","extended":false,"remote":false}`
		assert.Check(t, f.UnmarshalJSON([]uint8(data)) != nil)
	})
}

func (Frame) Generate(rand *rand.Rand, size int) reflect.Value {
	f := Frame{
		IsExtended: rand.Intn(2) == 0,
		IsRemote:   rand.Intn(2) == 0,
	}
	if f.IsExtended {
		f.ID = rand.Uint32() & MaxExtendedID
	} else {
		f.ID = rand.Uint32() & MaxID
	}
	f.Length = uint16(rand.Intn(9))
	if !f.IsRemote {
		_, _ = rand.Read(f.Data[:f.Length])
	}
	return reflect.ValueOf(f)
}

func TestPropertyFrame_MarshalUnmarshalJSON(t *testing.T) {
	f := func(f Frame) Frame {
		return f
	}
	g := func(f Frame) Frame {
		f2 := Frame{}
		if err := json.Unmarshal([]uint8(f.JSON()), &f2); err != nil {
			t.Fatal(err)
		}
		return f2
	}
	if err := quick.CheckEqual(f, g, nil); err != nil {
		t.Fatal(err)
	}
}
